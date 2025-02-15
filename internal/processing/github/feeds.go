package github

import (
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"vnbr.de/rssbridge/internal/util"
)

func NotificationFeed(token string, ctx *util.ScrapeCtx) (feed *feeds.Feed, err error) {
	feed = &feeds.Feed{}

	status, resp, err := Request("GET", "https://api.github.com/notifications", token, ctx)

	if err != nil {
		return
	}

	if status != 200 {
		// TODO: handle 401 (invalid token) and 403 (invalid permission) by
		// returning a valid feed with an error message. This way the user can
		// see the issue without looking at the server logs (if that is even
		// possible. e.g. on the go, shared instances).
		return feed, fmt.Errorf("github.com: received non-successful status code [%d]", status)
	}

	notifs := []Notification{}
	if err := json.Unmarshal(resp, &notifs); err != nil {
		return feed, fmt.Errorf("github.com: parse json: %s", err)
	}

	for _, notif := range notifs {
		status, resp, err = Request("GET", notif.Subject.Url, token, ctx)
		switch notif.Subject.Type {
		case "Discussion":
			dis := Issue{}
			if err := json.Unmarshal(resp, &dis); err != nil {
				return feed, fmt.Errorf("github.com: parse issue: %s", err)
			}
			feed.Add(&feeds.Item{
				Author:      &feeds.Author{Name: dis.User.Login},
				Content:     renderContent(dis.Body),
				Created:     dis.Created_at.Time,
				Description: notif.Subject.Title,
				Id:          notif.Id,
				Link:        &feeds.Link{Href: dis.Html_url},
				Source:      nil,
				Title:       fmt.Sprintf("%s (Post: %s #%d)", dis.Title, notif.Repository.Full_name, dis.Number),
				Updated:     notif.Updated_at.Time,
			})

		case "Issue":
			issue := Issue{}
			if err := json.Unmarshal(resp, &issue); err != nil {
				return feed, fmt.Errorf("github.com: parse issue: %s", err)
			}
			feed.Add(&feeds.Item{
				Author:      &feeds.Author{Name: issue.User.Login},
				Content:     renderContent(issue.Body),
				Created:     issue.Created_at.Time,
				Description: notif.Subject.Title,
				Id:          notif.Id,
				Link:        &feeds.Link{Href: issue.Html_url},
				Source:      nil,
				Title:       fmt.Sprintf("%s (Issue: %s #%d)", issue.Title, notif.Repository.Full_name, issue.Number),
				Updated:     notif.Updated_at.Time,
			})

		case "PullRequest":
			pr := Issue{}
			if err := json.Unmarshal(resp, &pr); err != nil {
				return feed, fmt.Errorf("github.com: parse pr: %s", err)
			}
			feed.Add(&feeds.Item{
				Author:      &feeds.Author{Name: pr.User.Login},
				Content:     renderContent(pr.Body),
				Created:     pr.Created_at.Time,
				Description: notif.Subject.Title,
				Id:          notif.Id,
				Link:        &feeds.Link{Href: pr.Html_url},
				Source:      nil,
				Title:       fmt.Sprintf("%s (PR: %s #%d)", pr.Title, notif.Repository.Full_name, pr.Number),
				Updated:     notif.Updated_at.Time,
			})

		case "Release":
			release := Release{}
			if err := json.Unmarshal(resp, &release); err != nil {
				return feed, fmt.Errorf("github.com: parse pr: %s", err)
			}
			feed.Add(&feeds.Item{
				Author:      &feeds.Author{Name: release.Author.Login},
				Content:     renderContent(release.Body),
				Created:     release.Created_at.Time,
				Description: notif.Subject.Title,
				Id:          notif.Id,
				Link:        &feeds.Link{Href: release.Html_url},
				Source:      nil,
				Title:       fmt.Sprintf("%s (Release: %s in %s)", release.Name, release.Tag_name, notif.Repository.Full_name),
				Updated:     notif.Updated_at.Time,
			})

		default:
			slog.Warn("github.com: skip notification of unknown type", "type", notif.Subject.Type)
		}

	}

	updated := time.Now()
	if len(feed.Items) > 0 {
		updated = slices.MaxFunc(feed.Items, func(a *feeds.Item, b *feeds.Item) int {
			return a.Updated.Compare(b.Updated)
		}).Updated
	}

	feed = &feeds.Feed{
		Items:   feed.Items,
		Link:    &feeds.Link{Href: "https://github.com/notifications"},
		Title:   "github.com notifications",
		Updated: updated,
	}

	return
}

func renderContent(content string) string {
	return strings.ReplaceAll(html.EscapeString(content), "\n", "<br>\n")
}
