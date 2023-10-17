// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestMMDBProvider(t *testing.T) {
	reader := NewMockMaxmindReader("./test/ipsum.mmdb", false)
	provider := NewMMDBProvider(NewMockWatcher(), reader)
	assert.False(t, provider.IsCached())

	inputs := common.DataInputs{}
	inputs[CATEGORY_IPADDR] = "1.2.3.4"
	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "ipsum.blacklist.isBlacklisted")
	assert.True(t, result.Bool())
}

func TestMMDBProviderResync(t *testing.T) {
	reader := NewMockMaxmindReader("./test/ipsum.mmdb", false)
	provider := NewMMDBProvider(NewMockWatcher(), reader)
	assert.False(t, provider.IsCached())

	inputs := common.DataInputs{}
	inputs[CATEGORY_IPADDR] = "1.2.3.4"
	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "ipsum.blacklist.isBlacklisted")
	assert.True(t, result.Bool())

	reader.Resync()

	data, err = provider.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Nil(t, err)

	result = gjson.GetBytes(data, "ipsum.blacklist.isBlacklisted")
	assert.True(t, result.Bool())
}

func TestMMDBProviderNoResults(t *testing.T) {
	reader := NewMockMaxmindReader("./test/ipsum.mmdb", true)
	provider := NewMMDBProvider(NewMockWatcher(), reader)

	inputs := common.DataInputs{}
	inputs[CATEGORY_IPADDR] = "4.3.2.1"
	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.Nil(t, data)
}

func TestMMDBProviderBadCategory(t *testing.T) {
	reader := NewMockMaxmindReader("./test/ipsum.mmdb", false)
	provider := NewMMDBProvider(NewMockWatcher(), reader)

	inputs := common.DataInputs{}
	inputs[CATEGORY_DOMAIN] = "nope.com"
	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, data)
}
