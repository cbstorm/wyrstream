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
	Endpoint: "/vods",
	Handlers: []func(*fiber.Ctx) error{
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			fetchArgs := dtos.NewFetchArgs()
			fetchArgs.ParseQueries(c.Queries())
			res, err := services.GetVodService().FetchVods(fetchArgs, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/vods/my_vods",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthRole(enums.AUTH_ROLE_USER),
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			fetchArgs := dtos.NewFetchArgs()
			fetchArgs.ParseQueries(c.Queries())
			fetchArgs.SetFilter("owner_id", req_ctx.GetObjId())
			res, err := services.GetVodService().FetchVods(fetchArgs, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   GET,
	Endpoint: "/vods/:id",
	Handlers: []func(*fiber.Ctx) error{
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input, err := dtos.NewGetOneInput().SetId(c.Params("id"))
			if err != nil {
				return common.ResponseError(c, err)
			}
			res, err := services.GetVodService().GetOneVod(input, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})

var _ = GetHttpServer().FeedRoute(&HTTPRoute{
	Method:   PUT,
	Endpoint: "/vods/:id",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthMiddleware,
		middlewares.BodyRequiredMiddleware,
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input, err := dtos.NewUpdateOneVodInput().SetId(c.Params("id"))
			if err != nil {
				return common.ResponseError(c, err)
			}
			if err := c.BodyParser(input.Data); err != nil {
				return common.ResponseError(c, exceptions.Err_BAD_REQUEST().SetMessage(err.Error()))
			}
			if err := input.Data.Validate(); err != nil {
				return common.ResponseError(c, exceptions.Err_BAD_REQUEST().SetMessage(err.Error()))
			}
			res, err := services.GetVodService().UpdateOneVod(input, req_ctx)
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
	Endpoint: "/vods/:id",
	Handlers: []func(*fiber.Ctx) error{
		middlewares.AuthRole(enums.AUTH_ROLE_ADMIN),
		func(c *fiber.Ctx) error {
			req_ctx := common.GetRequestContext(c)
			input, err := dtos.NewDeleteOneInput().SetId(c.Params("id"))
			if err != nil {
				return common.ResponseError(c, err)
			}
			res, err := services.GetVodService().DeleteOneVod(input, req_ctx)
			if err != nil {
				return common.ResponseError(c, err)
			}
			return common.ResponseOK(c, res)
		},
	},
})
