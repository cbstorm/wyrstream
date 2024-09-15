package http_server

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/control_service/middlewares"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/nats_service"
	"github.com/gofiber/fiber/v2"
)

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   POST,
	Endpoint: "/auth/user/login",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.LoggingMiddleware,
		middlewares.BodyRequiredMiddleware,
		func(c *fiber.Ctx) error {
			result, err := nats_service.GetNATSService().Request("auth.user.login", c.BodyRaw())
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, result)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   POST,
	Endpoint: "/auth/user/create_account",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.LoggingMiddleware,
		middlewares.BodyRequiredMiddleware,
		func(c *fiber.Ctx) error {
			result, err := nats_service.GetNATSService().Request("auth.user.create_account", c.BodyRaw())
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, result)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/auth/user/get_me",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthMiddleware,
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input := &dtos.UserGetMeInput{
				UserId: req_ctx.GetObjId(),
			}
			result, err := nats_service.GetNATSService().Request("auth.user.get_me", input)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, result)
		},
	},
})
