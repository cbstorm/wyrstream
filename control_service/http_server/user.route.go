package http_server

import (
	"github.com/gofiber/fiber/v2"
)

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/user/users",
	Handlers: []func(c *fiber.Ctx) error{
		func(ctx *fiber.Ctx) error {
			return ctx.Status(200).JSON("users")
		},
	},
})
