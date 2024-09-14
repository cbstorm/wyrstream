package middlewares

import (
	"fmt"

	"github.com/cbstorm/wyrstream/control_service/common"
)

func GetRequestContext(c common.IHttpContext) *common.RequestContext {
	req_ctx := c.UserContext()
	return req_ctx.Value(common.ReqContextKey{}).(*common.RequestContext)
}

var LoggingMiddleware HttpMiddleware = func(c common.IHttpContext) error {
	reqCtx := common.GetRequestContext(c)
	if reqCtx == nil {
		reqCtx = common.NewRequestContext().ParseHeader(c.GetReqHeaders()).GetIPForwardedFor(c.GetReqHeaders()).SetIP(c.IP()).SetMethod(c.Method()).SetPath(c.OriginalURL())
		c.Locals(common.ReqContextKey{}, reqCtx)
	}
	reqCtx.GetLogger().Info(fmt.Sprintf("%s %s %v", c.Method(), c.IP(), c.OriginalURL()))
	return c.Next()
}
