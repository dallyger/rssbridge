package ebay

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/feeds"

	"vnbr.de/rssbridge/internal/util"
)

func SearchDE(search string, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
	uri := &url.URL{
		Scheme: "https",
		Host: "www.ebay.de",
		Path: "sch/i.html",
		RawQuery: url.Values {
			"_nkw": { search },
			"_sop": { "12" },
		}.Encode(),
	}

	var feedErr error
	feed := &feeds.Feed{
		Link: &feeds.Link{Href: uri.String()},
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.ebay.de", "ebay.de"),
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

	c.OnHTML(".gf-legal", func(h *colly.HTMLElement) {
		feed.Copyright = h.Text
	});
	c.OnHTML("title", func(h *colly.HTMLElement) {
		feed.Title = h.Text
	});
	c.OnHTML("meta[name=\"description\"]", func(h *colly.HTMLElement) {
		feed.Description = h.Attr("content")
	});

	c.OnHTML("li.s-item", func(h *colly.HTMLElement) {

		// Everything below is not a full match and returns junk
		if h.DOM.PrevAll().Is("li.srp-river-answer") {
			return
		}

		itemHref, _ := h.DOM.Find(".s-item__link").First().Attr("href")
		title := strings.TrimSpace(h.DOM.Find(".s-item__title").First().Text())
		description := h.DOM.Find(".s-item__info").First().Text()
		price := strings.TrimSpace(h.DOM.Find(".s-item__price").First().Text())
		if price != "" {
			title = fmt.Sprintf("[%s] %s", price, title)
		}

		img, hasImg := h.DOM.Find("img[src]").First().Attr("src")
		if hasImg {
			description = fmt.Sprintf(
				"<img src=\"%s\"> <p>%s</p>",
				img,
				description,
			)
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Id: h.Attr("id"),
			Title: title,
			Link: &feeds.Link{Href: itemHref},
			Description: description,
		})
	});

	// Visit link found on page.
	// Only those links are visited which are in AllowedDomains.
	// Already visited links are skipped.
	c.OnHTML("a.pagination__next", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		c.Visit(h.Request.AbsoluteURL(link))
	})


	c.Visit(uri.String());

	return feed, feedErr

}

// vim: noexpandtab
