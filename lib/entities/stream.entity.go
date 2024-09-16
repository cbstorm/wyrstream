package entities

type StreamEntity struct {
	BaseEntity `bson:",inline"`
}

func NewStreamEntity() *StreamEntity {
	stream := &StreamEntity{}
	stream.NewId()
	stream.SetTime()
	return stream
}
