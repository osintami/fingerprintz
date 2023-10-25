// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestHandlerWhoami(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"
	r := common.BuildRequest(http.MethodGet, "/v1/data/whoami", nil, qParams)
	w := httptest.NewRecorder()

	server.GetWhoamiHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	response := &common.WhoamiInfo{}
	err := json.Unmarshal(w.Body.Bytes(), response)
	assert.Nil(t, err)
}

func TestHandlerWhoamiErrNoData(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	qParams["ip"] = "4.3.2.1"
	r := common.BuildRequest(http.MethodGet, "/v1/data/whoami", nil, qParams)
	w := httptest.NewRecorder()

	server.GetWhoamiHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerWhoamiMissingKeys(t *testing.T) {
	server := nodsServer(false)

	qParams := make(map[string]string)
	r := common.BuildRequest(http.MethodGet, "/v1/data/whoami", nil, qParams)
	w := httptest.NewRecorder()

	server.GetWhoamiHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrMissingInputs)
	assert.Equal(t, content, w.Body.String())
}
