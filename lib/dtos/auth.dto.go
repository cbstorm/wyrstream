package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id           primitive.ObjectID `json:"_id"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
}

type CreateAccountInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAccountReponse struct {
	Id        primitive.ObjectID `json:"_id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

type RefreshTokenInput struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshTokenOutput struct {
	NewToken        string `json:"new_token,omitempty"`
	NewRefreshToken string `json:"new_refresh_token,omitempty"`
}
