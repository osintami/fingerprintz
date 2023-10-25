// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestHttpExtractor(t *testing.T) {
	client := resty.New()
	job := createETLJob(client, "csv", "mmdb")
	he := NewHttpExtractor(job.Tools(), job.Source())

	httpmock.ActivateNonDefault(client.GetClient())
	content := "34.171.21\n1.2.3.4\n0.0.0.0\nnope\n"
	url := "https://raw.githubusercontent.com/osintami/fingerprintz/master/etlr/etl/test/test.csv?token=GHSAT0AAAAAACE25RHZPF6EBDUUB7JY4F7EZIS6QSA"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	err := he.Extract(job)
	assert.Nil(t, err)

	data, err := os.ReadFile("/tmp/test/test.csv")
	assert.Nil(t, err)
	assert.Equal(t, content, string(data))

	os.RemoveAll("/tmp/test")
}
