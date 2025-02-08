package printables

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/feeds"

	"vnbr.de/rssbridge/internal/util"
)

func SearchModels(ordering string, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
	var feedErr error
	feed := &feeds.Feed{}
	uri := &url.URL{
		Scheme: "https",
		Host: "www.printables.com",
		Path: "model",
		RawQuery: url.Values {
			"ordering": {ordering},
		}.Encode(),
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.printables.com"),
		colly.UserAgent(util.UserAgent()),
	);
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Forwarded-For", ctx.InboundIP)
		r.Headers.Set("X-Forwarded-Host", ctx.InboundHost)
		r.Headers.Set("X-Forwarded-Proto", ctx.InboundProto)
	})
	c.OnError(func(r *colly.Response, err error) {
		feedErr = err
	})

	// TODO: handle infinity scrolling

	c.OnHTML("meta[name=\"og:title\"]", func(h *colly.HTMLElement) {
		if feed.Title == "" {
			feed.Title = h.Attr("content")
		}
	});
	c.OnHTML("meta[name=\"og:description\"]", func(h *colly.HTMLElement) {
		feed.Description = h.Attr("content")
	});
	c.OnHTML("meta[name=\"og:url\"]", func(h *colly.HTMLElement) {
		feed.Link = &feeds.Link{Href: h.Attr("content")}
	});
	c.OnHTML("meta[name=\"og:image\"]", func(h *colly.HTMLElement) {
		feed.Image = &feeds.Image{Link: h.Attr("content")}
	});

	c.OnHTML("article.card", func(h *colly.HTMLElement) {
		title := strings.TrimSpace(h.DOM.Find("a.h").First().Text())
		description, _ := h.DOM.Find("div.stats-bar").First().Html()
		// TODO: include author and other metadata in description
		// TODO: include featured model badge
		// TODO: handle NSFW

		lnk, ok := h.DOM.Find("a.h").First().Attr("href")
		if ok == false {
			// TODO: log: link not found
		}

		// example lnk: https://www.printables.com/model/123456-title-here
		id := strings.Split(strings.Split(lnk, "-")[0], "/")[2]

		imgSrcset, hasImgSrcset := h.DOM.Find("a.card-image").First().Find("source").Attr("srcset")
		if hasImgSrcset {
			description = fmt.Sprintf(
				"<img src=\"%s\"><p>%s</p>",
				strings.Split(imgSrcset, " ")[2],
				description,
			)
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Id: id,
			Title: title,
			Link: &feeds.Link{Href: fmt.Sprintf("https://www.printables.com%s", lnk)},
			Description: description,
		})
	});

	c.Visit(uri.String());

	return feed, feedErr

}

// vim: noexpandtab
