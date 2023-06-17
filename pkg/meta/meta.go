package meta

import "net/http"

var (
	MetaErrorNotFound = &MetaError{
		HttpStatus: http.StatusNotFound,
	}
	MetaErrorInternalServer = &MetaError{
		HttpStatus: http.StatusInternalServerError,
	}
	MetaErrorUnauthorized = &MetaError{
		HttpStatus: http.StatusUnauthorized,
	}
	MetaErrorForbidden = &MetaError{
		HttpStatus: http.StatusForbidden,
	}
	MetaErrorBadRequest = &MetaError{
		HttpStatus: http.StatusBadRequest,
	}
	Error = &MetaError{}
)

type MetaError struct {
	HttpStatus int    `json:"-"`
	Code       int    `json:"-"`
	Message    string `json:"message"`
}

func (m MetaError) AppendError(code int, err error) *MetaError {
	if code == 0 {
		panic("error code is empty, meta error have to set code to append function")
	}
	if err == nil {
		panic("error is empty but is instructed to generate an error")
	}

	m.Code = code
	m.Message = err.Error()
	return &m
}

func (m MetaError) AppendMessage(code int, msg string) *MetaError {
	if code == 0 {
		panic("error code is empty, meta error have to set code to append function")
	}

	m.Code = code
	m.Message = msg

	return &m
}

func (m *MetaError) Error() string {
	return m.Message
}

func NewError(httpStatus int) MetaError {
	return MetaError{
		HttpStatus: httpStatus,
	}
}

func NewOK() int {
	return 0
}

func IsError(err error) (*MetaError, bool) {
	if err, ok := err.(*MetaError); ok {
		return err, true
	}

	return nil, false
}
