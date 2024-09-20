package repositories

import (
	"sync"
	"time"

	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/entities"
	"go.mongodb.org/mongo-driver/mongo"
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

func (r *StreamRepository) FindOneByStreamIdAndPublishKey(stream_id, publish_key string, out *entities.StreamEntity, opts ...CURDOptionFunc) (error, bool) {
	return r.FindOne(map[string]interface{}{
		"stream_id":   stream_id,
		"publish_key": publish_key,
	}, out, opts...)
}

func (r *StreamRepository) FindOneByStreamIdAndSubscribeKey(stream_id, subscribe_key string, out *entities.StreamEntity, opts ...CURDOptionFunc) (error, bool) {
	return r.FindOne(map[string]interface{}{
		"stream_id":     stream_id,
		"subscribe_key": subscribe_key,
	}, out, opts...)
}

func (r *StreamRepository) UpdatePublishStartByStreamId(stream_id, hls_url, thumbnail_url string, out *entities.StreamEntity, opts ...CURDOptionFunc) error {
	return r.WithTransaction(func(ctx mongo.SessionContext) error {
		// Update stream
		if err := r.UpdateOne(map[string]interface{}{
			"stream_id": stream_id,
		}, map[string]interface{}{
			"is_publishing": true,
			"published_at":  time.Now().UTC(),
			"hls_url":       hls_url,
			"thumbnail_url": thumbnail_url,
		}, out, WithContext(ctx)); err != nil {
			return err
		}
		// Insert stream log
		if err := GetStreamLogRepository().InsertOne(entities.NewStreamLogEntity().SetStreamId(out.Id).SetStartLog(), WithContext(ctx)); err != nil {
			return err
		}
		return nil
	}, opts...)
}

func (r *StreamRepository) UpdatePublishStopByStreamId(stream_id string, out *entities.StreamEntity, opts ...CURDOptionFunc) error {
	return r.WithTransaction(func(ctx mongo.SessionContext) error {
		// Update stream
		if err := r.UpdateOne(map[string]interface{}{
			"stream_id": stream_id,
		}, map[string]interface{}{
			"is_publishing": false,
			"stopped_at":    time.Now().UTC(),
		}, out, WithContext(ctx)); err != nil {
			return err
		}
		// Insert stream log
		if err := GetStreamLogRepository().InsertOne(entities.NewStreamLogEntity().SetStreamId(out.Id).SetStopLog(), WithContext(ctx)); err != nil {
			return err
		}
		return nil
	}, opts...)
}
