package main

import (
	"github.com/cbstorm/wyrstream/lib/dtos"
	nats_service "github.com/cbstorm/wyrstream/lib/nats_service"
)

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber("auth.user.login", func(m nats_service.IRequestMessage) (interface{}, error) {
	user_login_input := &dtos.UserLoginInput{}
	if err := m.JSONParse(user_login_input); err != nil {
		return nil, err
	}
	return GetAuthService().UserLogin(user_login_input)
}))

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber("auth.user.create_account", func(m nats_service.IRequestMessage) (interface{}, error) {
	user_create_account_input := &dtos.UserCreateAccountInput{}
	if err := m.JSONParse(user_create_account_input); err != nil {
		return nil, err
	}
	return GetAuthService().UserCreateAccount(user_create_account_input)
}))

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber("auth.user.get_me", func(m nats_service.IRequestMessage) (interface{}, error) {
	user_get_me_input := &dtos.UserGetMeInput{}
	if err := m.JSONParse(user_get_me_input); err != nil {
		return nil, err
	}
	return GetAuthService().UserGetMe(user_get_me_input)
}))
