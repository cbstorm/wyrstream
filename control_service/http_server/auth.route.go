package http_server

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/control_service/middlewares"
	"github.com/cbstorm/wyrstream/lib/nats_service"
)

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   POST,
	Endpoint: "/auth/user/login",
	Handlers: []func(common.IHttpContext) error{
		middlewares.BodyRequiredMiddleware,
		func(c common.IHttpContext) error {
			result, err := nats_service.GetNATSService().Request("auth.user.login", c.BodyRaw())
			if err != nil {
				return common.ResponseError(c, err)
			}
			return c.Status(200).JSON(result)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   POST,
	Endpoint: "/auth/user/create_account",
	Handlers: []func(common.IHttpContext) error{
		middlewares.BodyRequiredMiddleware,
		func(c common.IHttpContext) error {
			result, err := nats_service.GetNATSService().Request("auth.user.create_account", c.BodyRaw())
			if err != nil {
				return common.ResponseError(c, err)
			}
			return c.Status(200).JSON(result)
		},
	},
})
