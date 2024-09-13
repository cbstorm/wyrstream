package exceptions

var (
	Err_UNAUTHORIZED = func() *Exception {
		return NewException("UNAUTHORIZED").SetStatus(401)
	}
	Err_FORBIDEN = func() *Exception {
		return NewException("FORBIDEN").SetStatus(403)
	}
	Err_INVALID_TOKEN = func() *Exception {
		return NewException("INVALID_TOKEN").SetStatus(400)
	}
	Err_EXISTED_EMAIL = func() *Exception {
		return NewException("EXISTED_EMAIL").SetStatus(400)
	}
	Err_EXISTED_EMAIL_OR_USERNAME = func() *Exception {
		return NewException("EXISTED_EMAIL_OR_USERNAME").SetStatus(400)
	}
	Err_USER_NOT_FOUND = func() *Exception {
		return NewException("USER_NOT_FOUND").SetStatus(404)
	}
	Err_PASSWORD_INCORRECT = func() *Exception {
		return NewException("PASSWORD_INCORRECT").SetStatus(400)
	}
)
