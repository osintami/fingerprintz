// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestPhoneMexico(t *testing.T) {

	source := NewPhoneSource(mockToolbox())
	inputs := make(map[string]string)

	// test Mexico
	inputs[CATEGORY_PHONE] = "523220543290"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "CountryISOCode")
	assert.Equal(t, "MX", result.Str)

	result = gjson.GetBytes(data, "NationalNumber")
	assert.Equal(t, "322 054 3290", result.Str)

	result = gjson.GetBytes(data, "CountryCode")
	assert.Equal(t, float64(52), result.Num)

	result = gjson.GetBytes(data, "IsValid")
	assert.True(t, result.Bool())

	// test toll free
	inputs[CATEGORY_PHONE] = "1 888 234 5678"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Nil(t, err)

	result = gjson.GetBytes(data, "CountryISOCode")
	assert.Equal(t, "US", result.Str)

	result = gjson.GetBytes(data, "NationalNumber")
	assert.Equal(t, "(888) 234-5678", result.Str)

	result = gjson.GetBytes(data, "CountryCode")
	assert.Equal(t, float64(1), result.Num)

	result = gjson.GetBytes(data, "IsValid")
	assert.True(t, result.Bool())
}

func TestPhoneTollFree(t *testing.T) {

	source := NewPhoneSource(mockToolbox())
	inputs := make(map[string]string)

	// test toll free
	inputs[CATEGORY_PHONE] = "1 888 234 5678"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "CountryISOCode")
	assert.Equal(t, "US", result.Str)

	result = gjson.GetBytes(data, "NationalNumber")
	assert.Equal(t, "(888) 234-5678", result.Str)

	result = gjson.GetBytes(data, "CountryCode")
	assert.Equal(t, float64(1), result.Num)

	result = gjson.GetBytes(data, "IsValid")
	assert.True(t, result.Bool())

	result = gjson.GetBytes(data, "IsTollFree")
	assert.True(t, result.Bool())
}

func TestPhoneErrors(t *testing.T) {

	source := NewPhoneSource(mockToolbox())
	inputs := make(map[string]string)

	// test bogus number
	inputs[CATEGORY_PHONE] = "nope"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	result := gjson.GetBytes(data, "IsValid")
	assert.False(t, result.Bool())

	// test error path
	inputs[CATEGORY_IPADDR] = "1.2.3.4"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, data)

	// test bogus long number
	inputs[CATEGORY_PHONE] = "00000000000"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	result = gjson.GetBytes(data, "IsValid")
	assert.False(t, result.Bool())

	// test US phone number nope
	inputs[CATEGORY_PHONE] = "*"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	result = gjson.GetBytes(data, "IsValid")
	assert.False(t, result.Bool())

}
