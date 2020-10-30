package statuserror

import "errors"

// New -- creates a new error (wraps errors.New())
func New(code int, text string) StatusError {
	return statusError{
		err:  errors.New(text),
		code: code,
	}
}
