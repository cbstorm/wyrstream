package services

import (
	"sync"

	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/repositories"
)

var user_service *UserService
var user_service_sync sync.Once

func GetUserService() *UserService {
	if user_service == nil {
		user_service_sync.Do(func() {
			user_service = NewUserService()
		})
	}
	return user_service
}

type UserService struct {
	user_repository *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		user_repository: repositories.GetUserRepository(),
	}
}

func (svc *UserService) FetchUsers(fetchArgs *dtos.FetchArgs, reqCtx *common.RequestContext) (*repositories.FetchOutput[*entities.UserEntity], error) {
	users := make([]*entities.UserEntity, 0)
	return svc.user_repository.Fetch(fetchArgs, &users)
}

func (svc *UserService) GetOneUser(input *dtos.GetOneInput, reqCtx *common.RequestContext) (*entities.UserEntity, error) {
	user := entities.NewUserEntity()
	err, is_not_found := svc.user_repository.FindOneById(input.Id, user)
	if err != nil {
		return nil, err
	}
	if is_not_found {
		return nil, exceptions.Err_RESOURCE_NOT_FOUND()
	}
	return user, nil
}

func (svc *UserService) DeleteOneUser(input *dtos.DeleteOneInput, reqCtx *common.RequestContext) (*entities.UserEntity, error) {
	user := entities.NewUserEntity()
	return user, nil
}
