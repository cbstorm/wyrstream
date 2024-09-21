package main

import (
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
	hls_url := BuildHLSUrl(input.StreamId)
	thumbnail_url := BuildThumbnailUrl(input.StreamId)
	if err := utils.AssertDir(BuildHLSStreamDir(input.StreamId) + "/" + THUMBNAIL_DIR + "/"); err != nil {
		return err
	}
	stream := entities.NewStreamEntity()
	if err := repositories.GetStreamRepository().UpdatePublishStartByStreamId(input.StreamId, hls_url, thumbnail_url, stream); err != nil {
		return err
	}
	stream_url := BuildStreamURL(input.StreamServer, input.StreamServerApp, input.StreamId, stream.SubscribeKey)
	hls_cmd := NewProcessHLSCommand(input.StreamId).SetStartNumber(s.getStartFileNumber(input.StreamId)).SetInput(stream_url)
	GetProcessHLSCommandStore().Add(hls_cmd)
	go hls_cmd.Run()
	thumbnail_cmd := NewProcessThumbnailCommand(input.StreamId)
	GetProcessThumbnailCommandStore().Add(thumbnail_cmd)
	go thumbnail_cmd.Start()
	return nil
}

func (s *HLSService) ProcessStop(input *dtos.HLSPublishStopInput) error {
	if hls_cmd := GetProcessHLSCommandStore().Get(input.StreamId); hls_cmd != nil {
		hls_cmd.Cancel()
		GetProcessHLSCommandStore().Remove(input.StreamId)
	}
	if thumbnail_cmd := GetProcessThumbnailCommandStore().Get(input.StreamId); thumbnail_cmd != nil {
		thumbnail_cmd.Cancel()
		GetProcessThumbnailCommandStore().Remove(input.StreamId)
	}
	stream := entities.NewStreamEntity()
	if err := repositories.GetStreamRepository().UpdatePublishStopByStreamId(input.StreamId, stream); err != nil {
		return err
	}
	return nil
}

func (s *HLSService) getStartFileNumber(stream_id string) uint {
	files := GetListSegmentFilesByStreamId(stream_id)
	return uint(len(*files) + 1)
}
