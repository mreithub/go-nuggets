package tpl

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
)

// ErrorData -- contains data useful for error templates
type ErrorData struct {
	// Error -- error that happened
	Error error

	// Code -- set this to a non-zero value to manually set a HTTP response code
	Code int

	// Request -- HTTP Request
	Request *http.Request

	// RedirectLocation -- if set, sets the 'Location' response header (use this in conjunction with 3xx codes)
	RedirectLocation string

	// Data -- extra data
	Data interface{}
}

// GetStatusCode -- returns the HTTP status code we'll use for this ErrorData
//
// - if .Code != 0, it'll be returned
// - next, we'll use .Error's .Code() method if present
// - then check if .RedirectLocation is set and return http.StatusFound
// - defaults to http.StatusInternalServerError
func (d ErrorData) GetStatusCode() int {
	if d.Code != 0 {
		return d.Code
	}

	var rc = http.StatusInternalServerError
	if d.Error != nil {
		if e, ok := d.Error.(interface {
			Code() int
		}); ok {
			return e.Code()
		}
	}

	if d.RedirectLocation != "" {
		return http.StatusFound
	}

	return rc
}

// GetStatusText -- returns http.StatusText() matching our status code
func (d ErrorData) GetStatusText() string {
	return http.StatusText(d.GetStatusCode())
}

// SendError -- Send an error page to the client
//
// We'll use data.GetStatusCode() as HTTP response code and
// try the following template files (in order)
// - errors/{statusCode}.html
// - error.html
//
// if none of these is found, we'll default to a really simple HTML
func (t *Templates) SendError(w http.ResponseWriter, data ErrorData) error {
	var code = data.GetStatusCode()

	// set Location header
	if data.RedirectLocation != "" {
		w.Header().Set("Location", data.RedirectLocation)
	}

	// look for matching templates
	var tpl *template.Template
	var tplPath string
	for _, path := range []string{fmt.Sprintf("errors/%d.html", code), "error.html"} {
		if tpl = t.Get(path); tpl != nil {
			tplPath = path
			break
		}
	}

	var err error
	if tpl == nil {
		tpl, err = template.New("empty").Parse("{{define \"main\"}}<h1>{{.GetStatusCode}} {{.GetStatusText}}</h1><p>{{.Error}}</p>{{end}}")
		if err != nil {
			return fmt.Errorf("failed to parse default error template: %w", err)
		}
	}

	w.WriteHeader(code)
	if err = tpl.ExecuteTemplate(w, "main", data); err != nil {
		io.WriteString(w, fmt.Sprintf("ERROR: failed to render %d error page", code))
		return fmt.Errorf("failed to render error page '%s': %w", tplPath, err)
	}

	return nil
}
