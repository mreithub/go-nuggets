package tpl

import (
	"bytes"
	"net/http"
)

// FakeResponse -- imitates a http.ResponseWriter
type FakeResponse struct {
	code    int
	buff    bytes.Buffer
	headers http.Header
}

func (r *FakeResponse) Header() http.Header            { return r.headers }
func (r *FakeResponse) Write(data []byte) (int, error) { return r.buff.Write(data) }
func (r *FakeResponse) WriteHeader(code int)           { r.code = code }
