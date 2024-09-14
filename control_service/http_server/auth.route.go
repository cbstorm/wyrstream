package http_server

import "github.com/cbstorm/wyrstream/control_service/common"

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/auth/user/login",
	Handlers: []func(common.IHttpContext) error{
		func(ic common.IHttpContext) error {
			return ic.Status(200).JSON("auth")
		},
	},
	Enable: true,
})
