package services

import (
	"sync"

	"github.com/cbstorm/wyrstream/control_service/common"
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
	return svc.stream_repository.Fetch(fetchArgs, &streams)
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
	Stream := entities.NewStreamEntity()
	return Stream, nil
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
