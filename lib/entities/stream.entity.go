package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StreamEntity struct {
	BaseEntity      `bson:",inline"`
	PublisherId     primitive.ObjectID `bson:"publisher_id,omitempty" json:"publisher_id,omitempty"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	EnableRecord    bool               `bson:"enable_record" json:"enable_record"`
	StreamServerUrl string             `bson:"stream_server_url" json:"stream_server_url"`
	StreamId        string             `bson:"stream_id" json:"stream_id"`
	PublishKey      string             `bson:"publish_key" json:"-"`
	SubscribeKey    string             `bson:"subscribe_key" json:"-"`
	IsPublishing    bool               `bson:"is_publishing" json:"is_publishing"`
	PublishedAt     time.Time          `bson:"published_at,omitempty" json:"published_at,omitempty"`
	StoppedAt       time.Time          `bson:"stopped_at,omitempty" json:"stopped_at,omitempty"`
	HLSUrl          string             `bson:"hls_url" json:"hls_url"`
	ThumbnailUrl    string             `bson:"thumbnail_url" json:"thumbnail_url"`
	GuidanceCommand string             `bson:"guidance_command" json:"guidance_command"`

	PublishStreamUrl string              `bson:"-" json:"stream_url"`
	StreamLogs       *[]*StreamLogEntity `bson:"-" json:"stream_logs"`
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

func (e *StreamEntity) MakeGuidanceCommand() *StreamEntity {
	e.GuidanceCommand = fmt.Sprintf("ffmpeg -i <YOUR_INPUT> -c:v libx264 -b:v 2M -maxrate:v 2M -bufsize:v 1M -preset ultrafast -f mpegts \"%s?streamid=publish:/live/%s?key=%s\"", e.StreamServerUrl, e.StreamId, e.PublishKey)
	return e
}
