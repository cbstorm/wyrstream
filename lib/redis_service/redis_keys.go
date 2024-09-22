package redis_service

import "strings"

type RedisKey string

func (k RedisKey) Concat(s string) RedisKey {
	return RedisKey(string(k) + ":" + s)
}

func (k RedisKey) ConcatKey(s RedisKey) RedisKey {
	return RedisKey(string(k) + ":" + string(s))
}

func (k RedisKey) String() string {
	return string(k)
}
func (k RedisKey) PatternString() string {
	return k.String() + "*"
}
func (k RedisKey) TrimPrefix(s RedisKey) RedisKey {
	return RedisKey(strings.TrimPrefix(k.String(), s.String()+":"))
}

const (
	REDIS_KEY_AUTH_TOKEN           RedisKey = "AUTH_TOKEN"
	REDIS_KEY_AUTH_TOKEN_BLACKLIST RedisKey = "AUTH_TOKEN_BACKLIST"
)

const (
	REDIS_KEY_STREAM_SERVER_HEALTH RedisKey = "STREAM_SERVER_HEALTH"
)
