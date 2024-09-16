package middlewares

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/gofiber/fiber/v2"
)

type HttpMiddleware func(*fiber.Ctx) error

func AssertRequestContext(c *fiber.Ctx) *common.RequestContext {
	reqCtx := common.GetRequestContext(c)
	if reqCtx == nil {
		reqCtx = common.NewRequestContext().ParseHeader(c.GetReqHeaders()).GetIPForwardedFor(c.GetReqHeaders()).SetIP(c.IP()).SetMethod(c.Method()).SetPath(c.OriginalURL()).SetHttpContext(c.Context())
		c.Locals(common.ReqContextKey{}, reqCtx)
	}
	return reqCtx
}
