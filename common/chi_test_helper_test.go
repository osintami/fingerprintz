// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRequest(t *testing.T) {
	qParams := make(map[string]string)
	pParams := make(map[string]string)

	// GET
	qParams["ip"] = "1.2.3.4"
	pParams["category"] = "ip"
	r := BuildRequest(http.MethodGet, "/images/any.png", pParams, qParams)
	assert.Equal(t, "1.2.3.4", QueryParam(r, "ip"))
	assert.Equal(t, "ip", PathParam(r, "category"))

	// POST
	qParams["ip"] = "1.2.3.4"
	pParams["category"] = "ip"
	r = BuildRequest(http.MethodPost, "/images/any.png", pParams, qParams)
	assert.Equal(t, "ip", PathParam(r, "category"))

	data := make([]byte, 16)
	_, err := r.Body.Read(data)
	assert.Nil(t, err)
	assert.Equal(t, "{\"ip\":\"1.2.3.4\"}", string(data))
}

func TestBuildErrorResponse(t *testing.T) {
	out := BuildErrorResponse(ErrItemNotFound)
	assert.True(t, strings.Contains(out, ErrItemNotFound.Error()))
}
