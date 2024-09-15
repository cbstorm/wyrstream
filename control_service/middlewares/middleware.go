package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

type HttpMiddleware func(*fiber.Ctx) error
