package wraperror

import "net/http"

type ValidationError struct {
	messages map[string]interface{}
	Err      error
}

func NewValidationError(
	messages map[string]interface{},
	err error,
) *ValidationError {
	return &ValidationError{
		messages: messages,
		Err: NewApiDisplayableError(
			http.StatusBadRequest,
			messages,
			err,
		),
	}
}

func (e *ValidationError) Error() string {
	return e.Err.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}
