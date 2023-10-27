// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestDownloadHandler(t *testing.T) {
	svr := createServer(nil)
	r := common.BuildRequest(http.MethodGet, "/downloads?name=unwanted", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	w := httptest.NewRecorder()

	// TODO:  this file's location needs to be configurable
	os.WriteFile("unwanted.db", []byte(""), 0644)
	svr.DownloadHandler(w, r)

	// cleanup
	os.Remove("unwanted.db")
	os.Remove(svr.snapshotName("unwanted"))
}

func TestDownloadHandlerNonAdmin(t *testing.T) {
	svr := createServer(nil)
	r := common.BuildRequest(http.MethodGet, "/downloads?name=unwanted", nil, nil)
	w := httptest.NewRecorder()

	svr.DownloadHandler(w, r)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, common.BuildErrorResponse(ErrNotAuthorized), w.Body.String())
}

func TestDownloadHandlerNoDataFile(t *testing.T) {
	svr := createServer(nil)
	r := common.BuildRequest(http.MethodGet, "/downloads?name=unwanted", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	w := httptest.NewRecorder()

	svr.DownloadHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, common.BuildErrorResponse(ErrNotFound), w.Body.String())
}

func TestDownloadHandlerMissingName(t *testing.T) {
	svr := createServer(nil)
	r := common.BuildRequest(http.MethodGet, "/downloads", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	w := httptest.NewRecorder()

	svr.DownloadHandler(w, r)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Equal(t, common.BuildErrorResponse(ErrNotAuthorized), w.Body.String())
}

func TestDownloadHandlerFileNotModified(t *testing.T) {
	svr := createServer(nil)
	r := common.BuildRequest(http.MethodGet, "/downloads?name=unwanted", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	// NOTE:  this test may fail on a minute boundary
	r.Header.Add("If-Modified-Since", time.Now().Format(common.GO_DEFAULT_DATE))
	w := httptest.NewRecorder()

	// TODO:  this file's location needs to be configurable
	os.WriteFile("unwanted.db", []byte(""), 0644)
	svr.DownloadHandler(w, r)

	assert.Equal(t, http.StatusNotModified, w.Code)
	assert.Equal(t, common.BuildErrorResponse(ErrFileNotModified), w.Body.String())

	// cleanup
	os.Remove("unwanted.db")
}

func TestDownloadHandlerCorruptFile(t *testing.T) {
	fails := make(map[string]bool)
	fails["ToJSON"] = true
	svr := createServer(fails)
	r := common.BuildRequest(http.MethodGet, "/downloads?name=unwanted", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	w := httptest.NewRecorder()

	// TODO:  this file's location needs to be configurable
	os.WriteFile("unwanted.db", []byte(""), 0644)
	svr.DownloadHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// cleanup
	os.Remove("unwanted.db")
	os.Remove(svr.snapshotName("unwanted"))
}
