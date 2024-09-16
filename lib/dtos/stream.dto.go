package dtos

import (
	"fmt"

	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateOneStreamInput struct {
	Title       string `json:"title,omitempty" validate:"required,min_length=6,max_length=50"`
	Description string `json:"description,omitempty" validate:"required,min_length=6,max_length=255"`
}

func (i *CreateOneStreamInput) Validate() error {
	return utils.NewValidator(i).Validate()
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
