package main

import (
	"fmt"
	"sync"

	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/repositories"
	"github.com/cbstorm/wyrstream/lib/utils"
)

var hls_service *HLSService
var hls_service_sync sync.Once

func GetHLSService() *HLSService {
	if hls_service == nil {
		hls_service_sync.Do(func() {
			hls_service = &HLSService{
				logger:            logger.NewLogger("HLS_SERVICE"),
				stream_repository: repositories.GetStreamRepository(),
			}
		})
	}
	return hls_service
}

type HLSService struct {
	logger            *logger.Logger
	stream_repository *repositories.StreamRepository
}

func (s *HLSService) ProcessStart(input *dtos.HLSPublishStartInput) error {
	m3u8_file := "playlist.m3u8"
	hls_url := BuildHLSUrl(input.StreamId, m3u8_file)
	if err := utils.AssertDir(fmt.Sprintf("public/%s/%s", input.StreamId, m3u8_file)); err != nil {
		return err
	}
	stream := entities.NewStreamEntity()
	if err := repositories.GetStreamRepository().UpdatePublishStartByStreamId(input.StreamId, hls_url, stream); err != nil {
		return err
	}
	stream_url := fmt.Sprintf("%s?streamid=%s%s?key=%s", input.StreamServer, input.StreamServerApp, input.StreamId, stream.SubscribeKey)
	c := NewProcessHLSCommand(input.StreamId).SetStartNumber(1).SetInput(stream_url).SetOutput(m3u8_file)
	// c.Print()
	go c.Run()
	return nil
}
