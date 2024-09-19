package entities

type StreamLogEntity struct {
	BaseEntity `bson:",inline"`
	Log        string `bson:"log" json:"log"`
}

func NewStreamLogEntity() *StreamLogEntity {
	stream_log := &StreamLogEntity{}
	stream_log.New()
	return stream_log
}
