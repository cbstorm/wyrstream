package enums

type ETokenTypes string

const (
	TOKEN_TYPE_ACCESS_TOKEN  ETokenTypes = "ACCESS_TOKEN"
	TOKEN_TYPE_REFRESH_TOKEN ETokenTypes = "REFRESH_TOKEN"
	TOKEN_TYPE_DO_TEST_TOKEN ETokenTypes = "DO_TEST_TOKEN"
)

type EAuthRole string

const (
	AUTH_ROLE_USER  EAuthRole = "USER"
	AUTH_ROLE_ADMIN EAuthRole = "ADMIN"
)

func (e EAuthRole) String() string {
	return string(e)
}
