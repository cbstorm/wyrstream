package helpers

import (
	"fmt"
	"time"

	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/enums"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/redis_service"
	"github.com/golang-jwt/jwt"
)

func WithClaimOption(key string, value string) AuthTokenClaimOptionFunc {
	return func(claims *jwt.MapClaims) {
		(*claims)[key] = value
	}
}

type AuthTokenClaimOptionFunc func(*jwt.MapClaims)

type TokenHelper struct {
	person entities.IPerson
	role   enums.EAuthRole
}

func NewTokenHelper(person entities.IPerson, role enums.EAuthRole) *TokenHelper {
	return &TokenHelper{person: person, role: role}
}

func (th *TokenHelper) CreateAuthToken(token_type enums.ETokenTypes, opts ...AuthTokenClaimOptionFunc) (string, error) {
	redisInstance := redis_service.GetRedisService()
	redisTokenKey := th.hashTokenKey(token_type)
	if token_type == enums.TOKEN_TYPE_ACCESS_TOKEN {
		var existedToken string
		err, isNotFound := redisInstance.Get(redisTokenKey, &existedToken)
		if !isNotFound && err == nil {
			return existedToken, nil
		}
	}
	cfg := configs.GetConfig()
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	var tokenExp time.Duration
	if token_type == enums.TOKEN_TYPE_ACCESS_TOKEN {
		tokenExp = time.Minute * 10
		claims["tokenType"] = enums.TOKEN_TYPE_ACCESS_TOKEN
	}
	if token_type == enums.TOKEN_TYPE_REFRESH_TOKEN {
		tokenExp = time.Hour * 24
		claims["tokenType"] = enums.TOKEN_TYPE_REFRESH_TOKEN
	}
	claims["iss"] = cfg.APP_NAME
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(tokenExp).Unix()
	claims["_id"] = th.person.GetId().Hex()
	claims["name"] = th.person.GetName()
	claims["email"] = th.person.GetEmail()
	claims["role"] = th.role

	for _, optFunc := range opts {
		optFunc(&claims)
	}

	tokenString, err := token.SignedString([]byte(cfg.JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}
	err = redisInstance.SetWithTtl(redisTokenKey, tokenString, tokenExp)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyAuthToken(tokenString string) (map[string]interface{}, error) {
	cfg := configs.GetConfig()
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT_SECRET_KEY), nil
	})

	if err != nil {
		return nil, exceptions.Err_UNAUTHORIZED().SetMessage("UNAUTHORIZED")
	}

	if !token.Valid {
		return nil, exceptions.Err_UNAUTHORIZED().SetMessage("UNAUTHORIZED")
	}
	return claims, nil
}

func (th *TokenHelper) hashTokenKey(tokenType enums.ETokenTypes) string {
	return fmt.Sprintf("%s:%s", th.person.GetId().Hex(), tokenType)
}
