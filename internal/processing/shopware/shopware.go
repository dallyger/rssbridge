package shopware

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/gorilla/feeds"

	util "vnbr.de/rssbridge/internal/util"
)

func StorePluginChangelog(id string, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
	var feedErr error
	url := fmt.Sprintf("https://store.shopware.com/en/search?search=%s", id)
	feed := &feeds.Feed{}

	c := colly.NewCollector(
		colly.AllowedDomains("store.shopware.com"),
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

	c.OnHTML("meta[property=\"product:brand\"]", func(h *colly.HTMLElement) {
		author := h.Attr("content")
		if author != "" {
			feed.Author = &feeds.Author{Name: author}
		}
	});
	c.OnHTML("meta[itemprop=\"copyrightHolder\"]", func(h *colly.HTMLElement) {
		feed.Copyright = h.Attr("content")
	});
	c.OnHTML("meta[property=\"og:title\"]", func(h *colly.HTMLElement) {
		feed.Title = h.Attr("content")
	});
	c.OnHTML("meta[property=\"og:description\"]", func(h *colly.HTMLElement) {
		feed.Description = h.Attr("content")
	});
	c.OnHTML("meta[property=\"og:url\"]", func(h *colly.HTMLElement) {
		url = h.Attr("content")
		feed.Link = &feeds.Link{Href: h.Attr("content")}
	});
	c.OnHTML("meta[property=\"og:image\"]", func(h *colly.HTMLElement) {
		feed.Image = &feeds.Image{Link: h.Attr("content")}
	});
	c.OnHTML("h3.changelogs-header", func(h *colly.HTMLElement) {
		desc, err  :=  h.DOM.NextUntil("h4").Html()
		if err != nil {
			log.Fatal(err)
		}
		feed.Items = append(feed.Items, &feeds.Item{
			Id: h.Text,
			Title: h.Text,
			Link: &feeds.Link{Href: url},
			Description: desc,
		})
	});
	c.Visit(url);

	return feed, feedErr
}

// vim: noexpandtab
