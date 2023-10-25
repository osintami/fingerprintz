// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParams(t *testing.T) {
	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	// GET
	r := BuildRequest(http.MethodGet, "/v1/data/category/ip", nil, qParams)
	params := NewParameterHelper()
	keys, err := params.Keys(r)
	assert.Nil(t, err)
	assert.Equal(t, "1.2.3.4", keys["ip"])

	// POST
	r = BuildRequest(http.MethodPost, "/v1/data/category/ip", nil, qParams)
	keys, err = params.Keys(r)
	assert.Nil(t, err)
	assert.Equal(t, "1.2.3.4", keys["ip"])

	// Corrupted POST body
	r = httptest.NewRequest(
		http.MethodPost,
		"/v1/data/category/ip",
		strings.NewReader("..."))

	keys, err = params.Keys(r)
	assert.NotNil(t, err)
	assert.Nil(t, keys)
}
