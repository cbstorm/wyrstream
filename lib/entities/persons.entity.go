package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPerson interface {
	GetId() primitive.ObjectID
	GetName() string
	GetEmail() string
}

type PersonEntity struct {
	BaseEntity `bson:",inline"`
	Name       string `bson:"name" json:"name"`
	Email      string `bson:"email" json:"email"`
	Password   string `bson:"password" json:"-"`
}

func NewPersonEntity() *PersonEntity {
	person := &PersonEntity{}
	person.NewId()
	person.SetTime()
	return person
}

func (m *PersonEntity) GetId() primitive.ObjectID {
	return m.Id
}
func (m *PersonEntity) GetName() string {
	return m.Name
}
func (m *PersonEntity) GetEmail() string {
	return m.Email
}
