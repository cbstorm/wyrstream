package repositories

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/entities"
)

var vod_repository *VodRepository
var vod_repository_sync sync.Once

func GetVodRepository() *VodRepository {
	if vod_repository == nil {
		vod_repository_sync.Do(func() {
			db := database.GetDatabase()
			vod_collection := db.DB().Collection("vods")
			vod_repository = &VodRepository{
				CRUDRepository[*entities.VodEntity]{
					collection: vod_collection,
				},
			}
		})
	}
	return vod_repository
}

type VodRepository struct {
	CRUDRepository[*entities.VodEntity]
}
