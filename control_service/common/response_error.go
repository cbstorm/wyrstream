package common

import (
	"github.com/cbstorm/wyrstream/lib/exceptions"
)

func ResponseError(c IHttpContext, err error) error {
	exception := exceptions.NewFromError(err)
	return c.Status(exception.GetStatus()).JSON(exception)
}

func ResponseOK(c IHttpContext, data interface{}) error {
	return c.Status(200).JSON(data)
}
