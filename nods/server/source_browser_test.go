// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestBrowser(t *testing.T) {
	tools := mockToolbox()
	// NOTE:  the browser source object needs a third party YAML file full of REGEXES
	tools.DataPath = "./test/"
	source := NewBrowserSource(tools)
	assert.NotNil(t, source)

	inputs := make(map[string]string)
	inputs[CATEGORY_BROWSER] = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_BROWSER, inputs)

	assert.Nil(t, err)

	result := gjson.GetBytes(data, "UserAgent.Family")
	assert.Equal(t, "Googlebot", result.Str)

	result = gjson.GetBytes(data, "Device.Family")
	assert.Equal(t, "Spider", result.Str)

	// invalid category
	inputs[CATEGORY_EMAIL] = "1@2.com"
	data, err = source.CategoryInfo(context.TODO(), CATEGORY_EMAIL, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, data)
}

func TestBrowserNoRegexFile(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = ""
	source := NewBrowserSource(tools)
	assert.Nil(t, source)
}
