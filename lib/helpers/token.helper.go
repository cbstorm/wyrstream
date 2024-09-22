package helpers

import (
	"time"

	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/enums"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"github.com/cbstorm/wyrstream/lib/redis_service"
	"github.com/cbstorm/wyrstream/lib/utils"
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
	redis_instance := redis_service.GetRedisService()
	redis_token_key := th.hashTokenKey(token_type)
	if token_type == enums.TOKEN_TYPE_ACCESS_TOKEN {
		var existed_token string
		err, is_not_found := redis_instance.Get(redis_token_key, &existed_token)
		if !is_not_found && err == nil {
			return existed_token, nil
		}
	}
	cfg := configs.GetConfig()
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	var token_exp time.Duration
	if token_type == enums.TOKEN_TYPE_ACCESS_TOKEN {
		token_exp = time.Minute * 10
		claims["tokenType"] = enums.TOKEN_TYPE_ACCESS_TOKEN
	}
	if token_type == enums.TOKEN_TYPE_REFRESH_TOKEN {
		token_exp = time.Hour * 24
		claims["tokenType"] = enums.TOKEN_TYPE_REFRESH_TOKEN
	}
	claims["iss"] = cfg.APP_NAME
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(token_exp).Unix()
	claims["_id"] = th.person.GetId().Hex()
	claims["name"] = th.person.GetName()
	claims["email"] = th.person.GetEmail()
	claims["role"] = th.role

	for _, optFunc := range opts {
		optFunc(&claims)
	}

	token_string, err := token.SignedString([]byte(cfg.JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}
	err = redis_instance.SetWithTtl(redis_token_key, token_string, token_exp)
	if err != nil {
		return "", err
	}
	return token_string, nil
}

func VerifyAuthToken(token_string string) (map[string]interface{}, error) {
	cfg := configs.GetConfig()
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)

	token, err := jwt.ParseWithClaims(token_string, &claims, func(token *jwt.Token) (interface{}, error) {
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

func SaveBlacklist(token string, token_type enums.ETokenTypes) error {
	token_sum := utils.SHA512Sum(token)
	ttl := utils.TernaryOp(token_type == enums.TOKEN_TYPE_ACCESS_TOKEN, time.Minute*10, time.Hour*24)
	return redis_service.GetRedisService().SetWithTtl(redis_service.RedisKey(redis_service.REDIS_KEY_AUTH_TOKEN_BLACKLIST).Concat(token_sum), 1, ttl)
}

func CheckBlacklist(token string) error {
	token_sum := utils.SHA512Sum(token)
	var val string
	err, is_not_found := redis_service.GetRedisService().Get(redis_service.RedisKey(redis_service.REDIS_KEY_AUTH_TOKEN_BLACKLIST).Concat(token_sum), val)
	if err != nil {
		return err
	}
	if is_not_found {
		return nil
	}
	return exceptions.Err_UNAUTHORIZED().SetMessage("UNAUTHORIZED")
}

func (th *TokenHelper) hashTokenKey(token_type enums.ETokenTypes) redis_service.RedisKey {
	return redis_service.RedisKey(redis_service.REDIS_KEY_AUTH_TOKEN).Concat(th.person.GetId().Hex()).Concat(string(token_type))
}
