package tpl

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mreithub/go-nuggets/statuserror"
	"github.com/stretchr/testify/assert"
)

func TestErrorDataStatusCode(t *testing.T) {
	var d = ErrorData{
		Code:  http.StatusBadGateway,
		Error: statuserror.Errorf(http.StatusConflict, "hello: %w", errors.New("world")),
	}

	assert.Equal(t, "hello: world", d.Error.Error())
	assert.Equal(t, http.StatusBadGateway, d.GetStatusCode())

	d.Code = 0
	assert.Equal(t, http.StatusConflict, d.GetStatusCode())

	d.Error = errors.New("hello world")
	assert.Equal(t, http.StatusInternalServerError, d.GetStatusCode())
}

func TestSendErrorWithoutTEmplates(t *testing.T) {
	var templates = Templates{}

	var resp FakeResponse
	templates.SendError(&resp, ErrorData{
		Code:  http.StatusAccepted,
		Error: errors.New("hello"),
	})

	//
	//assert.Equal(t, "<h1>...", string(resp.buff.Bytes()))
}
