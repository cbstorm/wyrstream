package helpers

import (
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/repositories"
	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StreamsHelper struct {
	streams *[]*entities.StreamEntity
}

func NewStreamsHelper(streams *[]*entities.StreamEntity) *StreamsHelper {
	return &StreamsHelper{
		streams: streams,
	}
}

func (h *StreamsHelper) ResolveStreamLogs() error {
	if len(*h.streams) == 0 {
		return nil
	}
	stream_obj_ids := utils.Map(h.streams, func(e *entities.StreamEntity, i int) primitive.ObjectID {
		return e.Id
	})
	stream_logs := make([]*entities.StreamLogEntity, 0)
	if err := repositories.GetStreamLogRepository().Find(map[string]interface{}{
		"stream_obj_id": map[string]interface{}{
			"$in": stream_obj_ids,
		},
	}, &stream_logs); err != nil {
		return err
	}
	stream_logs_group_by_stream_id := utils.GroupBy(&stream_logs, func(e *entities.StreamLogEntity) string {
		return e.StreamObjId.Hex()
	})
	for _, v := range *h.streams {
		v.StreamLogs = (*stream_logs_group_by_stream_id)[v.Id.Hex()]
	}
	return nil
}
