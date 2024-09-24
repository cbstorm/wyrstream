package services

import (
	"log"
	"strconv"
	"sync"

	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/control_service/helpers"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/redis_service"
	"github.com/cbstorm/wyrstream/lib/repositories"
	"github.com/cbstorm/wyrstream/lib/utils"
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
	redis_service     *redis_service.RedisService
	stream_repository *repositories.StreamRepository
}

func NewStreamService() *StreamService {
	return &StreamService{
		stream_repository: repositories.GetStreamRepository(),
		redis_service:     redis_service.GetRedisService(),
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
	if reqCtx.PathIncludes("my_streams") {
		for _, v := range streams {
			v.MakeShownPublishKey().MakePublishStreamUrl()
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
	stream_server_url, err := svc.selectStreamServer()
	if err != nil {
		return nil, err
	}
	stream.StreamServerUrl = stream_server_url
	stream.GenerateStreamId().GeneratePublishKey().GenerateSubscribeKey().MakeShownPublishKey().MakePublishStreamUrl().MakeGuidanceCommand()
	if err := repositories.GetStreamRepository().InsertOne(stream); err != nil {
		return nil, err
	}
	return stream, nil
}

func (svc *StreamService) UpdateOneStream(input *dtos.UpdateOneStreamInput, reqCtx *common.RequestContext) (*entities.StreamEntity, error) {
	stream := entities.NewStreamEntity()
	err, is_not_found := svc.stream_repository.FindOne(map[string]interface{}{
		"_id":          input.Id,
		"publisher_id": reqCtx.GetObjId(),
	}, stream)
	if err != nil {
		return nil, err
	}
	if is_not_found {
		return nil, exceptions.Err_BAD_REQUEST().SetMessage("stream not found")
	}
	stream.Title = input.Data.Title
	stream.Description = input.Data.Description
	stream.EnableRecord = input.Data.EnableRecord
	stream.SetUpdatedAt()
	if err := svc.stream_repository.UpdateOneById(stream.Id, stream, stream); err != nil {
		return nil, err
	}
	return stream, nil
}

func (svc *StreamService) DeleteOneStream(input *dtos.DeleteOneInput, reqCtx *common.RequestContext) (*entities.StreamEntity, error) {
	stream := entities.NewStreamEntity()
	return stream, nil
}

func (svc *StreamService) selectStreamServer() (string, error) {
	server_keys := make([]redis_service.RedisKey, 0)
	if err := svc.redis_service.Keys(redis_service.REDIS_KEY_STREAM_SERVER_HEALTH, &server_keys); err != nil {
		return "", err
	}
	if len(server_keys) == 0 {
		return "", exceptions.Err_INTERNAL_SERVER_ERROR()
	}
	server_health := make([]string, 0)
	if err := svc.redis_service.MGet(server_keys, &server_health); err != nil {
		return "", err
	}
	min := 0
	min_idx := 0
	utils.ForEach(&server_health, func(a string, b int) {
		a_int, _ := strconv.ParseInt(a, 10, 32)
		log.Printf("idx:%d, a: %d, min: %d, s: %s", b, int(a_int), min, server_keys[b])
		if int(a_int) <= min {
			min = int(a_int)
			min_idx = b
		}
	})
	min_sv := server_keys[min_idx].TrimPrefix(redis_service.REDIS_KEY_STREAM_SERVER_HEALTH).String()

	return min_sv, nil
}
