// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestPixelFireHandler(t *testing.T) {
	svr := createServer()
	r := common.BuildRequest(http.MethodGet, "/images/any.png", nil, nil)

	r.Header.Add("User-Agent", "test user agent")

	w := httptest.NewRecorder()
	svr.PixelFireHandler(w, r)

	cookie := w.Result().Cookies()[0]
	assert.Equal(t, "osintami.com", cookie.Domain)
	assert.Equal(t, "osintami", cookie.Name)
	assert.Equal(t, 36, len(cookie.Value))

	cookie.Value = "nope_cookie_id"
	r.AddCookie(cookie)

	w = httptest.NewRecorder()
	svr.PixelFireHandler(w, r)
	cookie = w.Result().Cookies()[0]
	assert.Equal(t, "osintami.com", cookie.Domain)
	assert.Equal(t, "osintami", cookie.Name)
	assert.Equal(t, "nope_cookie_id", cookie.Value)
}
