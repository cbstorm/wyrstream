package middlewares

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/lib/enums"
	"github.com/gofiber/fiber/v2"
)

var AuthMiddleware HttpMiddleware = func(c *fiber.Ctx) error {
	reqCtx := AssertRequestContext(c)
	if err := reqCtx.Auth(); err != nil {
		return common.ResponseError(c, err)
	}
	return c.Next()
}

var AuthRole = func(role enums.EAuthRole) HttpMiddleware {
	return func(c *fiber.Ctx) error {
		reqCtx := AssertRequestContext(c)
		if role == enums.AUTH_ROLE_ADMIN {
			if err := reqCtx.AuthAdmin(); err != nil {
				return common.ResponseError(c, err)
			}
		}
		if role == enums.AUTH_ROLE_USER {
			if err := reqCtx.AuthUser(); err != nil {
				return common.ResponseError(c, err)
			}
		}
		return c.Next()
	}
}
