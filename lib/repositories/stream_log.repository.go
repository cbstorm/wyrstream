package repositories

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/entities"
)

var stream_log_repository *StreamLogRepository
var stream_log_repository_sync sync.Once

func GetStreamLogRepository() *StreamLogRepository {
	if stream_log_repository == nil {
		stream_log_repository_sync.Do(func() {
			db := database.GetDatabase()
			stream_log_collection := db.DB().Collection("stream_logs")
			stream_log_repository = &StreamLogRepository{
				CRUDRepository[*entities.StreamLogEntity]{
					collection: stream_log_collection,
				},
			}
		})
	}
	return stream_log_repository
}

type StreamLogRepository struct {
	CRUDRepository[*entities.StreamLogEntity]
}
