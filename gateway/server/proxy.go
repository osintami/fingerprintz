// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type IProxy interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, url *url.URL)
}

type ReverseProxy struct {
}

func NewReverseProxy() IProxy {
	return &ReverseProxy{}
}

func (x *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, url *url.URL) {
	p := httputil.NewSingleHostReverseProxy(url)
	p.Director = func(request *http.Request) {
		request.Host = url.Host
		request.URL.Scheme = url.Scheme
		request.URL.Host = url.Host
		request.URL.Path = url.Path
	}
	p.ServeHTTP(w, r)
}
