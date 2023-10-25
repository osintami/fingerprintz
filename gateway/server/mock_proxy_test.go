// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/url"
)

type MockReverseProxy struct {
}

func NewMockReverseProxy() IProxy {
	return &ReverseProxy{}
}

func (x *MockReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, url *url.URL) {
}
