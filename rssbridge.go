package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/feeds"
)

type Item struct {
	Title string
	Link string
	Description string
	Author string
	Created string
}

func main() {
	app := fiber.New()

	// TODO: Add rate limiting or caching or something to prevent (D)DoS'ing.

	app.Get("/store.shopware.com/:plugin.:ext", func (c *fiber.Ctx) error {
		plugin := c.Params("plugin")
		feed_type := strings.ToLower(c.Params("ext"))

		feed, err := shopware_store_plugin(plugin)
		if err != nil {
			log.Fatal(err)
			return c.SendStatus(500)
		}

		var response string
		var feed_err error
		switch feed_type {
			case "atom":
			response, feed_err = feed.ToAtom()
			case "rss":
			response, feed_err = feed.ToRss()
			case "json":
			response, feed_err = feed.ToJSON()
			default:
			c.SendStatus(404)
		}

		if feed_err != nil {
			log.Fatal(err)
			return c.SendStatus(500)
		}

		return c.SendString(response)
	})

	log.Fatal(app.Listen(":3000"))

}

func shopware_store_plugin(id string) (*feeds.Feed, error) {
	url := fmt.Sprintf("https://store.shopware.com/%s", id)
	feed := &feeds.Feed{}

	c := colly.NewCollector(
		colly.AllowedDomains("store.shopware.com"),
		// colly.CacheDir("./.cache/store.shopware.com"),
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
	c.OnHTML(".content--changelog h4", func(h *colly.HTMLElement) {
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

	return feed, nil
}

// vim: noexpandtab
