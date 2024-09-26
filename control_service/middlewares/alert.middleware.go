package middlewares

import (
	"github.com/cbstorm/wyrstream/lib/alert_service"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/gofiber/fiber/v2"
)

var Alert HttpMiddleware = func(c *fiber.Ctx) error {
	p := &dtos.AlertPayload{
		Method:  c.Method(),
		Url:     c.OriginalURL(),
		Payload: string(c.BodyRaw()),
	}
	alert_service.GetAlertService().Alert(p)
	return c.Next()
}
