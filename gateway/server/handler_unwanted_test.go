// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestUnwantedHandler(t *testing.T) {
	svr := createServer()
	r := common.BuildRequest(http.MethodGet, "/", nil, nil)
	w := httptest.NewRecorder()
	svr.UnwantedHandler(w, r)
	assert.Equal(t, http.StatusPaymentRequired, w.Code)

	// check cache of unwanted visitors
	unwantedInfo, found := svr.cache.Get("192.0.2.1")
	assert.True(t, found)
	assert.Equal(t, 1, unwantedInfo.(*UnwantedInfo).Count)

	// increment the unwanted count
	w = httptest.NewRecorder()
	svr.UnwantedHandler(w, r)
	assert.Equal(t, http.StatusPaymentRequired, w.Code)

	unwantedInfo, found = svr.cache.Get("192.0.2.1")
	assert.True(t, found)
	assert.Equal(t, 2, unwantedInfo.(*UnwantedInfo).Count)
}
