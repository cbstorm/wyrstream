package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StreamEntity struct {
	BaseEntity   `bson:",inline"`
	PublisherId  primitive.ObjectID `bson:"publisher_id,omitempty" json:"publisher_id,omitempty"`
	Title        string             `bson:"title" json:"title"`
	Description  string             `bson:"description" json:"description"`
	StreamId     string             `bson:"stream_id" json:"stream_id"`
	PublishKey   string             `bson:"publish_key" json:"-"`
	SubscribeKey string             `bson:"subscribe_key" json:"-"`
	IsPublishing bool               `bson:"is_publishing" json:"is_publishing"`
	PublishedAt  time.Time          `bson:"published_at" json:"published_at"`
	StoppedAt    time.Time          `bson:"stopped_at" json:"stopped_at"`
	HLSUrl       string             `bson:"hls_url" json:"hls_url"`
}

func NewStreamEntity() *StreamEntity {
	stream := &StreamEntity{}
	stream.New()
	return stream
}

func (e *StreamEntity) GenerateStreamId() *StreamEntity {
	timestamp_hex := fmt.Sprintf("%x", time.Now().UTC().Unix())
	counter_hex := fmt.Sprintf("%x", utils.GetCounter().Increase("stream"))
	e.StreamId = strings.ToUpper("STR" + timestamp_hex + counter_hex)
	return e
}

func (e *StreamEntity) GeneratePublishKey() *StreamEntity {
	e.PublishKey = utils.StringRand(30)
	return e
}

func (e *StreamEntity) GenerateSubscribeKey() *StreamEntity {
	e.SubscribeKey = utils.StringRand(30)
	return e
}
