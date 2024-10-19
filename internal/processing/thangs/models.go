package thangs

import (
	"dallyger/rssbridge/internal/util"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/feeds"
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

	c.OnHTML("section.model-card", func(h *colly.HTMLElement) {
		title := h.DOM.Find("h4").First().Text()
		desc, _ := h.DOM.Html()
		lnk := h.DOM.Find("a").Eq(1).AttrOr("href", "")

		lnk_chunks := strings.Split(lnk, "-")
		id := lnk_chunks[len(lnk_chunks)-1]

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

// vim: noexpandtab
