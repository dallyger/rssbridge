package kleinanzeigen

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/feeds"
)

func Search(search string) (*feeds.Feed, error) {
	url := fmt.Sprintf("https://www.kleinanzeigen.de/s-%s/k0", strings.ToLower(strings.ReplaceAll(search, " ", "-")))
	feed := &feeds.Feed{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.kleinanzeigen.de", "kleinanzeigen.de"),
		colly.UserAgent("dallyger/rssbridge"),
	);

	c.OnHTML("meta[name=\"author\"]", func(h *colly.HTMLElement) {
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

	c.OnHTML("article.aditem", func(h *colly.HTMLElement) {
		itemHref := fmt.Sprintf("https://www.kleinanzeigen.de/%s", h.Attr("data-href"))
		feed.Items = append(feed.Items, &feeds.Item{
			Id: h.Attr("data-adid"),
			Title: h.DOM.Find("a.ellipsis").First().Text(),
			Link: &feeds.Link{Href: itemHref},
			Description: h.DOM.Find("p.aditem-main--middle--description").First().Text(),
		})
	});

	// Visit link found on page.
	// Only those links are visited which are in AllowedDomains.
	// Already visited links are skipped.
	c.OnHTML("a.pagination-page", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		c.Visit(h.Request.AbsoluteURL(link))
	})

	c.Visit(url);

	return feed, nil

}

// vim: noexpandtab
