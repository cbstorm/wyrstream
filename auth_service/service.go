package main

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/repositories"
)

var instance *AuthService
var instance_sync sync.Once

func GetAuthService() *AuthService {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &AuthService{
				logger:          logger.NewLogger("AUTH_SERVICE"),
				user_repository: repositories.GetUserRepository(),
			}
		})
	}
	return instance
}

type AuthService struct {
	logger          *logger.Logger
	user_repository *repositories.UserRepository
}
