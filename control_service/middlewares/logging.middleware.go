package middlewares

import (
	"fmt"

	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/gofiber/fiber/v2"
)

func GetRequestContext(c *fiber.Ctx) *common.RequestContext {
	req_ctx := c.UserContext()
	return req_ctx.Value(common.ReqContextKey{}).(*common.RequestContext)
}

var LoggingMiddleware HttpMiddleware = func(c *fiber.Ctx) error {
	reqCtx := AssertRequestContext(c)
	reqCtx.GetLogger().Info(fmt.Sprintf("%s %s %v", c.Method(), c.IP(), c.OriginalURL()))
	return c.Next()
}
