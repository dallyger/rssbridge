package router

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gorilla/feeds"

	"vnbr.de/rssbridge/internal/processing/ebay"
	"vnbr.de/rssbridge/internal/processing/forgejo"
	"vnbr.de/rssbridge/internal/processing/github"
	kleinanzeigen "vnbr.de/rssbridge/internal/processing/kleinanzeigen"
	"vnbr.de/rssbridge/internal/processing/printables"
	shopware "vnbr.de/rssbridge/internal/processing/shopware"
	"vnbr.de/rssbridge/internal/processing/thangs"
	"vnbr.de/rssbridge/internal/util"
)

type Item struct {
	Title string
	Link string
	Description string
	Author string
	Created string
}

func Run() {
	app := fiber.New(fiber.Config{
		ProxyHeader: fiber.HeaderXForwardedFor,
		EnableTrustedProxyCheck: true,
		TrustedProxies: strings.Split(os.Getenv("TRUSTED_PROXIES"), ";"),
	})

	// Log incoming requests
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${ua} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// TODO: Add rate limiting or caching or something to prevent (D)DoS'ing.

	app.Get("/up", func (c *fiber.Ctx) error {
		return c.Status(200).JSON(&fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/codeberg.org/notifications.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return forgejo.NotificationFeed(ctx, "codeberg.org")
		}),
	)

	app.Get("/ebay.de.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return ebay.SearchDE(c.Query("query"), ctx)
		}),
	)

	app.Get("/forgejo/:domain/notifications.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return forgejo.NotificationFeed(ctx, c.Query("domain"))
		}),
	)

	app.Get("/gitea.com/notifications.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return forgejo.NotificationFeed(ctx, "gitea.com")
		}),
	)

	app.Get("/github.com/notifications.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return github.NotificationFeed(ctx)
		}),
	)

	app.Get("/kleinanzeigen.de", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return kleinanzeigen.Search(c.Query("query"), ctx)
		}),
	)

	app.Get("/kleinanzeigen.de.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return kleinanzeigen.Search(c.Query("query"), ctx)
		}),
	)

	app.Get("/printables.com/liked.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return printables.SearchModels("liked", ctx)
		}),
	)

	app.Get("/printables.com/trending.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return printables.SearchModels("", ctx)
		}),
	)

	// example: /store.shopware.com/swag136939272659f/shopware-6-sicherheits-plugin.html.rss
	app.Get("/store.shopware.com/:plugin/:slug.html.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return shopware.StorePluginChangelog(c.Params("plugin"), ctx)
		}),
	)

	// example: /store.shopware.com/swag136939272659f/shopware-6-sicherheits-plugin.rss
	app.Get("/store.shopware.com/:plugin/:slug.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return shopware.StorePluginChangelog(c.Params("plugin"), ctx)
		}),
	)

	// example: /store.shopware.com/swag136939272659f.rss
	app.Get("/store.shopware.com/:plugin.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return shopware.StorePluginChangelog(c.Params("plugin"), ctx)
		}),
	)

	app.Get("/thangs.com/downloads.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return thangs.Downloads(ctx)
		}),
	)

	app.Get("/thangs.com/popular.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return thangs.Popular(ctx)
		}),
	)

	app.Get("/thangs.com/recent.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return thangs.Recent(ctx)
		}),
	)

	app.Get("/thangs.com/trending.:ext", createFeedResponse(
		func (c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) {
			return thangs.Trending(ctx)
		}),
	)

	log.Fatal(app.Listen(":3000"))
}

func createFeedResponse(handler func(c *fiber.Ctx, ctx *util.ScrapeCtx) (*feeds.Feed, error) ) func(c *fiber.Ctx) error {
	return func (c *fiber.Ctx) error {

		ip := c.IP()
		if c.IsProxyTrusted() {
			ip = strings.Join(c.IPs(), ", ")
		}

		ctx := &util.ScrapeCtx{
			InboundIP: ip,
			InboundHost: string(c.Context().Host()),
			InboundProto: c.Protocol(),
			Token: c.Query("token", string(c.Request().URI().Password())),
		}

		feed_type := strings.ToLower(c.Params("ext", "rss"))
		if len(feed_type) > 4 {
			feed_type = feed_type[len(feed_type)-4:]
			feed_type = strings.TrimLeft(feed_type, ".")
		}
		feed, err := handler(c, ctx)

		if err != nil {
			log.Print(err)
			return c.SendStatus(500)
		}

		if feed == nil {
			slog.Error("Feed handler returned nil")
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
			log.Print(err)
			return c.SendStatus(500)
		}

		return c.SendString(response)
	}
}

// vim: noexpandtab
