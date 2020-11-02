package statuserror

import "errors"

// statusError -- simple statusError implementation
// can't hold extra data (will always return nil in that case)
type statusError struct {
	err  error
	code int
}

// Code -- returns the status code associated with this StatusError
func (e statusError) Code() int                    { return e.code }
func (e statusError) Data() map[string]interface{} { return nil }
func (e statusError) Error() string                { return e.err.Error() }
func (e statusError) Unwrap() error                { return errors.Unwrap(e.err) }

func (e statusError) WithData(data map[string]interface{}) Error {
	return statusErrorWithData{
		statusError: e,
		data:        data,
	}
}

func (e statusError) WithField(key string, value interface{}) Error {
	return statusErrorWithData{
		statusError: e,
	}.WithField(key, value)
}
