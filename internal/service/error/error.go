package errorService

func New(text, internal string) error {
	return &ErrorService{text, internal}
}

type ErrorService struct {
	s string
	i string
}

func (e *ErrorService) Error() string {
	return e.s
}

func (e *ErrorService) InternalError() string {
	return e.i
}
