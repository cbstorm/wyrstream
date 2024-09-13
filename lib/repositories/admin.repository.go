package repositories

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/entities"
)

var admin_repository *AdminRepository
var admin_repository_sync sync.Once

func GetAdminRepository() *AdminRepository {
	if admin_repository == nil {
		admin_repository_sync.Do(func() {
			db := database.GetDatabase()
			admin_collection := db.DB().Collection("admins")
			admin_repository = &AdminRepository{
				CRUDRepository[*entities.AdminEntity]{
					collection: admin_collection,
				},
			}
		})
	}
	return admin_repository
}

type AdminRepository struct {
	CRUDRepository[*entities.AdminEntity]
}
