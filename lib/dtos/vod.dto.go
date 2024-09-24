package dtos

import (
	"fmt"

	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateOneVodInput struct {
	Name string `json:"name,omitempty"`
}

func NewCreateOneVodInput() *CreateOneVodInput {
	return &CreateOneVodInput{}
}

type UpdateOneVodInput struct {
	Id   primitive.ObjectID `json:"_id,omitempty"`
	Data *UpdateOneVodData  `json:"data,omitempty"`
}

func (d *UpdateOneVodInput) SetId(id string) (*UpdateOneVodInput, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, exceptions.Err_BAD_REQUEST().SetMessage(fmt.Sprintf("id invalid with %s", id))
	}
	d.Id = objId
	return d, nil
}

type UpdateOneVodData struct {
	Title       string `json:"title,omitempty" validate:"required,min_length=6,max_length=50"`
	Description string `json:"description,omitempty" validate:"required,min_length=6,max_length=255"`
}

func (i *UpdateOneVodData) Validate() error {
	return utils.NewValidator(i).Validate()
}

func NewUpdateOneVodInput() *UpdateOneVodInput {
	return &UpdateOneVodInput{
		Data: &UpdateOneVodData{},
	}
}
