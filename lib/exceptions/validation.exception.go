package exceptions

var (
	Err_EMAIL_INVALID     = func() *Exception { return NewException("EMAIL_INVALID").SetStatus(400) }
	Err_USERNAME_INVALID  = func() *Exception { return NewException("USERNAME_INVALID").SetStatus(400) }
	Err_PASSWORD_INVALID  = func() *Exception { return NewException("PASSWORD_INVALID").SetStatus(400) }
	Err_USER_TYPE_INVALID = func() *Exception { return NewException("USER_TYPE_INVALID").SetStatus(400) }
)

var (
	Err_ID_INVALID = func() *Exception { return NewException("ID_INVALID").SetStatus(400) }
)
