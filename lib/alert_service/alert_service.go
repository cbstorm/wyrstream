package alert_service

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/nats_service"
)

var alert_service *AlertService
var alert_service_sync sync.Once

func GetAlertService() *AlertService {
	if alert_service == nil {
		alert_service_sync.Do(func() {
			alert_service = NewAlertService()
		})
	}
	return alert_service
}

type AlertService struct {
	nats_service *nats_service.NATS_Service
	logg         *logger.Logger
}

func NewAlertService() *AlertService {
	return &AlertService{
		nats_service: nats_service.GetNATSService(),
		logg:         logger.NewLogger("ALERT_SERVICE"),
	}
}

func (svc *AlertService) Alert(payload interface{}) {
	go func() {
		if _, err := svc.nats_service.Request(nats_service.ALERT, payload); err != nil {
			svc.logg.Error("Could not send alert with payload %v due to an error: %v", payload, err)
		}
	}()
}

func (svc *AlertService) StreamStopAlert(payload *dtos.StreamStopAlert) {
	svc.Alert(payload)
}
