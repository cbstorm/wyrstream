package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StreamLogEntity struct {
	BaseEntity  `bson:",inline"`
	StreamObjId primitive.ObjectID `bson:"stream_obj_id" json:"stream_obj_id"`
	Log         string             `bson:"log" json:"log"`
}

func NewStreamLogEntity() *StreamLogEntity {
	stream_log := &StreamLogEntity{}
	stream_log.New()
	return stream_log
}

func (e *StreamLogEntity) SetStreamId(stream_obj_id primitive.ObjectID) *StreamLogEntity {
	e.StreamObjId = stream_obj_id
	return e
}

func (e *StreamLogEntity) SetStartLog() *StreamLogEntity {
	e.Log = "START"
	return e
}

func (e *StreamLogEntity) SetStopLog() *StreamLogEntity {
	e.Log = "STOP"
	return e
}

func (e *StreamLogEntity) SetClosedLog() *StreamLogEntity {
	e.Log = "Closed"
	return e
}
