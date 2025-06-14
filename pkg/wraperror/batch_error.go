package wraperror

type BatchError struct {
	Code  string
	Retry bool
	err   error
}

func NewBatchError(code string, err error) *BatchError {
	return &BatchError{
		Code:  code,
		Retry: false,
		err:   err,
	}
}

func NewBatchErrorAndRetry(code string, err error) *BatchError {
	return &BatchError{
		Code:  code,
		Retry: true,
		err:   err,
	}
}

func (err *BatchError) Error() string {
	return err.err.Error()
}

func (err *BatchError) Unwrap() error {
	return err.err
}
