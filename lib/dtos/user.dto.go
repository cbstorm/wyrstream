package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserReponse struct {
	Id        primitive.ObjectID `json:"_id,omitempty"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty"`
}
