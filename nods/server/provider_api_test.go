// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestAPIProvider(t *testing.T) {
	tools := mockToolbox()
	tools.Client = resty.New()

	apiKey := "OSINTAMI_API_KEY"

	headers := make([]HeaderInfo, 1)
	headers[0].Key = "content"
	headers[0].Value = "json/application"

	remote := &API{
		Name:    "osintami",
		URL:     "https://nope.com",
		Method:  http.MethodGet,
		Auth:    &BasicAuthInfo{Name: "test", Pass: "test"},
		ApiKey:  &apiKey,
		Headers: headers}
	provider := NewAPIProvider(tools, remote)
	assert.True(t, provider.IsCached())

	inputs := common.DataInputs{}
	inputs[CATEGORY_DOMAIN] = "nope.com"

	// set up mock HTTP stuff
	httpmock.ActivateNonDefault(tools.Client.GetClient())
	defer httpmock.DeactivateAndReset()
	content := "{\"IsFake\":true}"
	url := "https://nope.com"
	httpmock.RegisterResponder(
		http.MethodGet,
		url,
		httpmock.NewStringResponder(http.StatusOK, content))

	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "IsFake")
	assert.True(t, result.Bool())
}

func TestAPIProviderGreynoiseFailure(t *testing.T) {
	tools := mockToolbox()
	tools.Client = resty.New()

	remote := &API{
		Name:   GREYNOISE,
		URL:    "https://nope.com",
		Method: http.MethodGet}
	provider := NewAPIProvider(tools, remote)

	inputs := common.DataInputs{}
	inputs[CATEGORY_IPADDR] = "1.2.3.4"

	// set up mock HTTP stuff
	httpmock.ActivateNonDefault(tools.Client.GetClient())
	defer httpmock.DeactivateAndReset()
	content := "{\"IsFake\":false}"
	url := "https://nope.com"
	httpmock.RegisterResponder(
		http.MethodGet,
		url,
		httpmock.NewStringResponder(http.StatusNotFound, content))

	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "IsFake")
	assert.False(t, result.Bool())
}

func TestAPIProviderAPIVoidFailure(t *testing.T) {
	tools := mockToolbox()
	tools.Client = resty.New()

	remote := &API{
		Name:   APIVOID,
		URL:    "https://nope.com",
		Method: http.MethodGet}
	provider := NewAPIProvider(tools, remote)

	inputs := common.DataInputs{}
	inputs[CATEGORY_IPADDR] = "1.2.3.4"

	// set up mock HTTP stuff
	httpmock.ActivateNonDefault(tools.Client.GetClient())
	defer httpmock.DeactivateAndReset()
	content := "{\"error\":\"we don't use HTTP error codes for missing data\"}"
	url := "https://nope.com"
	httpmock.RegisterResponder(
		http.MethodGet,
		url,
		httpmock.NewStringResponder(http.StatusOK, content))

	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.Nil(t, data)
}

func TestAPIProviderAPICallFailure(t *testing.T) {
	tools := mockToolbox()
	tools.Client = resty.New()

	remote := &API{
		Name:   "any",
		URL:    "https://nope.com",
		Method: http.MethodGet}
	provider := NewAPIProvider(tools, remote)

	inputs := common.DataInputs{}
	inputs[CATEGORY_IPADDR] = "1.2.3.4"

	// set up mock HTTP stuff
	httpmock.ActivateNonDefault(tools.Client.GetClient())
	defer httpmock.DeactivateAndReset()
	content := "{}"
	url := "https://nope.com"
	httpmock.RegisterResponder(
		http.MethodGet,
		url,
		httpmock.NewStringResponder(http.StatusTooManyRequests, content))

	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_DOMAIN, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.Nil(t, data)
}
