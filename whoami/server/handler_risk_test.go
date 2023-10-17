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

func TestHandlerRiskNodsFailure(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(false, false), NewMockNods(true))

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"
	qParams["fingerprint"] = "nope"

	r := common.BuildRequest(http.MethodPost, "/v1/data/risk", nil, qParams)
	w := httptest.NewRecorder()

	server.PostRiskHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandlerRisktDecodeFailure(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(true, false), NewMockNods(false))
	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"
	qParams["fingerprint"] = "nope"

	r := common.BuildRequest(http.MethodGet, "/v1/data/risk", nil, qParams)
	w := httptest.NewRecorder()

	server.GetRiskHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandlerRisktEncodeFailure(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(false, true), NewMockNods(false))
	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"
	qParams["fingerprint"] = "nope"

	r := common.BuildRequest(http.MethodGet, "/v1/data/risk", nil, qParams)
	w := httptest.NewRecorder()

	server.GetRiskHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandleRiskNoInputs(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(false, false), NewMockNods(false))

	qParams := make(map[string]string)

	r := common.BuildRequest(http.MethodGet, "/v1/data/risk", nil, qParams)
	w := httptest.NewRecorder()

	server.PostRiskHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleRiskNoFingerprint(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(false, false), NewMockNods(false))

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	r := common.BuildRequest(http.MethodGet, "/v1/data/risk", nil, qParams)
	w := httptest.NewRecorder()

	server.PostRiskHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerRisk(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(false, false), NewMockNods(false))

	qParams := make(map[string]string)
	qParams["email"] = "1@2.com"
	// location is South Korea
	qParams["ip"] = "1.2.3.4"
	qParams["ua"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
	qParams["hw"] = "01234567-89ABCDEF-01234567-89ABCDEF"
	qParams["pid"] = "osintami"

	r := common.BuildRequest(http.MethodGet, "/v1/data/fingerprint", nil, qParams)
	w := httptest.NewRecorder()

	server.GetFingerprintHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	fingerprint1 := w.Body.String()

	// change location to USA and run the print through NODS
	qParams["ip"] = "4.3.2.1"
	qParams["fingerprint"] = fingerprint1

	r = common.BuildRequest(http.MethodGet, "/v1/data/risk", nil, qParams)
	w = httptest.NewRecorder()

	server.GetRiskHandler(w, r)

	info := &FingerprintInfo{}
	json.Unmarshal(w.Body.Bytes(), info)
	assert.True(t, info.ImpossibleTravel)
	assert.False(t, info.Match.IpAddr)
	assert.False(t, info.Match.Location)
}

func TestHandlerRiskSameWhom(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(false, false), NewMockNods(false))

	qParams := make(map[string]string)
	qParams["email"] = "1@2.com"
	// location is South Korea
	qParams["ip"] = "1.2.3.4"
	qParams["ua"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
	qParams["hw"] = "01234567-89ABCDEF-01234567-89ABCDEF"
	qParams["pid"] = "osintami"

	r := common.BuildRequest(http.MethodGet, "/v1/data/fingerprint", nil, qParams)
	w := httptest.NewRecorder()

	server.GetFingerprintHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	fingerprint1 := w.Body.String()

	// change location to USA and run the print through NODS
	qParams["ip"] = "1.2.3.4"
	qParams["fingerprint"] = fingerprint1

	r = common.BuildRequest(http.MethodGet, "/v1/data/risk", nil, qParams)
	w = httptest.NewRecorder()

	server.GetRiskHandler(w, r)

	info := &FingerprintInfo{}
	json.Unmarshal(w.Body.Bytes(), info)
	assert.Equal(t, fingerprint1, info.Fingerprint)
	assert.False(t, info.ImpossibleTravel)
	assert.True(t, info.Match.IpAddr)
	assert.True(t, info.Match.Location)
}
