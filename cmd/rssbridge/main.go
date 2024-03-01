package main

import (
	shopware "dallyger/rssbridge/internal/processing"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	// Log incoming requests
	app.Use(logger.New())

	// TODO: Add rate limiting or caching or something to prevent (D)DoS'ing.

	app.Get("/up", func (c *fiber.Ctx) error {
		return c.Status(200).JSON(&fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/store.shopware.com/:plugin.:plugin_ext.:ext", func (c *fiber.Ctx) error {
		plugin := fmt.Sprintf("%s.%s", c.Params("plugin"), c.Params("plugin_ext"))
		feed_type := strings.ToLower(c.Params("ext"))

		feed, err := shopware.StorePluginChangelog(plugin)
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

// vim: noexpandtab
