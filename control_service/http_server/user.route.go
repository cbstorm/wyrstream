package http_server

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/control_service/middlewares"
	"github.com/cbstorm/wyrstream/control_service/services"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/enums"
	"github.com/gofiber/fiber/v2"
)

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/users",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthRole(enums.AUTH_ROLE_ADMIN),
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			fetchArgs := dtos.NewFetchArgs()
			fetchArgs.ParseQueries(c.Queries())
			res, err := services.GetUserService().FetchUsers(fetchArgs, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/users/:id",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthRole(enums.AUTH_ROLE_ADMIN),
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input, err := dtos.NewGetOneInput().SetId(c.Params("id"))
			if err != nil {
				return common.ResponseError(c, err)
			}
			res, err := services.GetUserService().GetOneUser(input, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   DELETE,
	Endpoint: "/users/:id",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthRole(enums.AUTH_ROLE_ADMIN),
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input, err := dtos.NewDeleteOneInput().SetId(c.Params("id"))
			if err != nil {
				return common.ResponseError(c, err)
			}
			res, err := services.GetUserService().DeleteOneUser(input, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})
