package http_server

import (
	"github.com/gofiber/fiber/v2"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

type HTTPRoute struct {
	Method   Method
	Endpoint string
	Handlers []func(*fiber.Ctx) error
	Disable  bool
}
