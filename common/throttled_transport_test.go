// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThrottledTransport(t *testing.T) {
	// NOTE:  set this sufficiently high if you are debugging
	transport := NewThrottledTransport(5*time.Second, 1, http.DefaultTransport)

	r := httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		strings.NewReader("{}"))

	// NOTE:  this returns an error because we can't mock a remote server with httpmock
	//  because httpmock overrides the DefaultTransport too
	resp, err := transport.RoundTrip(r)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// NOTE:  use the mock context to force a timeout
	r = httptest.NewRequest(
		http.MethodGet,
		"http://localhost",
		strings.NewReader("{}")).WithContext(NewMockContext())

	resp, err = transport.RoundTrip(r)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}
