package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IEntity interface {
	NewId() primitive.ObjectID
	GetId() primitive.ObjectID
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreatedAt() time.Time
	SetUpdatedAt() time.Time
	SetTime() time.Time
}

type BaseEntity struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

func (e *BaseEntity) NewId() primitive.ObjectID {
	e.Id = primitive.NewObjectID()
	return e.Id
}

func (e *BaseEntity) GetId() primitive.ObjectID {
	return e.Id
}

func (e *BaseEntity) GetCreatedAt() time.Time {
	return e.CreatedAt
}

func (e *BaseEntity) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}

func (e *BaseEntity) SetCreatedAt() time.Time {
	e.CreatedAt = time.Now().UTC()
	return e.CreatedAt
}

func (e *BaseEntity) SetUpdatedAt() time.Time {
	e.UpdatedAt = time.Now().UTC()
	return e.UpdatedAt
}

func (e *BaseEntity) SetTime() time.Time {
	t := time.Now().UTC()
	e.CreatedAt = t
	e.UpdatedAt = t
	return t
}
