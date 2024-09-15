package http_server

import "github.com/cbstorm/wyrstream/control_service/common"

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/user/users",
	Handlers: []func(common.IHttpContext) error{
		func(ctx common.IHttpContext) error {
			return ctx.Status(200).JSON("users")
		},
	},
})
