package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type VodEntity struct {
	BaseEntity   `bson:",inline"`
	OwnerId      primitive.ObjectID `bson:"owner_id,omitempty" json:"omitempty"`
	Title        string             `bson:"title" json:"title"`
	Description  string             `bson:"description" json:"description"`
	HLSUrl       string             `bson:"hls_url" json:"hls_url"`
	ThumbnailUrl string             `bson:"thumbnail_url" json:"thumbnail_url"`
}

func NewVodEntity() *VodEntity {
	vod := &VodEntity{}
	vod.New()
	return vod
}
