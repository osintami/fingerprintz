// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestServerConfig(t *testing.T) {
	cfg := &ServerConfig{}
	err := common.LoadEnv(true, false, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "/home/osintami/nods/", cfg.ConfigPath)
	assert.Equal(t, "/home/osintami/data/", cfg.SchemaPath)
	assert.Equal(t, "/home/osintami/data/", cfg.DataPath)
	assert.Equal(t, "/home/osintami/logs/", cfg.LogPath)
	assert.Equal(t, "/", cfg.PathPrefix)
	assert.Equal(t, "127.0.0.1:8082", cfg.ListenAddr)
	assert.Equal(t, "INFO", cfg.LogLevel)
}

func TestNodsServer(t *testing.T) {
	NewNormalizedDataServer(nil, nil, nil, nil)
}

func TestParamsFromGET(t *testing.T) {
	svr := NewNormalizedDataServer(nil, nil, nil, nil)

	ua := "-._~ — :/?#[]@ !$&'()*+,;= \\&*"

	pParams := make(map[string]string)
	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"
	qParams["phone"] = "18001234567"
	qParams["domain"] = "nope.com"
	qParams["browser"] = ua

	r := common.BuildRequest(http.MethodGet, "/v1/data/any", pParams, qParams)
	w := httptest.NewRecorder()
	params, err := svr.params.Keys(r)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, ua, params["browser"])
	assert.Equal(t, qParams["ip"], params["ip"])
	assert.Equal(t, qParams["domain"], params["domain"])
	assert.Equal(t, qParams["phone"], params["phone"])
}

func TestParamsFromPOST(t *testing.T) {
	svr := NewNormalizedDataServer(nil, nil, nil, nil)

	ua := "-._~ — :/?#[]@ !$&'()*+,;= \\&*"

	pParams := make(map[string]string)
	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"
	qParams["phone"] = "18001234567"
	qParams["domain"] = "nope.com"
	qParams["browser"] = ua

	r := common.BuildRequest(http.MethodPost, "/v1/data/any", pParams, qParams)
	w := httptest.NewRecorder()
	params, err := svr.params.Keys(r)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, ua, params["browser"])
	assert.Equal(t, qParams["ip"], params["ip"])
	assert.Equal(t, qParams["domain"], params["domain"])
	assert.Equal(t, qParams["phone"], params["phone"])
}

func TestParamsFromPOSTWithBadContent(t *testing.T) {
	svr := NewNormalizedDataServer(nil, nil, nil, nil)

	r := httptest.NewRequest(
		http.MethodPost,
		"/v1/data/any",
		strings.NewReader(string("...")))

	params, err := svr.params.Keys(r)
	assert.NotNil(t, err)
	assert.Nil(t, params)
}
