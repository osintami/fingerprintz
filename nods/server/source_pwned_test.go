// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestPwnedBreaches(t *testing.T) {
	source := NewPwnedSource(mockToolbox(), NewMockPwned(false, false, false))
	assert.True(t, source.IsCached())

	inputs := make(map[string]string)
	inputs[CATEGORY_EMAIL] = "test@gmail.com"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_EMAIL, inputs)

	assert.Nil(t, err)

	result := gjson.GetBytes(data, "PastebinCount")
	assert.Equal(t, float64(1), result.Num)

	result = gjson.GetBytes(data, "BreachCount")
	assert.Equal(t, float64(1), result.Num)

	result = gjson.GetBytes(data, "BreachAgeInDays")
	assert.GreaterOrEqual(t, float64(0), result.Num)
}

func TestPwnedPassword(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// password API
	content := "{\"PastebinCount\":3940}"
	url := "https://api.pwnedpasswords.com/range/76272"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(200, content))

	source := NewPwnedSource(mockToolbox(), NewMockPwned(false, false, false))

	inputs := make(map[string]string)
	inputs[CATEGORY_PASSWORD] = "nope"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_PASSWORD, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "PastebinCount")
	assert.Equal(t, float64(0), result.Num)
}

func TestPwnedCallToPasswordFailed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// password API
	content := "{\"PastebinCount\":3940}"
	url := "https://api.pwnedpasswords.com/range/76272"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(200, content))

	source := NewPwnedSource(mockToolbox(), NewMockPwned(false, false, true))

	inputs := make(map[string]string)
	inputs[CATEGORY_PASSWORD] = "nope"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_PASSWORD, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)

	result := gjson.GetBytes(data, "PastebinCount")
	assert.Equal(t, float64(0), result.Num)
}

func TestPwnedBadCategory(t *testing.T) {
	source := NewPwnedSource(mockToolbox(), NewMockPwned(false, false, false))

	// invalid category
	inputs := make(map[string]string)
	inputs[CATEGORY_PHONE] = "18001234567"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, data)
}

func TestPwnedCallToBreachesFail(t *testing.T) {
	tools := mockToolbox()
	source := NewPwnedSource(tools, NewMockPwned(false, true, false))

	inputs := make(map[string]string)
	inputs[CATEGORY_EMAIL] = "1@2.com"

	// network call failed
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_EMAIL, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.NotNil(t, data)
}

func TestPwnedCallToPastesFail(t *testing.T) {
	tools := mockToolbox()
	source := NewPwnedSource(tools, NewMockPwned(true, false, false))

	inputs := make(map[string]string)
	inputs[CATEGORY_EMAIL] = "1@2.com"

	// network call failed
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_EMAIL, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.NotNil(t, data)
}

func TestPwnedNotBreachedOrPwned(t *testing.T) {
	tools := mockToolbox()
	source := NewPwnedSource(tools, NewMockPwned(false, false, true))

	inputs := make(map[string]string)
	inputs[CATEGORY_EMAIL] = "1@2.com"

	// network call failed
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_EMAIL, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "PastebinCount")
	assert.Equal(t, float64(0), result.Num)

	result = gjson.GetBytes(data, "BreachCount")
	assert.Equal(t, float64(0), result.Num)
}
