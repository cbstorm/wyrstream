package middlewares

import (
	"github.com/cbstorm/wyrstream/control_service/common"
)

type HttpMiddleware func(common.IHttpContext) error
