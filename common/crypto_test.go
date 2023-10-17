// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	key := fmt.Sprintf("%032s", "test-aes-key")
	aes := NewAES([]byte(key))
	edata, err := aes.Encrypt([]byte("hello world"))
	assert.Nil(t, err)

	ddata, err := aes.Decrypt(edata)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", string(ddata))
}

func TestCryptoBadKey(t *testing.T) {
	aes := NewAES(nil)
	assert.Nil(t, aes)
}

func TestCryptoBadData(t *testing.T) {
	key := fmt.Sprintf("%032s", "test-aes-key")

	aes := NewAES([]byte(key))
	data, err := aes.Decrypt(nil)
	assert.NotNil(t, err)
	assert.Equal(t, []byte("[]"), data)
}

func TestSHA256(t *testing.T) {
	out := SHA256ToBase16Lowercase("hello world")
	assert.Equal(t, out, SHA256ToBase16Lowercase("hello world"))
	assert.NotEqual(t, out, SHA256ToBase16Lowercase("world hello"))
}
