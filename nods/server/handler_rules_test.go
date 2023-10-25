// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestHandlerRuleNoData(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"
	qParams["rule"] = "[ip/ipsum/blacklist.isBlacklisted] || [ip/uhb/blacklist.isBlacklisted]"

	r := common.BuildRequest(http.MethodGet, "/v1/data/rule", nil, qParams)
	w := httptest.NewRecorder()

	server.GetEvaluateHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	content := "{\"Item\":\"ip/uhb/blacklist.isBlacklisted\",\"Result\":{\"Type\":0,\"Bool\":true},\"Keys\":{\"ip\":\"1.2.3.4\"},\"Error\":\"\"}\n"
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerRuleMissingKeys(t *testing.T) {
	server := nodsServer(false)

	r := buildItemRequest(nil, nil)
	w := httptest.NewRecorder()

	server.GetEvaluateHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrMissingInputs)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerRuleBadRule(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	r := common.BuildRequest(http.MethodGet, "/v1/data/rule", nil, qParams)
	w := httptest.NewRecorder()

	server.GetEvaluateHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	content := common.BuildErrorResponse(ErrMissingInputs)
	assert.Equal(t, content, w.Body.String())
}
