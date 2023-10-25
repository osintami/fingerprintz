// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestHandlerSchemaUser(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["role"] = "user"
	r := common.BuildRequest(http.MethodGet, "/v1/data/items", nil, qParams)
	w := httptest.NewRecorder()

	server.DictionaryHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	assert.True(t, strings.HasPrefix(w.Body.String(), "[{\"Item\":"))
	assert.True(t, strings.HasSuffix(w.Body.String(), "\"}]\n"))
}

func TestHandlerSchemaUserCSV(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["role"] = "user"
	qParams["csv"] = "true"
	r := common.BuildRequest(http.MethodGet, "/v1/data/items", nil, qParams)

	w := httptest.NewRecorder()
	server.DictionaryHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	assert.True(t, strings.HasPrefix(w.Body.String(), "CATEGORY,   ITEM,       SOURCE,                                ENABLED,   TYPE,      DESCRIPTION\n"))
}

func TestHandlerSchemaAdmin(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["role"] = "admin"
	r := common.BuildRequest(http.MethodGet, "/v1/data/items", nil, qParams)
	w := httptest.NewRecorder()

	server.DictionaryHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	assert.True(t, strings.HasPrefix(w.Body.String(), "[{\"Item\":"))
	assert.True(t, strings.HasSuffix(w.Body.String(), "\"}]\n"))
}

func TestHandlerSchemaAdminCSV(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["role"] = "admin"
	qParams["csv"] = "true"

	r := common.BuildRequest(http.MethodGet, "/v1/data/items", nil, qParams)
	w := httptest.NewRecorder()

	server.DictionaryHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	assert.True(t, strings.HasPrefix(w.Body.String(), "CATEGORY,   ITEM,       SOURCE,                    ENABLED,   TYPE,      DESCRIPTION\n"))
}

func TestHandlerSchemaRoleMissing(t *testing.T) {
	server := nodsServer(false)

	r := common.BuildRequest(http.MethodGet, "/v1/data/items", nil, nil)
	w := httptest.NewRecorder()

	server.DictionaryHandler(w, r)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	content := common.BuildErrorResponse(ErrInvalidUserRole)
	assert.Equal(t, content, w.Body.String())
}
