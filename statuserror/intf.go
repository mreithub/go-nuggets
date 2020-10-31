package statuserror

// Error - error implementation with a status code as well as optional data
type Error interface {
	error

	// Code -- returns the HTTP status code for this error
	//
	// when calling FromError() which a
	Code() int

	// Data -- returns any extra data stored with this error (may be used for error template pages)
	// may return nil
	Data() map[string]interface{}

	// Unwrap -- returns an underlying error if available
	Unwrap() error

	WithField(key string, value interface{}) Error
	WithData(map[string]interface{}) Error
}
