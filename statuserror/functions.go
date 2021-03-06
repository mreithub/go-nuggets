package statuserror

import (
	"errors"
	"fmt"
	"net/http"
)

// New -- creates a new error (wraps errors.New())
func New(code int, text string) Error {
	return statusError{
		err:  errors.New(text),
		code: code,
	}
}

// Errorf -- creates and returns an error (like fmt.Errorf() but also holds an integer status code)
func Errorf(code int, format string, args ...interface{}) Error {
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
func FromError(err error) Error {
	if serr, ok := err.(Error); ok {
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
func From(code int, err error) Error {
	// TODO make this work with the interface, not just the struct (right now we'll loose Unwrap() information)
	if serr, ok := err.(statusErrorWithData); ok {
		serr.code = code
		return serr
	}
	if serr, ok := err.(statusError); ok {
		serr.code = code
		return serr
	}

	return statusError{
		code: code,
		err:  err,
	}
}

// GetCode -- returns the result of an error's Code() method (if present)
//
// if not found, http.StatusInternalServerError will be returned
func GetCode(err error) int {
	if codeErr, ok := err.(interface {
		Code() int
	}); ok {
		return codeErr.Code()
	}

	return http.StatusInternalServerError
}

// GetData -- will call Data() if present on the given error value
func GetData(err error) map[string]interface{} {
	if dataErr, ok := err.(interface {
		Data() map[string]interface{}
	}); ok {
		return dataErr.Data()
	}

	return nil
}
