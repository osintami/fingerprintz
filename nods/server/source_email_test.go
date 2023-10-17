// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestEmail(t *testing.T) {

	source := NewEmailSource(mockToolbox())

	inputs := make(map[string]string)
	inputs[CATEGORY_EMAIL] = "12345678+alias@gmail.com"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_EMAIL, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "IsWeirdUserName")
	assert.Equal(t, true, result.Bool())

	result = gjson.GetBytes(data, "IsEmailAlias")
	assert.Equal(t, true, result.Bool())

	// test error path
	inputs[CATEGORY_PHONE] = "18001234567"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, data)
}

func TestEmailMethods(t *testing.T) {
	email := &EmailSource{}
	check := "abc12345!*#XYZ^"

	assert.Equal(t, 4, email.nonAlphanumerics(check))
	assert.Equal(t, 5, email.numerics(check))
	assert.True(t, email.isNonHumanDomainName(check))
	assert.False(t, email.isNonHumanUserName("iamnormal"))
	assert.True(t, email.isNonHumanDomainName(check))
	assert.True(t, email.isRiskyTopLevelDomain("1.xyz"))
	assert.False(t, email.isRiskyTopLevelDomain("nope.com"))
}
