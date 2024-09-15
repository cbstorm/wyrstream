package middlewares

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/gofiber/fiber/v2"
)

var AuthMiddleware HttpMiddleware = func(c *fiber.Ctx) error {
	reqCtx := common.GetRequestContext(c)
	if reqCtx == nil {
		reqCtx = common.NewRequestContext().ParseHeader(c.GetReqHeaders()).GetIPForwardedFor(c.GetReqHeaders()).SetIP(c.IP()).SetMethod(c.Method()).SetPath(c.OriginalURL())
		c.Locals(common.ReqContextKey{}, reqCtx)
	}
	err := reqCtx.Auth()
	if err != nil {
		return common.ResponseError(c, err)
	}
	return c.Next()
}
