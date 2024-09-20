package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

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
