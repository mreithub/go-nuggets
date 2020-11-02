package statuserror

type statusErrorWithData struct {
	statusError
	data map[string]interface{}
}

// Data -- returns extra data stored with this error
func (e statusErrorWithData) Data() map[string]interface{} { return e.data }

func (e statusErrorWithData) WithData(data map[string]interface{}) Error {
	e.data = data
	return e
}

func (e statusErrorWithData) WithField(key string, value interface{}) Error {
	// TODO this may update the original in-place. Do we want that?
	if e.data == nil {
		e.data = make(map[string]interface{})
	}

	e.data[key] = value
	return e
}
