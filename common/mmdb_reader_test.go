// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxmindReaderInvalidFileOnOpen(t *testing.T) {
	reader, err := NewMaxmindReader("")
	assert.NotNil(t, err)
	assert.NotNil(t, reader)
	// internal MMDB reader is nil
	out, err := reader.Lookup(net.ParseIP("1.2.3.4"))
	assert.Equal(t, ErrNoDataPresent, err)
	assert.Nil(t, out)
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

func TestMaxmindReader(t *testing.T) {
	reader, err := NewMaxmindReader("./test/test.mmdb")
	assert.Nil(t, err)
	// value found
	out, err := reader.Lookup(net.ParseIP("1.2.3.4"))
	assert.Nil(t, err)
	fmt.Println("******", string(out))
}

func TestMaxmindReaderNoData(t *testing.T) {
	reader, err := NewMaxmindReader("./test/test.mmdb")
	assert.Nil(t, err)
	// no value found
	out, err := reader.Lookup(net.ParseIP("4.3.2.1"))
	assert.Equal(t, ErrNoDataPresent, err)
	assert.Nil(t, out)
}
