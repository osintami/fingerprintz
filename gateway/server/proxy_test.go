// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osintami/fingerprintz/common"
)

func TestProxy(t *testing.T) {
	proxy := NewReverseProxy()
	r := common.BuildRequest(http.MethodGet, "/data/items", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, r, r.URL)
}
