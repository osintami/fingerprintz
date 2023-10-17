// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestFileExtractor(t *testing.T) {
	err := os.MkdirAll("/tmp/test", 01777)
	assert.Nil(t, err)

	// source file
	err = os.WriteFile("/tmp/test.in", []byte("..."), 0644)
	assert.Nil(t, err)

	fe := NewFileExtractor("/tmp/test.in")
	etlr := createETLJob(resty.New(), "csv", "mmdb")

	err = fe.Extract(etlr)
	assert.Nil(t, err)

	// extracted file
	_, err = os.Stat("/tmp/test/test.csv")
	assert.Nil(t, err)

	// original content
	content, err := os.ReadFile("/tmp/test.in")
	assert.Nil(t, err)

	// extracted content
	data, err := os.ReadFile("/tmp/test/test.csv")
	assert.Nil(t, err)

	// compare original/extracted content
	assert.Equal(t, string(content), string(data))

	// cleanup
	os.RemoveAll("/tmp/test")
	os.RemoveAll("/tmp/test.in")
}

func TestFileExtractorFileDNE(t *testing.T) {
	fe := NewFileExtractor("/tmp/nope.in")
	etlr := createETLJob(resty.New(), "csv", "mmdb")

	err := fe.Extract(etlr)
	assert.NotNil(t, err)

	// cleanup
	os.RemoveAll("/tmp/test")
}
