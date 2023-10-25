// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCategoryNoData(t *testing.T) {
	server := nodsServer(true)

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	pParams := make(map[string]string)
	pParams["category"] = "ip"

	r := common.BuildRequest(http.MethodGet, "/v1/data/category/ip", pParams, qParams)
	w := httptest.NewRecorder()

	server.GetCategoryHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrItemNotFound)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerCategoryMissingKeys(t *testing.T) {
	server := nodsServer(false)

	r := common.BuildRequest(http.MethodGet, "/v1/data/category/ip", nil, nil)
	w := httptest.NewRecorder()

	server.PostCategoryHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrMissingInputs)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerCategoryInvalidCategory(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	pParams := make(map[string]string)
	pParams["category"] = "rule"
	r := common.BuildRequest(http.MethodGet, "/v1/data/category/rule", pParams, qParams)
	w := httptest.NewRecorder()

	server.PostCategoryHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrInvalidCategoryParam)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerCategoryIP(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	pParams := make(map[string]string)
	pParams["category"] = "ip"

	r := common.BuildRequest(http.MethodGet, "/v1/data/category/ip", pParams, qParams)
	w := httptest.NewRecorder()

	server.PostCategoryHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	content := "[{\"Item\":\"ip/ipsum/blacklist.isBlacklisted\",\"Result\":{\"Type\":1,\"Bool\":true},\"Keys\":{\"ip\":\"1.2.3.4\"},\"Error\":\"\"},{\"Item\":\"ip/uhb/blacklist.isBlacklisted\",\"Result\":{\"Type\":1,\"Bool\":true},\"Keys\":{\"ip\":\"1.2.3.4\"},\"Error\":\"\"}]\n"
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerCategoryCSV(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"
	qParams["csv"] = "true"

	pParams := make(map[string]string)
	pParams["category"] = "ip"

	r := common.BuildRequest(http.MethodGet, "/v1/data/category/ip", pParams, qParams)
	w := httptest.NewRecorder()

	server.GetCategoryHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	content := "category, source, item,                    key,     value\nip,       ipsum,  blacklist.isBlacklisted, 1.2.3.4, true\nip,       uhb,    blacklist.isBlacklisted, 1.2.3.4, true\n"
	assert.Equal(t, content, w.Body.String())
}
