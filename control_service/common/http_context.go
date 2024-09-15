package common

import (
	"context"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

type IHttpContext interface {
	Locals(key interface{}, value ...interface{}) interface{}
	GetReqHeaders() map[string][]string
	MultipartForm() (*multipart.Form, error)
	Status(status int) *fiber.Ctx
	JSON(data interface{}, ctype ...string) error
	Queries() map[string]string
	Query(key string, defaultValue ...string) string
	UserContext() context.Context
	Method(override ...string) string
	Next() error
	IP() string
	OriginalURL() string
	BodyRaw() []byte
	BodyParser(out interface{}) error
}
