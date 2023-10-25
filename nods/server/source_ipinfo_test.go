// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestIpInfo(t *testing.T) {
	source := NewIpInfoSource(mockToolbox(), NewMockIpInfo(false))
	assert.True(t, source.IsCached())

	inputs := make(map[string]string)
	inputs[CATEGORY_IPADDR] = "89.187.175.145"

	// standard path
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Nil(t, err)
	result := gjson.GetBytes(data, "region")
	assert.Equal(t, "Texas", result.Str)
}

func TestIpInfoBadCategory(t *testing.T) {
	source := NewIpInfoSource(mockToolbox(), NewMockIpInfo(false))

	// invalid category
	inputs := make(map[string]string)
	inputs[CATEGORY_PHONE] = "18001234567"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, data)
}

func TestIpInfoCallFail(t *testing.T) {
	tools := mockToolbox()
	source := NewIpInfoSource(tools, NewMockIpInfo(true))

	inputs := make(map[string]string)
	inputs[CATEGORY_IPADDR] = "89.187.175.145"

	// network call failed
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.NotNil(t, err)
	assert.NotNil(t, data)
}
