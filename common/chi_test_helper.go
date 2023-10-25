// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

func BuildRequest(method string, path string, pathParams map[string]string, queryParams map[string]string) *http.Request {

	body := strings.NewReader("")

	if queryParams != nil && method == http.MethodPost {
		content, _ := json.Marshal(queryParams)
		body = strings.NewReader(string(content))
	}

	r := httptest.NewRequest(
		method,
		path,
		body)

	if pathParams != nil {
		// chi requires any path parameters to be set in the route context
		rctx := chi.NewRouteContext()
		for k, v := range pathParams {
			rctx.URLParams.Add(k, v)
		}
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	}

	if queryParams != nil && method == http.MethodGet {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		r.URL.RawQuery = params.Encode()
	}
	return r
}

func BuildErrorResponse(err error) string {
	return fmt.Sprintf("%s%s%s", "{\"error\":\"", err.Error(), "\"}\n")
}
