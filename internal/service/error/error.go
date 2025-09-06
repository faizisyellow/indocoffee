package errorService

import (
	"errors"
)

func New(err, internal error) error {
	return &ErrorService{err, internal}
}

type ErrorService struct {
	E        error
	Internal error
}

func (err *ErrorService) Error() string {
	return err.E.Error()
}

func (err *ErrorService) InternalError() string {
	return err.Internal.Error()
}

func GetError(err error) *ErrorService {
	if err == nil {
		return &ErrorService{errors.New("encounter internal error"), errors.New("errors is nill")}
	}

	errValue, ok := err.(*ErrorService)
	if !ok {
		return &ErrorService{err, errors.New("errors not from service")}
	}

	return errValue
}
