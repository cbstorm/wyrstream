package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UserRefreshTokenOutput struct {
	NewToken        string `json:"new_token,omitempty"`
	NewRefreshToken string `json:"new_refresh_token,omitempty"`
}
