package dtos

import (
	"fmt"

	"github.com/cbstorm/wyrstream/lib/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateOneStreamInput struct {
	Name string `json:"name,omitempty"`
}

func NewCreateOneStreamInput() *CreateOneStreamInput {
	return &CreateOneStreamInput{}
}

type UpdateOneStreamInput struct {
	Id   primitive.ObjectID   `json:"_id,omitempty"`
	Data *UpdateOneStreamData `json:"data,omitempty"`
}

func (d *UpdateOneStreamInput) SetId(id string) (*UpdateOneStreamInput, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, exceptions.Err_BAD_REQUEST().SetMessage(fmt.Sprintf("id invalid with %s", id))
	}
	d.Id = objId
	return d, nil
}

type UpdateOneStreamData struct {
	Name string `json:"name,omitempty"`
}

func NewUpdateOneStreamInput() *UpdateOneStreamInput {
	return &UpdateOneStreamInput{}
}
