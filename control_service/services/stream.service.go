package services

import (
	"sync"

	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/control_service/helpers"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/repositories"
)

var stream_service *StreamService
var stream_service_sync sync.Once

func GetStreamService() *StreamService {
	if stream_service == nil {
		stream_service_sync.Do(func() {
			stream_service = NewStreamService()
		})
	}
	return stream_service
}

type StreamService struct {
	stream_repository *repositories.StreamRepository
}

func NewStreamService() *StreamService {
	return &StreamService{
		stream_repository: repositories.GetStreamRepository(),
	}
}

func (svc *StreamService) FetchStreams(fetchArgs *dtos.FetchArgs, reqCtx *common.RequestContext) (*repositories.FetchOutput[*entities.StreamEntity], error) {
	streams := make([]*entities.StreamEntity, 0)
	res, err := svc.stream_repository.Fetch(fetchArgs, &streams)
	helper := helpers.NewStreamsHelper(res.Result)
	if fetchArgs.IsIncludes("stream_logs") {
		if err := helper.ResolveStreamLogs(); err != nil {
			return nil, err
		}
	}
	return res, err
}

func (svc *StreamService) GetOneStream(input *dtos.GetOneInput, reqCtx *common.RequestContext) (*entities.StreamEntity, error) {
	stream := entities.NewStreamEntity()
	err, is_not_found := svc.stream_repository.FindOneById(input.Id, stream)
	if err != nil {
		return nil, err
	}
	if is_not_found {
		return nil, exceptions.Err_RESOURCE_NOT_FOUND()
	}
	return stream, nil
}

func (svc *StreamService) CreateOneStream(input *dtos.CreateOneStreamInput, reqCtx *common.RequestContext) (*entities.StreamEntity, error) {
	stream := entities.NewStreamEntity()
	stream.PublisherId = reqCtx.GetObjId()
	stream.Title = input.Title
	stream.Description = input.Description
	stream.EnableRecord = input.EnableRecord
	stream.GenerateStreamId().GeneratePublishKey().GenerateSubscribeKey().MakeGuidanceCommand()
	if err := repositories.GetStreamRepository().InsertOne(stream); err != nil {
		return nil, err
	}
	return stream, nil
}

func (svc *StreamService) UpdateOneStream(input *dtos.UpdateOneStreamInput, reqCtx *common.RequestContext) (*entities.StreamEntity, error) {
	stream := entities.NewStreamEntity()
	stream.SetTime()
	return stream, nil
}

func (svc *StreamService) DeleteOneStream(input *dtos.DeleteOneInput, reqCtx *common.RequestContext) (*entities.StreamEntity, error) {
	stream := entities.NewStreamEntity()
	return stream, nil
}
