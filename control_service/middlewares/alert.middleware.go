package middlewares

import (
	"github.com/cbstorm/wyrstream/control_service/services"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/gofiber/fiber/v2"
)

var Alert HttpMiddleware = func(c *fiber.Ctx) error {
	reqCtx := AssertRequestContext(c)
	p := &dtos.AlertPayload{
		Method:  c.Method(),
		Url:     c.OriginalURL(),
		Payload: string(c.BodyRaw()),
	}
	go func() {
		if err := services.GetAlertService().SendAlert(p); err != nil {
			reqCtx.GetLogger().Error("Could not send alert with err: %v", err)
		}
	}()
	return c.Next()
}
