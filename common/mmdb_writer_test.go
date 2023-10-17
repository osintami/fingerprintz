// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"encoding/json"
	"net"
	"os"
	"testing"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestMaxmindWriter(t *testing.T) {
	// create an mmdb file with one entry
	_, cidr, _ := net.ParseCIDR("1.2.3.4/32")

	writer := NewMaxmindWriter("test")
	entry := mmdbtype.Map{
		"test": mmdbtype.Map{
			"blacklist": mmdbtype.Map{
				"isBlacklisted": mmdbtype.Bool(true),
			},
		},
	}
	err := writer.Insert(cidr, entry)
	assert.Nil(t, err)

	err = writer.Close("/tmp/test.mmdb")
	assert.Nil(t, err)

	// read the mmdb file we just created
	reader, err := NewMaxmindReader("/tmp/test.mmdb")
	assert.Nil(t, err)
	data, err := reader.Lookup(net.ParseIP("1.2.3.4"))
	assert.Nil(t, err)

	out, err := json.Marshal(data)
	assert.Nil(t, err)
	result := gjson.GetBytes(out, "test.blacklist.isBlacklisted")
	assert.True(t, result.Bool())

	// reload the underlying data
	reader.Resync()

	// re-check we can read the data
	data, err = reader.Lookup(net.ParseIP("1.2.3.4"))
	assert.Nil(t, err)

	out, err = json.Marshal(data)
	assert.Nil(t, err)
	result = gjson.GetBytes(out, "test.blacklist.isBlacklisted")
	assert.True(t, result.Bool())

	// cleanup
	os.Remove("/tmp/test/mmdb")
}

func TestMaxmindWriterNoData(t *testing.T) {
	// create an mmdb file with one entry
	_, cidr, _ := net.ParseCIDR("1.2.3.4/32")

	writer := NewMaxmindWriter("test")
	entry := mmdbtype.Map{
		"test": mmdbtype.Map{
			"blacklist": mmdbtype.Map{
				"isBlacklisted": mmdbtype.Bool(true),
			},
		},
	}
	err := writer.Insert(cidr, entry)
	assert.Nil(t, err)

	err = writer.Close("/tmp/test.mmdb")
	assert.Nil(t, err)

	// read the mmdb file we just created
	reader, err := NewMaxmindReader("/tmp/test.mmdb")
	assert.Nil(t, err)
	data, err := reader.Lookup(net.ParseIP("4.3.2.1"))
	assert.Equal(t, ErrNoDataPresent, err)
	assert.Nil(t, data)
}

func TestMaxmindWriterEmptyDatabaseType(t *testing.T) {
	// test close invalid file
	writer := NewMaxmindWriter("")
	assert.Nil(t, writer)
}

func TestMaxmindWriterInvalidFileOnClose(t *testing.T) {
	// test close invalid file
	writer := NewMaxmindWriter("test")
	err := writer.Close("././././.")
	assert.NotNil(t, err)
}
