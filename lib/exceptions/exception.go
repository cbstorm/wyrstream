package exceptions

import "errors"

type Exception struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Name    string `json:"name"`
}

func NewException(name string) *Exception {
	return &Exception{
		Name: name,
	}
}

func (e *Exception) Error() string {
	return e.Message
}

func (e *Exception) SetStatus(code int) *Exception {
	e.Status = code
	return e
}

func (e *Exception) SetName(name string) *Exception {
	e.Name = name
	return e
}

func (e *Exception) SetMessage(message string) *Exception {
	e.Message = message
	return e
}

func (e *Exception) GetStatus() int {
	return e.Status
}

func NewFromError(e error) *Exception {
	var exception *Exception
	if errors.As(e, &exception) {
		return exception
	}
	return Err_INTERNAL_SERVER_ERROR().SetMessage(e.Error())
}
