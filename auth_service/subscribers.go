package main

import (
	"github.com/cbstorm/wyrstream/lib/dtos"
	nats_service "github.com/cbstorm/wyrstream/lib/nats_service"
)

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber("auth.login", func(m nats_service.IMessage) (interface{}, error) {
	user_login_input := &dtos.UserLoginInput{}
	if err := m.JSONParse(user_login_input); err != nil {
		return nil, err
	}
	return GetAuthService().UserLogin(user_login_input)
}))
