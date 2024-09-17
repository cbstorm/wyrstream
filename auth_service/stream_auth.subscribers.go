package main

import (
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/nats_service"
)

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber(nats_service.AUTH_STREAM_CHECK_PUBLISH_KEY, func(m nats_service.IRequestMessage) (interface{}, error) {
	input := &dtos.CheckStreamKeyInput{}
	if err := m.JSONParse(input); err != nil {
		return nil, err
	}
	if err := GetAuthService().CheckStreamPublishKey(input); err != nil {
		return &dtos.CheckStreamKeyResponse{Message: err.Error()}, nil
	}
	return &dtos.CheckStreamKeyResponse{Ok: true}, nil
}))

var _ = nats_service.GetNATSService().AddSubcriber(nats_service.NewSubscriber(nats_service.AUTH_STREAM_CHECK_SUBSCRIBE_KEY, func(m nats_service.IRequestMessage) (interface{}, error) {
	input := &dtos.CheckStreamKeyInput{}
	if err := m.JSONParse(input); err != nil {
		return nil, err
	}
	if err := GetAuthService().CheckStreamSubscribeKey(input); err != nil {
		return &dtos.CheckStreamKeyResponse{Message: err.Error()}, nil
	}
	return &dtos.CheckStreamKeyResponse{Ok: true}, nil
}))
