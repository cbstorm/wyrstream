package entities

import (
	"crypto/rand"
	"encoding/base64"
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
	HlsUrl       string             `bson:"hls_url" json:"hls_url"`
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
	buf := make([]byte, 256)
	rand, _ := rand.Read(buf)
	rand_hex := fmt.Sprintf("%x", rand)
	data := []byte(e.PublisherId.Hex() + e.StreamId + rand_hex)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)
	e.PublishKey = strings.ToUpper(string(dst))
	return e
}

func (e *StreamEntity) GenerateSubscribeKey() *StreamEntity {
	buf := make([]byte, 512)
	rand, _ := rand.Read(buf)
	rand_hex := fmt.Sprintf("%x", rand)
	data := []byte(e.Id.Hex() + e.PublisherId.Hex() + e.StreamId + rand_hex)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)
	e.SubscribeKey = strings.ToUpper(string(dst))
	return e
}
