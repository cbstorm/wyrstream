package main

import "github.com/cbstorm/wyrstream/lib/nats_service"

var _ = nats_service.GetNATSService().AddSubcriber(
	nats_service.NewSubscriber(nats_service.HLS_PUBLISH_START,
		func(im nats_service.IRequestMessage) (interface{}, error) {

			//Listen sub
			nats_service.GetNATSService().AddSubcriber(
				nats_service.NewSubscriber(nats_service.HLS_PUBLISH_STOP,
					func(im nats_service.IRequestMessage) (interface{}, error) {
						return nil, nil
					},
				))

			return nil, nil
		},
	))
