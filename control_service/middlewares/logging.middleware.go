package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var LoggingMiddleware HttpMiddleware = func(c *fiber.Ctx) error {
	reqCtx := AssertRequestContext(c)
	reqCtx.GetLogger().Info(fmt.Sprintf("%s %s %v", c.Method(), c.IP(), c.OriginalURL()))
	return c.Next()
}
