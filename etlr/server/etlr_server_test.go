// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/etlr/etl"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func TestETLServer(t *testing.T) {
	server := NewETLrServer(NewMockETLManager())

	pParams := make(map[string]string)

	// refresh with valid source
	pParams["vendor"] = "test"
	r := common.BuildRequest(http.MethodGet, "/v1/refresh", pParams, nil)
	w := httptest.NewRecorder()
	server.RefreshHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"success\"}\n", w.Body.String())

	// refresh with the super secret and scary ALL source
	pParams["vendor"] = "ALL"
	r = common.BuildRequest(http.MethodGet, "/v1/refresh", pParams, nil)
	w = httptest.NewRecorder()
	server.RefreshHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"success\"}\n", w.Body.String())

	// refresh with invalid source
	pParams["vendor"] = "nope"
	r = common.BuildRequest(http.MethodGet, "/v1/refresh", pParams, nil)
	w = httptest.NewRecorder()
	server.RefreshHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"vendor not found\"}\n", w.Body.String())
}

type MockETLManager struct{}

func NewMockETLManager() etl.IETLManager {
	return &MockETLManager{}
}
func (x *MockETLManager) ScheduleCronJobs() *cron.Cron { return nil }
func (x *MockETLManager) RefreshAll()                  {}
func (x *MockETLManager) Refresh(vendorName string) error {
	if vendorName == "test" {
		return nil
	}
	return etl.ErrVendorNotFound
}
func (x *MockETLManager) FindJob(sourceName string) *etl.ETLJob { return &etl.ETLJob{} }
func (x *MockETLManager) Source(sourceName string) *etl.Source  { return &etl.Source{} }
