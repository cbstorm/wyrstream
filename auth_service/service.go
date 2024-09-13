package main

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/logger"
)

var instance *AuthService
var instance_sync sync.Once

func GetAuthService() *AuthService {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &AuthService{
				logger: logger.NewLogger("AUTH_SERVICE"),
			}
		})
	}
	return instance
}

func NewAuthService() {

}

type AuthService struct {
	logger *logger.Logger
}
