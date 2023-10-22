// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestHandlerFingerprintNoInputs(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(false, false), NewMockNods(false), "")

	qParams := make(map[string]string)

	r := common.BuildRequest(http.MethodPost, "/v1/data/fingerprint", nil, qParams)
	w := httptest.NewRecorder()

	server.PostFingerprintHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerFingerprint(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(false, false), NewMockNods(false), "")

	qParams := make(map[string]string)
	qParams["email"] = "1@2.com"
	qParams["ip"] = "1.2.3.4"
	qParams["ua"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
	qParams["hw"] = "01234567-89ABCDEF-01234567-89ABCDEF"
	qParams["pid"] = "osintami"

	r := common.BuildRequest(http.MethodGet, "/v1/data/fingerprint", nil, qParams)
	w := httptest.NewRecorder()

	server.GetFingerprintHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	signature := w.Body.String()

	fingerprint, err := server.signer.DecodeJWT(signature)
	assert.Nil(t, err)
	assert.Equal(t, common.SHA256ToBase16Lowercase(qParams["email"]), fingerprint.EHash)
	assert.Equal(t, qParams["ip"], fingerprint.IpAddr)
	assert.Equal(t, qParams["ua"], fingerprint.UserAgent)
	assert.Equal(t, qParams["hw"], fingerprint.DeviceId)
	assert.Equal(t, fmt.Sprintf("0X%04X", BIT_BLACKLIST|BIT_CLOUD_NODE|BIT_KNOWN_BOT|BIT_PROXY_IN_USE|BIT_TOR_IN_USE|BIT_VPN_IN_USE), fingerprint.NetworkId)
	assert.Equal(t, qParams["pid"], fingerprint.PartnerId)
}

func TestHandlerFingerprintSignFailure(t *testing.T) {
	server := NewWhoamiServer(NewMockJWTSigner(false, true), NewMockNods(false), "")
	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	r := common.BuildRequest(http.MethodGet, "/v1/data/fingerprint", nil, qParams)
	w := httptest.NewRecorder()

	server.GetFingerprintHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
