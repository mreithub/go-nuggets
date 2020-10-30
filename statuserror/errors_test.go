package statuserror

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var err = New(http.StatusNotFound, "hello")

	assert.Equal(t, http.StatusNotFound, err.Code())
	assert.Equal(t, "hello", err.Error())
}
