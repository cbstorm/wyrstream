package http_server

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/control_service/middlewares"
	"github.com/cbstorm/wyrstream/control_service/services"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/enums"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/gofiber/fiber/v2"
)

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/streams",
	Handlers: []func(*fiber.Ctx) error{
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			fetchArgs := dtos.NewFetchArgs()
			fetchArgs.ParseQueries(c.Queries())
			fetchArgs.SetFilter("is_publishing", true)
			res, err := services.GetStreamService().FetchStreams(fetchArgs, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/streams/my_streams",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthRole(enums.AUTH_ROLE_USER),
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			fetchArgs := dtos.NewFetchArgs()
			fetchArgs.ParseQueries(c.Queries())
			fetchArgs.SetFilter("publisher_id", req_ctx.GetObjId())
			res, err := services.GetStreamService().FetchStreams(fetchArgs, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/streams/:id",
	Handlers: []func(*fiber.Ctx) error{
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input, err := dtos.NewGetOneInput().SetId(c.Params("id"))
			if err != nil {
				return common.ResponseError(c, err)
			}
			res, err := services.GetStreamService().GetOneStream(input, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   POST,
	Endpoint: "/streams",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.Alert,
		middlewares.AuthRole(enums.AUTH_ROLE_USER),
		middlewares.BodyRequiredMiddleware,
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input := dtos.NewCreateOneStreamInput()
			if err := c.BodyParser(input); err != nil {
				return common.ResponseError(c, exceptions.Err_BAD_REQUEST().SetMessage(err.Error()))
			}
			if err := input.Validate(); err != nil {
				return common.ResponseError(c, exceptions.Err_BAD_REQUEST().SetMessage(err.Error()))
			}
			res, err := services.GetStreamService().CreateOneStream(input, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   PUT,
	Endpoint: "/streams/:id",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthRole(enums.AUTH_ROLE_USER),
		middlewares.BodyRequiredMiddleware,
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input, err := dtos.NewUpdateOneStreamInput().SetId(c.Params("id"))
			if err != nil {
				return common.ResponseError(c, err)
			}
			if err := c.BodyParser(input.Data); err != nil {
				e := exceptions.Err_BAD_REQUEST().SetMessage(err.Error())
				return common.ResponseError(c, e)
			}
			if err := input.Data.Validate(); err != nil {
				e := exceptions.Err_BAD_REQUEST().SetMessage(err.Error())
				return common.ResponseError(c, e)
			}
			res, err := services.GetStreamService().UpdateOneStream(input, req_ctx)
			if err != nil {
				e := exceptions.NewFromError(err)
				return common.ResponseError(c, e)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   DELETE,
	Endpoint: "/streams/:id",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthMiddleware,
		middlewares.AuthRole(enums.AUTH_ROLE_ADMIN),
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input, err := dtos.NewDeleteOneInput().SetId(c.Params("id"))
			if err != nil {
				return common.ResponseError(c, err)
			}
			res, err := services.GetStreamService().DeleteOneStream(input, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})
