package main

import (
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/nats_service"
)

var _ = nats_service.GetNATSService().AddSubcriber(
	nats_service.NewSubscriber(nats_service.HLS_PUBLISH_START,
		func(m nats_service.IRequestMessage) (interface{}, error) {
			start_input := &dtos.HLSPublishStartInput{}
			if err := m.JSONParse(start_input); err != nil {
				return nil, err
			}
			//Listen stop subscriber
			listen_publish_stop := nats_service.NewSubscriber(nats_service.HLS_PUBLISH_STOP.Concat(start_input.StreamId),
				func(msg nats_service.IRequestMessage) (interface{}, error) {
					stop_input := &dtos.HLSPublishStopInput{}
					if err := msg.JSONParse(stop_input); err != nil {
						return nil, err
					}
					// Process stop command
					if err := GetHLSService().ProcessStop(stop_input); err != nil {
						return nil, err
					}
					return nil, nil
				},
			)
			nats_service.GetNATSService().AddSubcriber(listen_publish_stop)
			nats_service.GetNATSService().Start(listen_publish_stop.GetId())
			// Process start command
			if err := GetHLSService().ProcessStart(start_input); err != nil {
				return nil, err
			}
			return &dtos.HLSPublishStartResponse{Ok: true}, nil
		},
	))
