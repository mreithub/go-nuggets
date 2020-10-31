package statuserror

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var err = New(http.StatusNotFound, "hello")

	assert.Equal(t, http.StatusNotFound, err.Code())
	assert.Equal(t, "hello", err.Error())
	assert.Nil(t, err.Data())

	var err2 = err.WithField("hello", "world")

	assert.Nil(t, err.Data())
	assert.Equal(t, map[string]interface{}{
		"hello": "world",
	}, err2.Data())

	// err2 will be altered here - figure out if that's what we want or not
	var err3 = err2.WithField("answer", 42)
	assert.Equal(t, map[string]interface{}{
		"hello":  "world",
		"answer": 42,
	}, err2.Data())

	assert.Equal(t, err2, err3)
}

func TestUnwrap(t *testing.T) {
	var err = Errorf(http.StatusForbidden, "reason: %w", errors.New("lol"))
	assert.Equal(t, "reason: lol", err.Error())
	assert.Equal(t, http.StatusForbidden, err.Code())
	assert.Equal(t, errors.New("lol"), err.Unwrap())
	assert.Equal(t, errors.New("lol"), errors.Unwrap(err))

	var err2 = FromError(err.(error))
	assert.Equal(t, "reason: lol", err2.Error())
	assert.Equal(t, http.StatusForbidden, err2.Code())
	assert.Equal(t, errors.New("lol"), err2.Unwrap())
	assert.Equal(t, errors.New("lol"), errors.Unwrap(err2))

	var err3 = From(http.StatusUnauthorized, err.(error))
	assert.Equal(t, "reason: lol", err3.Error())
	assert.Equal(t, http.StatusUnauthorized, err3.Code())
	assert.Equal(t, errors.New("lol"), err3.Unwrap())
	assert.Equal(t, errors.New("lol"), errors.Unwrap(err3))

}
