// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestSpamhaus(t *testing.T) {

	rbl := NewMockRealtimeBlackholeList()
	source := NewSpamhausSource(mockToolbox(), rbl)

	// test IP lookup
	inputs := make(map[string]string)
	inputs[CATEGORY_IPADDR] = "1.2.3.4"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "IPAddr")
	assert.Equal(t, "1.2.3.4", result.Str)

	result = gjson.GetBytes(data, "Domain")
	assert.Equal(t, "", result.Str)

	result = gjson.GetBytes(data, "IsBlacklisted")
	assert.True(t, result.Bool())

	// test DOMAIN lookup
	inputs[CATEGORY_DOMAIN] = "nunya.com"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Nil(t, err)

	result = gjson.GetBytes(data, "IPAddr")
	assert.Equal(t, "", result.Str)

	result = gjson.GetBytes(data, "Domain")
	assert.Equal(t, "nunya.com", result.Str)

	result = gjson.GetBytes(data, "IsBlacklisted")
	assert.True(t, result.Bool())

	// test error path
	inputs[CATEGORY_PHONE] = "18001234567"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, data)
}
