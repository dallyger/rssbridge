package forgejo

import (
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"vnbr.de/rssbridge/internal/util"
)

func NotificationFeed(ctx *util.ScrapeCtx, domain string) (feed *feeds.Feed, err error) {
	c := Client{ctx, domain}
	feed = &feeds.Feed{
		Title:   fmt.Sprintf("%s notifications", domain),
		Link:    &feeds.Link{Href: fmt.Sprintf("https://%s/notifications", domain)},
		Updated: time.Now(),
	}

	notifs := []Notification{}
	if msg, ok := fetch(c, "notifications", &notifs); !ok {
		if msg == "" {
			msg = "Please check server logs for more information."
		}
		feed.Items = []*feeds.Item{
			{
				Title:   "Something went wrong",
				Content: renderContent(msg),
				Updated: time.Now(),
			},
		}
		return
	}

	for _, notif := range notifs {
		switch notif.Subject.Type {

		case "Issue", "Pull":
			id := strconv.FormatUint(uint64(notif.Id), 10)
			no := urlToNum(notif.Subject.Url)
			title := fmt.Sprintf("%s (%s: %s #%s)", notif.Subject.Title, notif.Subject.Type, notif.Repository.FullName, no)
			comment := Comment{}
			if notif.Subject.LatestCommentUrl != "" {
				if msg, ok := fetch(c, notif.Subject.LatestCommentUrl, &comment); !ok {
					// Failed fetching comment
					feed.Add(&feeds.Item{
						Id:          id,
						Title:       title,
						Description: msg,
						Link:        &feeds.Link{Href: notif.Subject.HtmlUrl},
						Updated:     notif.UpdatedAt.Time,
					})
					continue
				} else {
					// Fetched comment
					feed.Add(&feeds.Item{
						Author:      &feeds.Author{Name: comment.User.Login},
						Content:     renderContent(comment.Body),
						Created:     comment.CreatedAt.Time,
						Description: notif.Subject.Title,
						Id:          id,
						Link:        &feeds.Link{Href: comment.HtmlUrl},
						Title:       title,
						Updated:     notif.UpdatedAt.Time,
					})
					continue
				}
			}
			// No comment
			feed.Add(&feeds.Item{
				Id:      id,
				Title:   title,
				Updated: notif.UpdatedAt.Time,
			})

		default:
			slog.Warn("forgejo: parse notification of unknown type", "domain", domain, "type", notif.Subject.Type)
			feed.Add(&feeds.Item{
				Id:          strconv.FormatUint(uint64(notif.Id), 10),
				Title:       fmt.Sprintf("%s (%s: %s)", notif.Subject.Title, notif.Subject.Type, notif.Repository.FullName),
				Content:     "Description unavailable.\nUnsupported notification type.",
				Description: notif.Subject.Title,
				Updated:     notif.UpdatedAt.Time,
				Link:        &feeds.Link{Href: notif.Subject.HtmlUrl},
			})
		}

	}

	if len(feed.Items) > 0 {
		updated := feed.Updated
		updated = slices.MaxFunc(feed.Items, func(a *feeds.Item, b *feeds.Item) int {
			return a.Updated.Compare(b.Updated)
		}).Updated
		feed.Updated = updated
	}

	return
}

func fetch[T any](c Client, url string, into *T) (string, bool) {
	status, resp, err := c.Request("GET", url)

	if err != nil {
		slog.Warn("forgejo: fetch data", "domain", c.Domain, "error", err)
		return "", false
	}

	if slices.Contains([]int{401, 403}, status) {
		msg := Error{}
		if err := json.Unmarshal(resp, &msg); err != nil {
			slog.Warn("forgejo: parse error message", "domain", c.Domain, "error", err)
			return "", false
		}
		return msg.Message, false
	}

	if status != 200 {
		slog.Warn("forgejo: received non-successful status", "domain", c.Domain, "status", status)
		return "", false
	}

	if err := json.Unmarshal(resp, &into); err != nil {
		slog.Warn("forgejo: parse failure", "domain", c.Domain, "type", reflect.TypeOf(into).Name(), "error", err)
		return "", false
	}

	return "", true
}

func renderContent(content string) string {
	return strings.ReplaceAll(html.EscapeString(content), "\n", "<br>\n")
}

func urlToNum(url string) string {
	chunks := strings.Split(url, "/")
	return chunks[len(chunks)-1]
}
