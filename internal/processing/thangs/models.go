package thangs

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gorilla/feeds"

	"vnbr.de/rssbridge/internal/util"
)

func Downloads(ctx *util.ScrapeCtx) (*feeds.Feed, error) {
	return Models("downloads", ctx)
}

func Popular(ctx *util.ScrapeCtx) (*feeds.Feed, error) {
	return Models("likes", ctx)
}

func Recent(ctx *util.ScrapeCtx) (*feeds.Feed, error) {
	return Models("date", ctx)
}

func Trending(ctx *util.ScrapeCtx) (*feeds.Feed, error) {
	return Models("trending", ctx)
}

func Models(sort string, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
	var feedErr error
	feed := &feeds.Feed{}
	uri := &url.URL{
		Scheme: "https",
		Host:   "thangs.com",
		Path:   "",
		RawQuery: url.Values{
			"sort": {sort},
		}.Encode(),
	}

	c := colly.NewCollector(
		colly.AllowedDomains("thangs.com"),
		colly.UserAgent(util.UserAgent()),
	)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Forwarded-For", ctx.InboundIP)
		r.Headers.Set("X-Forwarded-Host", ctx.InboundHost)
		r.Headers.Set("X-Forwarded-Proto", ctx.InboundProto)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Fprintf(os.Stderr, "thangs.com: %s\n", err)
		feedErr = err
	})

	// gather feed metadata
	c.OnHTML("head", func(h *colly.HTMLElement) {
		// ensure we're on the index page
		if strings.Contains(h.Request.URL.Path, "page=") {
			return
		}

		feed.Title = h.DOM.Find("meta[property=\"og:title\"]").First().AttrOr("content", "thangs.com")
		feed.Description = h.DOM.Find("meta[property=\"og:description\"]").First().AttrOr("content", "")
		feed.Link = &feeds.Link{
			Href: h.DOM.Find("link[rel=\"canonical\"]").First().AttrOr("href", ""),
		}
		feed.Image = &feeds.Image{
			Link: h.DOM.Find("link[rel=\"apple-touch-icon\"]").First().AttrOr("href", ""),
		}
	})

	// walk pagination
	c.OnHTML("a", func(h *colly.HTMLElement) {
		// The index page contains links for the first 4 pages. Those should be
		// enough, because anything below should be not as relevant.

		// ensure we're on the index page
		if strings.Contains(h.Request.URL.Path, "page=") {
			return
		}

		// ensure it's a pagination link
		if !strings.Contains(h.DOM.AttrOr("class", ""), "PaginatorLink-") {
			return
		}

		if lnk := h.DOM.AttrOr("href", ""); lnk != "" {
			c.Visit(h.Request.AbsoluteURL(lnk))
		}
	})

	c.OnHTML("section.model-card", func(h *colly.HTMLElement) {
		var id string
		var lnk string
		title := h.DOM.Find("h4").First().Text()
		desc, _ := h.DOM.Html()

		h.DOM.Find("a").Each(func(i int, s *goquery.Selection) {
			if id == "" && modelIdFromUrl(s.AttrOr("href", "")) != "" {
				id = modelIdFromUrl(s.AttrOr("href", ""))
				lnk = s.AttrOr("href", "")
			}
		})

		if id == "" {
			slog.Error("Failed extracting thangs.com model data", "title", title, "url", lnk)
			return
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Id:    id,
			Title: title,
			Link: &feeds.Link{
				Href: fmt.Sprintf("https://www.thangs.com%s", lnk),
			},
			Description: desc,
		})
	})

	c.Visit(uri.String())

	return feed, feedErr
}

func modelIdFromUrl(lnk string) string {
	lnk_chunks := strings.Split(lnk, "-")

	if len(lnk_chunks) == 1 {
		return ""
	}

	id := lnk_chunks[len(lnk_chunks)-1]

	if _, err := strconv.Atoi(id); err == nil && len(id) > 0 {
		return id
	}

	return ""
}

// vim: noexpandtab
