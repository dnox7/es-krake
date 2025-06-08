package wraperror

type APIError struct {
	httpStatus int
	message    interface{}
	err        error
}

func NewAPIError(
	httpStatus int,
	message interface{},
	err error,
) *APIError {
	return &APIError{
		httpStatus: httpStatus,
		message:    message,
		err:        err,
	}
}

func (e *APIError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	if message, ok := e.message.(string); ok {
		return message
	}

	return "Unknown error"
}

func (e *APIError) Unwrap() error {
	return e.err
}

func (e *APIError) Message() interface{} {
	return e.message
}

func (e *APIError) HttpStatus() int {
	return e.httpStatus
}
