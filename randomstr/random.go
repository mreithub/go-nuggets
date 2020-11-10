package randomstr

import (
	"crypto/rand"
	"errors"
)

// we're using base32 characters here
// having exactly 2^5 chars has the advantagethat we can simply take 5 bits of each rand.Read() byte
// also, this base32 set skips some characters that can be confused with others
var alphanum = []rune("abcdefghijklmnopqrstuvwxyz234567")

// Generate -- Generates a random string identifier of the given length
//
// each generated character holds 5bit of information (i.e. we choose from 32 valid characters)
//
// uses crypto/rand internally
func Generate(length int) (string, error) {
	var noise = make([]byte, length)
	if n, err := rand.Read(noise); err != nil {
		return "", err
	} else if n != length {
		// rand.Read() will always return len(buff) bytes according to its documentation but better safe than sorry
		return "", errors.New("rand.Read() returned too few bytes")
	}

	var rc = make([]rune, length)
	// now that we have random 8-bit values, strip them down to 5 bit values and use them as index for our custom char array
	for i := 0; i < length; i++ {
		var c = alphanum[noise[i]%32]
		rc[i] = c
	}
	return string(rc), nil
}

// MustGenerate -- wraps Generate(), panicing on error
func MustGenerate(length int) string {
	var rc, err = Generate(length)
	if err != nil {
		panic(err)
	}
	return rc
}
