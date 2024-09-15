package middlewares

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/gofiber/fiber/v2"
)

var BodyRequiredMiddleware HttpMiddleware = func(c *fiber.Ctx) error {
	if len(c.BodyRaw()) <= 0 {
		e := exceptions.Err_BAD_REQUEST().SetMessage("Request body is required")
		return common.ResponseError(c, e)
	}
	return c.Next()
}
