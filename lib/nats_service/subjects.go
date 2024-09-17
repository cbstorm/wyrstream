package nats_service

type NATS_Subject string

// User auth subjects
const (
	AUTH_USER_LOGIN          NATS_Subject = "auth.user.login"
	AUTH_USER_CREATE_ACCOUNT NATS_Subject = "auth.user.create_account"
	AUTH_USER_GET_ME         NATS_Subject = "auth.user.get_me"
	AUTH_USER_REFRESH_TOKEN  NATS_Subject = "auth.user.refresh_token"
)

const (
	AUTH_STREAM_CHECK_PUBLISH_KEY   NATS_Subject = "auth.stream.check_publish_key"
	AUTH_STREAM_CHECK_SUBSCRIBE_KEY NATS_Subject = "auth.stream.check_subscribe_key"
)

// stream.auth.check_publish_key
// stream.auth.check_subscribe_key
