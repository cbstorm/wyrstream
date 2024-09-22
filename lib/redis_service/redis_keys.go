package redis_service

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

const (
	REDIS_KEY_AUTH_TOKEN           = "AUTH_TOKEN"
	REDIS_KEY_AUTH_TOKEN_BLACKLIST = "AUTH_TOKEN_BACKLIST"
)
