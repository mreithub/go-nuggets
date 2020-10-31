package statuserror

import (
	"errors"
	"fmt"
	"net/http"
)

// New -- creates a new error (wraps errors.New())
func New(code int, text string) StatusError {
	return statusError{
		err:  errors.New(text),
		code: code,
	}
}

func Errorf(code int, format string, args ...interface{}) StatusError {
	return statusError{
		err:  fmt.Errorf(format, args...),
		code: code,
	}
}

// FromError -- returns a StatusError for the given error
//
// if err is a StatusError, it'll be returned as-is
// Otherwise, we'll try to extract the error code (using the Code() method if available) or default to http.StatusInternalServerError
// Then we'll call From() with the given status code
func FromError(err error) StatusError {
	if serr, ok := err.(StatusError); ok {
		// already a StatusError -> simply return it
		return serr
	}

	var code = http.StatusInternalServerError

	if codeErr, ok := err.(interface {
		Code() int
	}); ok {
		code = codeErr.Code()
	}

	return From(code, err)
}

// From -- returns a StatusError for the given error (with Code() set to code)
func From(code int, err error) StatusError {
	// TODO if err is a StatusError, use its data
	if serr, ok := err.(statusError); ok { // TODO make this work with the interface, not just the struct (right now we'll loose Unwrap() information)
		return statusError{
			code: code, err: serr.err,
			data: serr.data,
		}
	}

	return statusError{
		code: code,
		err:  err,
	}
}
