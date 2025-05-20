package wraperror

type ApiDisplayableError struct {
	httpStatus int
	message    interface{}
	err        error
}

func NewApiDisplayableError(
	httpStatus int,
	message interface{},
	err error,
) *ApiDisplayableError {
	return &ApiDisplayableError{
		httpStatus: httpStatus,
		message:    message,
		err:        err,
	}
}

func (e *ApiDisplayableError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	if message, ok := e.message.(string); ok {
		return message
	}

	return "Unknown error"
}

func (e *ApiDisplayableError) Unwrap() error {
	return e.err
}

func (e *ApiDisplayableError) Message() interface{} {
	return e.message
}

func (e *ApiDisplayableError) HttpStatus() int {
	return e.httpStatus
}
