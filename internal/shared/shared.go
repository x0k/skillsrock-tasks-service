package shared

import "errors"

var ErrNotFound = errors.New("not found")

type ServiceError struct {
	Expected bool
	Err      error
	Msg      string
}

func NewServiceError(err error, msg string) *ServiceError {
	return &ServiceError{
		Expected: true,
		Err:      err,
		Msg:      msg,
	}
}

func NewUnexpectedError(err error, msg string) *ServiceError {
	return &ServiceError{
		Expected: false,
		Err:      err,
		Msg:      msg,
	}
}
