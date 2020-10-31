package statuserror

// Error - error implementation with a status code as well as optional data
type Error interface {
	error

	// Code -- returns the HTTP status code for this error
	//
	// when calling FromError() which a
	Code() int

	// Unwrap -- returns an underlying error if available
	Unwrap() error

	WithField(key string, value interface{}) Error
	WithData(map[string]interface{}) Error
}
