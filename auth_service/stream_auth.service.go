package main

import (
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/repositories"
)

func (s *AuthService) CheckStreamPublishKey(input *dtos.CheckStreamKeyInput) error {
	stream := entities.NewStreamEntity()
	err, is_not_found := repositories.GetStreamRepository().FindOneByStreamIdAndPublishKey(input.StreamId, input.Key, stream)
	if err != nil {
		return err
	}
	if is_not_found {
		return exceptions.Err_BAD_REQUEST().SetMessage("stream not found")
	}
	return nil
}

func (s *AuthService) CheckStreamSubscribeKey(input *dtos.CheckStreamKeyInput) error {
	stream := entities.NewStreamEntity()
	err, is_not_found := repositories.GetStreamRepository().FindOneByStreamIdAndSubscribeKey(input.StreamId, input.Key, stream)
	if err != nil {
		return err
	}
	if is_not_found {
		return exceptions.Err_BAD_REQUEST().SetMessage("stream not found")
	}
	return nil
}
