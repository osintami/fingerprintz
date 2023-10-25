// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestWhoIs(t *testing.T) {
	tools := mockToolbox()
	source := NewWhoisSource(tools, NewMockWhois(false, false))

	inputs := make(map[string]string)
	inputs[CATEGORY_DOMAIN] = "osintami.com"

	// standard path
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Nil(t, err)
	result := gjson.GetBytes(data, "DomainAgeInDays")
	assert.Greater(t, result.Num, float64(300))

	// unsupported category
	inputs[CATEGORY_PHONE] = "18001234567"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, data)

	// top level domain
	inputs[CATEGORY_DOMAIN] = "com"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	result = gjson.GetBytes(data, "DomainAgeInDays")
	assert.Equal(t, float64(0), result.Num)
}

func TestWhoIsFailParse(t *testing.T) {
	tools := mockToolbox()
	source := NewWhoisSource(tools, NewMockWhois(true, false))

	inputs := make(map[string]string)
	inputs[CATEGORY_DOMAIN] = "osintami.com"

	_, err := source.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.NotNil(t, err)
}

func TestWhoIsFailGet(t *testing.T) {
	tools := mockToolbox()
	source := NewWhoisSource(tools, NewMockWhois(false, true))

	inputs := make(map[string]string)
	inputs[CATEGORY_DOMAIN] = "osintami.com"

	_, err := source.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.NotNil(t, err)
}
