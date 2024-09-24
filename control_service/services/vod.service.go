package services

import (
	"sync"

	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/repositories"
)

var vod_service *VodService
var vod_service_sync sync.Once

func GetVodService() *VodService {
	if vod_service == nil {
		vod_service_sync.Do(func() {
			vod_service = NewVodService()
		})
	}
	return vod_service
}

type VodService struct {
	vod_repository *repositories.VodRepository
}

func NewVodService() *VodService {
	return &VodService{
		vod_repository: repositories.GetVodRepository(),
	}
}

func (svc *VodService) FetchVods(fetchArgs *dtos.FetchArgs, reqCtx *common.RequestContext) (*repositories.FetchOutput[*entities.VodEntity], error) {
	vods := make([]*entities.VodEntity, 0)
	return svc.vod_repository.Fetch(fetchArgs, &vods)
}

func (svc *VodService) GetOneVod(input *dtos.GetOneInput, reqCtx *common.RequestContext) (*entities.VodEntity, error) {
	vod := entities.NewVodEntity()
	err, is_not_found := svc.vod_repository.FindOneById(input.Id, vod)
	if err != nil {
		return nil, err
	}
	if is_not_found {
		return nil, exceptions.Err_RESOURCE_NOT_FOUND()
	}
	return vod, nil
}

func (svc *VodService) UpdateOneVod(input *dtos.UpdateOneVodInput, reqCtx *common.RequestContext) (*entities.VodEntity, error) {
	vod := entities.NewVodEntity()
	err, is_not_found := svc.vod_repository.FindOne(map[string]interface{}{
		"_id":      input.Id,
		"owner_id": reqCtx.GetObjId(),
	}, vod)
	if err != nil {
		return nil, err
	}
	if is_not_found {
		return nil, exceptions.Err_RESOURCE_NOT_FOUND().SetMessage("vod not found")
	}
	vod.Title = input.Data.Title
	vod.Description = input.Data.Description
	vod.SetUpdatedAt()
	if err := svc.vod_repository.UpdateOneById(vod.Id, vod, vod); err != nil {
		return nil, err
	}
	return vod, nil
}

func (svc *VodService) DeleteOneVod(input *dtos.DeleteOneInput, reqCtx *common.RequestContext) (*entities.VodEntity, error) {
	vod := entities.NewVodEntity()
	return vod, nil
}
