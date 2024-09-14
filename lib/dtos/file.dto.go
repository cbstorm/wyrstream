package dtos

import (
	"fmt"
	"strings"

	"github.com/cbstorm/wyrstream/lib/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileInput struct {
	Name     string
	Size     int64
	Bytes    []byte
	Feat     string
	MimeType string
}

type CreateOneFileInput struct {
	Files *[]*FileInput
}

func (f *FileInput) GetExt() string {
	sp := strings.Split(f.Name, ".")
	return sp[len(sp)-1]
}

func (f *FileInput) GetFileType() string {
	sp := strings.Split(f.MimeType, "/")
	return sp[0]
}

func NewCreateOneFileInput() *CreateOneFileInput {
	return &CreateOneFileInput{}
}

type UpdateOneFileInput struct {
	Id   primitive.ObjectID `json:"_id,omitempty"`
	Data *UpdateOneFileData `json:"data,omitempty"`
}

func (d *UpdateOneFileInput) SetId(id string) (*UpdateOneFileInput, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, exceptions.Err_BAD_REQUEST().SetMessage(fmt.Sprintf("id invalid with %s", id))
	}
	d.Id = objId
	return d, nil
}

type UpdateOneFileData struct {
	Name string `json:"name,omitempty"`
}

func NewUpdateOneFileInput() *UpdateOneFileInput {
	return &UpdateOneFileInput{}
}

type UploadS3Result struct {
	FileName string
	Url      string
	Err      error
	FileSize int64
	MimeType string
	Feat     string
	FileType string
	SizeType string
}
