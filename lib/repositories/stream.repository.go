package repositories

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/entities"
)

var stream_repository *StreamRepository
var stream_repository_sync sync.Once

func GetStreamRepository() *StreamRepository {
	if stream_repository == nil {
		stream_repository_sync.Do(func() {
			db := database.GetDatabase()
			stream_collection := db.DB().Collection("streams")
			stream_repository = &StreamRepository{
				CRUDRepository[*entities.StreamEntity]{
					collection: stream_collection,
				},
			}
		})
	}
	return stream_repository
}

type StreamRepository struct {
	CRUDRepository[*entities.StreamEntity]
}
