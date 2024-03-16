package kleinanzeigen

import (
	"dallyger/rssbridge/internal/util"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/feeds"
)

func Search(search string, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
	feed := &feeds.Feed{}
	uri := &url.URL{
		Scheme: "https",
		Host: "www.kleinanzeigen.de",
		Path: "s-suchanfrage.html",
		RawQuery: url.Values {
			"keywords": {search},
			"categoryId": {""},
			"locationStr": {""},
			"locationId": {""},
			"radius": {"0"},
			"sortingField": {"SORTING_DATE"},
			"pageNum": {"1"},
			"action": {"find"},
			"maxPrice": {""},
			"minPrice": {""},
		}.Encode(),
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.kleinanzeigen.de", "kleinanzeigen.de"),
		// Kleinanzeigen will close the connection if it encounters a link
		// TODO: set version in UA
		colly.UserAgent("Mozilla/5.0 (compatible; rssbridge/0.0.0)"),
	);
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Forwarded-For", ctx.InboundIP)
		r.Headers.Set("X-Forwarded-Host", ctx.InboundHost)
		r.Headers.Set("X-Forwarded-Proto", ctx.InboundProto)
	})

	c.OnHTML("meta[name=\"author\"]", func(h *colly.HTMLElement) {
		author := h.Attr("content")
		if author != "" {
			feed.Author = &feeds.Author{Name: author}
		}
	});
	c.OnHTML("meta[itemprop=\"copyrightHolder\"]", func(h *colly.HTMLElement) {
		feed.Copyright = h.Attr("content")
	});
	c.OnHTML("title", func(h *colly.HTMLElement) {
		feed.Title = h.Text
	});
	c.OnHTML("meta[property=\"og:description\"]", func(h *colly.HTMLElement) {
		feed.Description = h.Attr("content")
	});
	c.OnHTML("meta[property=\"og:url\"]", func(h *colly.HTMLElement) {
		feed.Link = &feeds.Link{Href: h.Attr("content")}
	});
	c.OnHTML("meta[property=\"og:image\"]", func(h *colly.HTMLElement) {
		feed.Image = &feeds.Image{Link: h.Attr("content")}
	});

	c.OnHTML("article.aditem", func(h *colly.HTMLElement) {
		itemHref := fmt.Sprintf("https://www.kleinanzeigen.de/%s", h.Attr("data-href"))
		title := strings.TrimSpace(h.DOM.Find("a.ellipsis").First().Text())
		description := h.DOM.Find("p.aditem-main--middle--description").First().Text()
		price := strings.TrimSpace(h.DOM.Find("p.aditem-main--middle--price-shipping--price").First().Text())
		if price != "" {
			title = fmt.Sprintf("[%s] %s", price, title)
		}

		date := strings.TrimSpace(h.DOM.Find(".icon-calendar-open").Parent().Text())
		var created time.Time
		if strings.HasPrefix(date, "Heute") || strings.HasPrefix(date, "Gestern") {
			created, _ = time.Parse("15:04", strings.SplitN(date, " ", 2)[1])
			created = created.AddDate(time.Now().Year(), int(time.Now().Month()) - 1, time.Now().Day() - 1)
		} else {
			created, _ = time.Parse("02.01.2006", date)
		}

		imgSrcset, hasImgSrcset := h.DOM.Find("img[srcset]").First().Attr("srcset")
		if hasImgSrcset {
			description = fmt.Sprintf(
				"<p><img src=\"%s\"> %s</p>",
				imgSrcset,
				description,
			)
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Id: h.Attr("data-adid"),
			Title: title,
			Link: &feeds.Link{Href: itemHref},
			Description: description,
			Created: created,
			Updated: created,
		})
	});

	// Visit link found on page.
	// Only those links are visited which are in AllowedDomains.
	// Already visited links are skipped.
	c.OnHTML("a.pagination-page", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		c.Visit(h.Request.AbsoluteURL(link))
	})

	c.Visit(uri.String());

	return feed, nil

}

// vim: noexpandtab
