package main

import (
	"sync"

	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/enums"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/helpers"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/repositories"
	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var instance *AuthService
var instance_sync sync.Once

func GetAuthService() *AuthService {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &AuthService{
				logger:          logger.NewLogger("AUTH_SERVICE"),
				user_repository: repositories.GetUserRepository(),
			}
		})
	}
	return instance
}

type AuthService struct {
	logger          *logger.Logger
	user_repository *repositories.UserRepository
}

func (svc *AuthService) UserLogin(input *dtos.UserLoginInput) (*dtos.UserLoginResponse, error) {
	if !utils.IsValidEmailAddress(input.Email) {
		return nil, exceptions.Err_EMAIL_INVALID().SetMessage("Your email address invalid")
	}
	user := entities.NewUserEntity()
	err, is_not_found := svc.user_repository.FindOneByEmail(input.Email, user)
	if err != nil {
		return nil, err
	}
	if is_not_found {
		return nil, exceptions.Err_RESOURCE_NOT_FOUND().SetMessage("Your email did not exist or password does not match")
	}
	if match := utils.BcryptMatch(user.Password, input.Password); !match {
		return nil, exceptions.Err_BAD_REQUEST().SetMessage("Your email did not exist or password does not match")
	}
	token_helper := helpers.NewTokenHelper(user, enums.AUTH_ROLE_USER)
	access_token, err := token_helper.CreateAuthToken(enums.TOKEN_TYPE_ACCESS_TOKEN)
	if err != nil {
		return nil, err
	}
	refresh_token, err := token_helper.CreateAuthToken(enums.TOKEN_TYPE_REFRESH_TOKEN, helpers.WithClaimOption("token", utils.MD5Sum(access_token)))
	if err != nil {
		return nil, err
	}
	out := &dtos.UserLoginResponse{
		Id:           user.Id,
		Name:         user.Name,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
	return out, nil
}

func (svc *AuthService) UserCreateAccount(input *dtos.UserCreateAccountInput) (*dtos.UserCreateAccountReponse, error) {
	existed_user := entities.NewUserEntity()
	err, is_not_found := svc.user_repository.FindOneByEmail(input.Email, existed_user)
	if err != nil {
		return nil, err
	}
	if !is_not_found {
		return nil, exceptions.Err_EXISTED_EMAIL().SetMessage("Your email already existed")
	}
	user := entities.NewUserEntity()
	user.Name = input.Name
	user.Email = input.Email
	user.Password = input.Password
	user_helper := helpers.NewUserHelper(user)
	if err := user_helper.HashPassword(); err != nil {
		return nil, err
	}
	if err := svc.user_repository.InsertOne(user); err != nil {
		return nil, err
	}
	out := &dtos.UserCreateAccountReponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return out, nil
}

func (svc *AuthService) UserGetMe(input *dtos.UserGetMeInput) (*dtos.UserReponse, error) {
	user := entities.NewUserEntity()
	err, is_not_found := svc.user_repository.FindOneById(input.UserId, user)
	if err != nil {
		return nil, err
	}
	if is_not_found {
		return nil, exceptions.Err_RESOURCE_NOT_FOUND().SetMessage("User not found")
	}
	out := &dtos.UserReponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return out, nil
}

func (svc *AuthService) UserRefreshToken(input *dtos.UserRefreshTokenInput) (*dtos.UserRefreshTokenOutput, error) {
	payload, err := helpers.VerifyAuthToken(input.RefreshToken)
	if err != nil {
		return nil, err
	}
	if payload["token"] != utils.MD5Sum(input.Token) {
		return nil, exceptions.Err_UNAUTHORIZED().SetMessage("UNAUTHORIZED")
	}
	id := payload["_id"].(string)
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, exceptions.Err_UNAUTHORIZED().SetMessage("UNAUTHORIZED")
	}
	user := entities.NewUserEntity()
	err, is_not_found := svc.user_repository.FindOneById(objId, user)
	if err != nil {
		return nil, err
	}
	if is_not_found {
		return nil, exceptions.Err_UNAUTHORIZED().SetMessage("UNAUTHORIZED")
	}
	token_helper := helpers.NewTokenHelper(user, enums.AUTH_ROLE_USER)
	token, err := token_helper.CreateAuthToken(enums.TOKEN_TYPE_ACCESS_TOKEN)
	if err != nil {
		return nil, err
	}
	refresh_token, err := token_helper.CreateAuthToken(enums.TOKEN_TYPE_REFRESH_TOKEN, helpers.WithClaimOption("token", utils.MD5Sum(token)))
	if err != nil {
		return nil, err
	}
	output := &dtos.UserRefreshTokenOutput{
		NewToken:        token,
		NewRefreshToken: refresh_token,
	}
	return output, nil
}
