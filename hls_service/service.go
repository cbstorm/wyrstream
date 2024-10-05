package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/cbstorm/wyrstream/lib/alert_service"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/minio_service"
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
	if err := s.stream_repository.UpdatePublishStartByStreamId(input.StreamId, input.StreamServer, hls_url, thumbnail_url, stream); err != nil {
		return err
	}
	stream_url := BuildStreamURL(input.StreamServer, input.StreamServerApp, input.StreamId, stream.SubscribeKey)
	hls_cmd := NewProcessHLSCommand(input.StreamId, stream.EnableRecord).SetStartNumber(stream.HLSSegmentCount + 1).SetInput(stream_url)
	GetProcessHLSCommandStore().Add(hls_cmd)
	go hls_cmd.Run()
	thumbnail_cmd := NewProcessThumbnailCommand(input.StreamId)
	GetProcessThumbnailCommandStore().Add(thumbnail_cmd)
	go thumbnail_cmd.Start()
	return nil
}

func (s *HLSService) ProcessStop(input *dtos.HLSPublishStopInput) error {
	logg := s.logger.Child(input.StreamId)
	if hls_cmd := GetProcessHLSCommandStore().Get(input.StreamId); hls_cmd != nil {
		hls_cmd.Cancel()
		GetProcessHLSCommandStore().Remove(input.StreamId)
	}
	if thumbnail_cmd := GetProcessThumbnailCommandStore().Get(input.StreamId); thumbnail_cmd != nil {
		thumbnail_cmd.Cancel()
		GetProcessThumbnailCommandStore().Remove(input.StreamId)
	}
	hls_segment_count := uint(len(GetListSegmentFilesByStreamId(input.StreamId)))
	stream := entities.NewStreamEntity()
	if err := s.stream_repository.UpdatePublishStopByStreamId(input.StreamId, hls_segment_count, stream); err != nil {
		return err
	}
	if stream.EnableRecord {
		if err := s.putSegmentsToStorage(input.StreamId); err != nil {
			logg.Error("Could not put segments to storage due to an error: %v", err)
			return err
		}
	}
	thumbnail_url, err := s.putThumbnailToStorage(input.StreamId)
	if err != nil {
		logg.Error("Could not put thumbnail to storage due to an error: %v", err)
	}
	if err := s.stream_repository.UpdateOne(map[string]interface{}{"stream_id": input.StreamId}, map[string]interface{}{
		"thumbnail_url": thumbnail_url,
		"ready_for_vod": true,
	}, stream); err != nil {
		s.logger.Error("Could not update the thumbnail url due to an error: %v", err)
	}

	if err := s.cleanStreamDir(input.StreamId); err != nil {
		logg.Error("Could not clean the directory due to an error: %v", err)
		return err
	}
	go alert_service.GetAlertService().Alert(&dtos.StreamStopAlert{StreamId: stream.StreamId, Title: stream.Title})
	return nil
}

func (s *HLSService) putSegmentsToStorage(stream_id string) error {
	segments := GetListSegmentFilesByStreamId(stream_id)
	if len(segments) == 0 {
		return fmt.Errorf("hls segments is empty")
	}
	seg_objs := utils.Map(segments, func(a string, b int) minio_service.MinIOFObject {
		return &minio_service.HLSSegmentObject{
			StreamId: stream_id,
			Name:     a,
			Path:     fmt.Sprintf("%s/%s", BuildHLSStreamDir(stream_id), a),
		}
	})
	res := minio_service.GetMinioService().FPutObjects(seg_objs)
	for _, v := range *res {
		if v.Error != nil {
			s.logger.Error("Could not fput object %s due to an  error: %v", v.ObjectName, v.Error)
		}
	}
	return nil
}

func (s *HLSService) putThumbnailToStorage(stream_id string) (string, error) {
	p := BuildThumbnailFilePath(stream_id)
	obj := &minio_service.StreamThumbnailObject{
		StreamId: stream_id,
		Path:     p,
	}
	if !obj.EnsurePath() {
		return "", fmt.Errorf("thumbnail at %s doesn't exist", p)
	}
	return minio_service.GetMinioService().FPutObject(obj)
}

func (s *HLSService) cleanStreamDir(stream_id string) error {
	p := BuildHLSStreamDir(stream_id)
	if err := os.RemoveAll(p); err != nil {
		s.logger.Error("Could not remove dir %s due to an error: %v", p, err)
		return err
	}
	return nil
}
