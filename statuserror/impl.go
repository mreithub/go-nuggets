package statuserror

import "errors"

type statusError struct {
	err  error
	code int
	data map[string]interface{}
}

// Code -- returns the status code associated with this StatusError
func (e statusError) Code() int     { return e.code }
func (e statusError) Error() string { return e.err.Error() }
func (e statusError) Unwrap() error { return errors.Unwrap(e.err) }

func (e statusError) WithData(data map[string]interface{}) StatusError {
	e.data = data
	return e
}

func (e statusError) WithField(key string, value interface{}) StatusError {
	if e.data == nil {
		e.data = make(map[string]interface{})
	}
	e.data[key] = value
	return e
}
