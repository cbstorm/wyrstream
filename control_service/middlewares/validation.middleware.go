package middlewares

import (
	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/lib/exceptions"
)

var BodyRequiredMiddleware HttpMiddleware = func(c common.IHttpContext) error {
	if len(c.BodyRaw()) <= 0 {
		e := exceptions.Err_BAD_REQUEST().SetMessage("Request body is required")
		return common.ResponseError(c, e)
	}
	return c.Next()
}
