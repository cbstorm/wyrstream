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
		"stream_id":     stream_id,
		"is_publishing": false,
		"publish_key":   publish_key,
	}, out, opts...)
}

func (r *StreamRepository) FindOneByStreamIdAndSubscribeKey(stream_id, subscribe_key string, out *entities.StreamEntity, opts ...CURDOptionFunc) (error, bool) {
	return r.FindOne(map[string]interface{}{
		"stream_id":     stream_id,
		"subscribe_key": subscribe_key,
		"is_publishing": true,
	}, out, opts...)
}

func (r *StreamRepository) UpdatePublishStartByStreamId(stream_id, stream_server_url, hls_url, thumbnail_url string, out *entities.StreamEntity, opts ...CURDOptionFunc) error {
	return r.WithTransaction(func(ctx mongo.SessionContext) error {
		// Update stream
		if err := r.UpdateOne(
			map[string]interface{}{
				"stream_id": stream_id,
			},
			map[string]interface{}{
				"is_publishing":     true,
				"published_at":      time.Now().UTC(),
				"hls_url":           hls_url,
				"thumbnail_url":     thumbnail_url,
				"stream_server_url": stream_server_url,
			},
			out,
			WithContext(ctx),
		); err != nil {
			return err
		}
		// Insert stream log
		if err := GetStreamLogRepository().InsertOne(entities.NewStreamLogEntity().SetStreamId(out.Id).SetStartLog(), WithContext(ctx)); err != nil {
			return err
		}
		return nil
	}, opts...)
}

func (r *StreamRepository) UpdatePublishStopByStreamId(stream_id string, hls_segment_count uint, out *entities.StreamEntity, opts ...CURDOptionFunc) error {
	return r.WithTransaction(func(ctx mongo.SessionContext) error {
		// Update stream
		if err := r.UpdateOne(
			map[string]interface{}{
				"stream_id": stream_id,
			},
			map[string]interface{}{
				"is_publishing": false,
				"stopped_at":    time.Now().UTC(),
				"updatedAt":     time.Now().UTC(),
			},
			out,
			WithIncr(map[string]interface{}{"hls_segment_count": hls_segment_count}),
			WithContext(ctx)); err != nil {
			return err
		}
		// Insert stream log
		if err := GetStreamLogRepository().InsertOne(entities.NewStreamLogEntity().SetStreamId(out.Id).SetStopLog(), WithContext(ctx)); err != nil {
			return err
		}
		return nil
	}, opts...)
}

func (r *StreamRepository) UpdateStreamToClosed(stream *entities.StreamEntity, opts ...CURDOptionFunc) error {
	return r.WithTransaction(func(ctx mongo.SessionContext) error {
		if err := r.UpdateOneById(stream.Id, map[string]interface{}{
			"is_closed":     true,
			"is_publishing": false,
			"updatedAt":     time.Now().UTC(),
		}, stream, WithContext(ctx)); err != nil {
			return err
		}
		if err := GetStreamLogRepository().InsertOne(entities.NewStreamLogEntity().SetStreamId(stream.Id).SetClosedLog(), WithContext(ctx)); err != nil {
			return err
		}
		return nil
	}, opts...)
}
