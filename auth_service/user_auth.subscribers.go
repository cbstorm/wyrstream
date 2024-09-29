package main

import (
	"github.com/cbstorm/wyrstream/lib/dtos"
	nats_service "github.com/cbstorm/wyrstream/lib/nats_service"
)

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber("auth.user.login", func(m nats_service.IRequestMessage) (interface{}, error) {
	input := &dtos.UserLoginInput{}
	if err := m.JSONParse(input); err != nil {
		return nil, err
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return GetAuthService().UserLogin(input)
}))

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber("auth.user.create_account", func(m nats_service.IRequestMessage) (interface{}, error) {
	input := &dtos.UserCreateAccountInput{}
	if err := m.JSONParse(input); err != nil {
		return nil, err
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return GetAuthService().UserCreateAccount(input)
}))

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber("auth.user.get_me", func(m nats_service.IRequestMessage) (interface{}, error) {
	input := &dtos.UserGetMeInput{}
	if err := m.JSONParse(input); err != nil {
		return nil, err
	}
	return GetAuthService().UserGetMe(input)
}))

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber("auth.user.refresh_token", func(m nats_service.IRequestMessage) (interface{}, error) {
	input := &dtos.UserRefreshTokenInput{}
	if err := m.JSONParse(input); err != nil {
		return nil, err
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return GetAuthService().UserRefreshToken(input)
}))
