package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type StreamEntity struct {
	BaseEntity  `bson:",inline"`
	PublisherId primitive.ObjectID `bson:"publisher_id,omitempty" json:"publisher_id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
}

func NewStreamEntity() *StreamEntity {
	stream := &StreamEntity{}
	stream.New()
	return stream
}
