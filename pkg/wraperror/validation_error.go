package wraperror

import "net/http"

type validationError struct {
	messages map[string]interface{}
	Err      error
}

func ValidationError(
	messages map[string]interface{},
	err error,
) *validationError {
	return &validationError{
		messages: messages,
		Err: NewAPIError(
			http.StatusBadRequest,
			messages,
			err,
		),
	}
}

func (e *validationError) Error() string {
	return e.Err.Error()
}

func (e *validationError) Unwrap() error {
	return e.Err
}
