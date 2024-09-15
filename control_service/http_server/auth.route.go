package http_server

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/control_service/middlewares"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/exceptions"
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
				e := exceptions.Err_BAD_REQUEST().SetMessage(err.Error())
				return common.ResponseError(c, e)
			}
			response := &dtos.UserLoginResponse{}
			if err := result.Decode(response); err != nil {
				return err
			}
			return c.Status(200).JSON(response)
		},
	},
	Enable: true,
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/auth/user/create_account",
	Handlers: []func(common.IHttpContext) error{
		func(c common.IHttpContext) error {
			return c.Status(200).JSON("auth")
		},
	},
	Enable: true,
})
