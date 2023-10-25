// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/osintami/fingerprintz/log"
	"github.com/stretchr/testify/assert"
)

func TestIpAddr(t *testing.T) {
	// default test IP
	r := httptest.NewRequest(http.MethodGet, "http:/127.0.0.1/", nil)
	ip := IpAddr(r)
	assert.Equal(t, "192.0.2.1", ip)

	// without port
	r = httptest.NewRequest(http.MethodGet, "http:/127.0.0.1/", nil)
	r.RemoteAddr = "1.2.3.4"
	ip = IpAddr(r)
	assert.Equal(t, "1.2.3.4", ip)

	// with port
	r = httptest.NewRequest(http.MethodGet, "http:/127.0.0.1/", nil)
	r.RemoteAddr = "1.2.3.4:1234"
	ip = IpAddr(r)
	assert.Equal(t, "1.2.3.4", ip)

	// special case IPv6
	r = httptest.NewRequest(http.MethodGet, "http:/127.0.0.1/", nil)
	r.RemoteAddr = "[::1]"
	ip = IpAddr(r)
	assert.Equal(t, "127.0.0.1", ip)
}

func TestInitZeroLog(t *testing.T) {
	err := InitZeroLog("TRACE")
	assert.Nil(t, err)

	err = InitZeroLog("NOPE")
	assert.NotNil(t, err)
}

func TestInitFileLog(t *testing.T) {
	// test defaults
	assert.Equal(t, log.LOG_LEVEL, log.LOG_ERROR)
	assert.True(t, log.LOG_STDERR)

	// use API to init file logger
	InitFileLog("/tmp/", "test.log", "TRACE", false)

	// test defaults are changed
	assert.Equal(t, log.LOG_TRACE, log.LOG_LEVEL)
	assert.False(t, log.LOG_STDERR)

	// cleanup
	os.Remove("/tmp/test.log")
}

type TestEnvConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"nope"`
}

func TestLoadEnv(t *testing.T) {
	cfg := &TestEnvConfig{}

	// create both personal and production .env files
	os.WriteFile(".env", []byte("LOG_LEVEL=bidness"), 0644)
	os.WriteFile(".env.personal", []byte("LOG_LEVEL=nunya"), 0644)

	// test that env.personal is loaded and not .env
	err := LoadEnv(true, true, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "nunya", cfg.LogLevel)

	// cleanup
	os.Remove(".env.personal")

	// test that production is loaded with .env.personal gone
	err = LoadEnv(true, true, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "bidness", cfg.LogLevel)

	// cleanup

	os.Remove(".env")
}

func TestLoadEnvEmtpy(t *testing.T) {
	cfg := &TestEnvConfig{}

	// create empty .env
	os.WriteFile(".env", []byte(""), 0644)

	// test that default is loaded
	err := LoadEnv(true, true, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "nope", cfg.LogLevel)

	// cleanup

	os.Remove(".env")
}
func TestLoadEnvOverride(t *testing.T) {
	cfg := &TestEnvConfig{}

	// test that OS variable already set overrides .env
	os.WriteFile(".env.personal", []byte("LOG_LEVEL=nunya"), 0644)
	os.Setenv("LOG_LEVEL", "osintami")

	err := LoadEnv(false, true, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "osintami", cfg.LogLevel)

	// cleanup
	os.Remove(".env.personal")
}

func TestLoadEnvMissingFile(t *testing.T) {
	cfg := &TestEnvConfig{}

	// error paths
	err := LoadEnv(true, true, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "nope", cfg.LogLevel)
}

func TestLoadEnvBadPointer(t *testing.T) {
	err := LoadEnv(true, false, TestEnvConfig{})
	assert.NotNil(t, err)
}

func TestLoadEnvCorrupFile(t *testing.T) {
	cfg := &TestEnvConfig{}

	os.WriteFile(".env.personal", []byte("fj;lk1450^780y2 @!34 46."), 0644)
	err := LoadEnv(true, true, cfg)
	assert.Nil(t, err)

	// cleanup
	os.Remove(".env.personal")
}

func TestLoadJson(t *testing.T) {

	type TestJsonConfig struct {
		A string
		B bool
		C int
	}
	cfg := &TestJsonConfig{}

	// test missing file
	err := LoadJson("/tmp/x.json", cfg)
	assert.NotNil(t, err)

	// create and test valid JSON data for Config
	os.WriteFile("/tmp/test.json", []byte("{\"A\":\"nunya\",\"B\":true,\"C\":82}"), 0644)

	err = LoadJson("/tmp/test.json", cfg)
	assert.Nil(t, err)
	assert.Equal(t, "nunya", cfg.A)
	assert.True(t, cfg.B)
	assert.Equal(t, 82, cfg.C)

	// create and test invalid JSON data for Config
	os.WriteFile("/tmp/test.json", []byte("."), 0644)

	err = LoadJson("/tmp/test.json", cfg)
	assert.NotNil(t, err)

	// cleanup
	os.Remove("/tmp/test.json")
}

func TestPrintEnvironment(t *testing.T) {
	os.Setenv("TEST_API_KEY", "xxx")
	PrintEnvironment()
}

func TestQueryParam(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "https://api.osintami.com/data/items", strings.NewReader(""))
	assert.Nil(t, err)

	params := url.Values{}
	params.Add("key", "xxx")
	r.URL.RawQuery = params.Encode()

	key := QueryParam(r, "key")
	assert.Equal(t, "xxx", key)
}

func TestPathParam(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "https://api.osintami.com/data/items", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("key", "xxx")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	key := PathParam(r, "key")
	assert.Equal(t, "xxx", key)
}

func TestSendError(t *testing.T) {
	err := errors.New("test error")
	w := httptest.NewRecorder()

	SendError(w, err, http.StatusInternalServerError)
	res := w.Result()

	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"error\":\"test error\"}\n", string(data))
}

func TestSendPrettyJSON(t *testing.T) {
	w := httptest.NewRecorder()

	type ResponseInfo struct {
		Message string
		Value   int
	}
	responseInfo := &ResponseInfo{Message: "nunya", Value: 82}

	SendPrettyJSON(w, responseInfo)
	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\n    \"Message\": \"nunya\",\n    \"Value\": 82\n}\n", string(data))
}

func TestSendJSON(t *testing.T) {
	w := httptest.NewRecorder()

	type ResponseInfo struct {
		Message string
		Value   int
	}
	responseInfo := &ResponseInfo{Message: "nunya", Value: 82}

	SendJSON(w, responseInfo)
	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"Message\":\"nunya\",\"Value\":82}\n", string(data))
}

func TestListenAndServe(t *testing.T) {
	err := ListenAndServe("0.0.0.0", "", "", nil)
	assert.NotNil(t, err)

	err = ListenAndServe("0.0.0.0", "cert", "key", nil)
	assert.NotNil(t, err)
}
