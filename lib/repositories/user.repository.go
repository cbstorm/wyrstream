package repositories

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/entities"
)

var user_repository *UserRepository
var user_repository_sync sync.Once

func GetUserRepository() *UserRepository {
	if user_repository == nil {
		user_repository_sync.Do(func() {
			db := database.GetDatabase()
			user_collection := db.DB().Collection("users")
			user_repository = &UserRepository{
				CRUDRepository[*entities.UserEntity]{
					collection: user_collection,
				},
			}
		})
	}
	return user_repository
}

type UserRepository struct {
	CRUDRepository[*entities.UserEntity]
}
