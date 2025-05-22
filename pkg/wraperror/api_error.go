package wraperror

type apiError struct {
	httpStatus int
	message    interface{}
	err        error
}

func APIError(
	httpStatus int,
	message interface{},
	err error,
) *apiError {
	return &apiError{
		httpStatus: httpStatus,
		message:    message,
		err:        err,
	}
}

func (e *apiError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	if message, ok := e.message.(string); ok {
		return message
	}

	return "Unknown error"
}

func (e *apiError) Unwrap() error {
	return e.err
}

func (e *apiError) Message() interface{} {
	return e.message
}

func (e *apiError) HttpStatus() int {
	return e.httpStatus
}
