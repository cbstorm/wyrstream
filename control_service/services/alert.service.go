package services

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/dtos"
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
}

func NewAlertService() *AlertService {
	return &AlertService{
		nats_service: nats_service.GetNATSService(),
	}
}

func (svc *AlertService) SendAlert(payload *dtos.AlertPayload) error {
	if _, err := nats_service.GetNATSService().Request(nats_service.ALERT, payload); err != nil {
		return err
	}
	return nil
}
