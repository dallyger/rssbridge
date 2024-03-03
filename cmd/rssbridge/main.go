package main

import (
	kleinanzeigen "dallyger/rssbridge/internal/processing/kleinanzeigen"
	shopware "dallyger/rssbridge/internal/processing/shopware"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	app := fiber.New(fiber.Config{
		ProxyHeader: fiber.HeaderXForwardedFor,
		EnableTrustedProxyCheck: true,
		TrustedProxies: strings.Split(os.Getenv("TRUSTED_PROXIES"), ";"),
	})

	// Log incoming requests
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${ua} | ${error}\n",
	}))

	// TODO: Add rate limiting or caching or something to prevent (D)DoS'ing.

	app.Get("/up", func (c *fiber.Ctx) error {
		return c.Status(200).JSON(&fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/kleinanzeigen.de", createFeedResponse(func (c *fiber.Ctx) (*feeds.Feed, error) {
		query := c.Query("query")
		return kleinanzeigen.Search(query)
	}))

	app.Get("/kleinanzeigen.de.:ext", createFeedResponse(func (c *fiber.Ctx) (*feeds.Feed, error) {
		query := c.Query("query")
		return kleinanzeigen.Search(query)
	}))

	app.Get("/store.shopware.com/:plugin.:ext", createFeedResponse(func (c *fiber.Ctx) (*feeds.Feed, error) {
		plugin := fmt.Sprintf("%s.%s", c.Params("plugin"), "html")
		return shopware.StorePluginChangelog(plugin)
	}))

	app.Get("/store.shopware.com/:plugin.html.:ext", createFeedResponse(func (c *fiber.Ctx) (*feeds.Feed, error) {
		plugin := fmt.Sprintf("%s.%s", c.Params("plugin"), "html")
		return shopware.StorePluginChangelog(plugin)
	}))

	log.Fatal(app.Listen(":3000"))
}

func createFeedResponse(handler func(c *fiber.Ctx) (*feeds.Feed, error) ) func(c *fiber.Ctx) error {
	return func (c *fiber.Ctx) error {
		feed_type := strings.ToLower(c.Params("ext", "rss"))
		feed, err := handler(c)

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
	}
}

// vim: noexpandtab
