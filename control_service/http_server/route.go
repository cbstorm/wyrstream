package http_server

import (
	"github.com/cbstorm/wyrstream/control_service/common"
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
	Handlers []func(common.IHttpContext) error
	Disable  bool
}
