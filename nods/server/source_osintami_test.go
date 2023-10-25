// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestOsintamiSourceMissingInputs(t *testing.T) {
	source := NewInternalSource(NewMockDataRouter(false), mockToolbox())
	assert.False(t, source.IsCached())

	inputs := make(map[string]string)
	// missing inputs, all should return the same error and data
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_BROWSER, inputs)
	assert.Equal(t, ErrMissingInputs, err)
	assert.Equal(t, "{}", string(data))

	data, err = source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, ErrMissingInputs, err)
	assert.Equal(t, "{}", string(data))

	data, err = source.CategoryInfo(context.TODO(), CATEGORY_EMAIL, inputs)
	assert.Equal(t, ErrMissingInputs, err)
	assert.Equal(t, "{}", string(data))

	data, err = source.CategoryInfo(context.TODO(), CATEGORY_RULE, inputs)
	assert.Equal(t, ErrMissingInputs, err)
	assert.Equal(t, "{}", string(data))
}

func TestOsintamiSourceInvalidCategory(t *testing.T) {
	source := NewInternalSource(NewMockDataRouter(false), mockToolbox())
	inputs := make(map[string]string)
	inputs["ip"] = "1.2.3.4"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, data)
}

func TestOsintamiSourceInvalidInput(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = "./test/"
	source := NewInternalSource(NewMockDataRouter(false), tools)

	inputs := make(map[string]string)

	inputs[CATEGORY_BROWSER] = "*"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_BROWSER, inputs)
	assert.Nil(t, err)
	result := gjson.GetBytes(data, "UserAgent.Family")
	assert.Equal(t, "Other", result.Str)

	inputs[CATEGORY_PHONE] = "*"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.NotNil(t, err)
	result = gjson.GetBytes(data, "IsValid")
	assert.False(t, result.Bool())

	inputs[CATEGORY_EMAIL] = "*"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_EMAIL, inputs)
	assert.Nil(t, err)
	result = gjson.GetBytes(data, "isValidEmail")
	assert.False(t, result.Bool())

	inputs[CATEGORY_RULE] = "*"
	inputs[common.INPUT_RULE] = "*"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_RULE, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.Equal(t, "{}", string(data))
}
