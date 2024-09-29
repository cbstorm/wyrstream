package dtos

import (
	"time"

	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserLoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (i *UserLoginInput) Validate() error {
	return utils.NewValidator(i).Validate()
}

type UserLoginResponse struct {
	Id           primitive.ObjectID `json:"_id"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
}

type UserCreateAccountInput struct {
	Name     string `json:"name" validate:"required,min_length=6,max_length=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min_length=8"`
}

func (i *UserCreateAccountInput) Validate() error {
	return utils.NewValidator(i).Validate()
}

type UserCreateAccountReponse struct {
	Id        primitive.ObjectID `json:"_id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

type UserGetMeInput struct {
	UserId primitive.ObjectID `json:"user_id"`
}

type UserRefreshTokenInput struct {
	AccesssToken string `json:"access_token,omitempty" validate:"required"`
	RefreshToken string `json:"refresh_token,omitempty" validate:"required"`
}

func (i *UserRefreshTokenInput) Validate() error {
	return utils.NewValidator(i).Validate()
}

type UserRefreshTokenOutput struct {
	NewToken        string `json:"new_token,omitempty"`
	NewRefreshToken string `json:"new_refresh_token,omitempty"`
}

type CheckStreamKeyInput struct {
	StreamServer string `json:"stream_server"`
	StreamId     string `json:"stream_id"`
	Key          string `json:"key"`
}

type CheckStreamKeyResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}
