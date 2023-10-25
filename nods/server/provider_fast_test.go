// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestFastProvider(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = "./test/"
	provider, err := NewFastDBProvider(tools, "fakefilter")
	assert.Nil(t, err)
	assert.False(t, provider.IsCached())

	inputs := common.DataInputs{}
	inputs[CATEGORY_DOMAIN] = "nope.com"
	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "Result.IsFake")
	assert.True(t, result.Bool())
}

func TestFastProviderNoResults(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = "./test/"
	provider, err := NewFastDBProvider(tools, "fakefilter")
	assert.Nil(t, err)
	assert.False(t, provider.IsCached())

	inputs := common.DataInputs{}
	inputs[CATEGORY_DOMAIN] = "nunya.com"
	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.Nil(t, data)
}

// TODO:  test changing fakefilter.json schema and see if it really reloads properly
