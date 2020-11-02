package statuserror

// Error - error implementation with a status code as well as optional data
//
//
// Error is immutable, all modifying methods return a new instance
//
// Internally, there are two StatusError implementations: the default one not having a Data field (.Data() will always return nil).
// As soon as you call .WithData() or .WithField() a new instance with Data field will be returned.
// Having a map[] field does break go's comparison operator though. This is relevant when reusing error instances and comparing them, e.g.:
//
// ```
// var ErrNotFound = statuserror.New(http.statusNotFound, "whatever you're looking for isn't here")
//
//  func main() {
//	  // ...
//	  var err = someFunction()
//    if err == ErrNotFound {
//      // we can only do this with data-less StatusErrors
//    }
//  }
// ```
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

	// WithField -- add a data field to this error
	WithField(key string, value interface{}) Error
	// WithData -- returns a copy of this error with Data set to data
	WithData(data map[string]interface{}) Error
}
