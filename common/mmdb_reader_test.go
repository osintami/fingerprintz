// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxmindReaderInvalidFileOnOpen(t *testing.T) {
	reader, err := NewMaxmindReader("")
	assert.NotNil(t, err)
	assert.NotNil(t, reader)
}

func TestMaxmindReaderFilename(t *testing.T) {
	reader, _ := NewMaxmindReader("nope")
	assert.Equal(t, "nope", reader.FileName())
}

func TestMaxmindReaderResyncFail(t *testing.T) {
	reader, err := NewMaxmindReader("")
	assert.NotNil(t, err)
	reader.Resync()
}
