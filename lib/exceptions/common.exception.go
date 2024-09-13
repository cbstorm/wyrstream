package exceptions

var (
	Err_RESOURCE_NOT_FOUND = func() *Exception {
		return NewException("RESOURCE_NOT_FOUND").SetStatus(404)
	}
	Err_BAD_REQUEST = func() *Exception {
		return NewException("BAD_REQUEST").SetStatus(400)
	}
	Err_INTERNAL_SERVER_ERROR = func() *Exception {
		return NewException("INTERNAL_SERVER_ERROR").SetStatus(500)
	}
)
