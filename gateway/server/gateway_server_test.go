// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func createServer() *GatewayServer {
	failMap := make(map[string]bool)
	return NewGatewayServer(
		NewMockReverseProxy(),
		NewMockNods(failMap),
		common.NewPersistentCache("/tmp/unwanted.db"),
		NewMockAccounts(failMap),
		NewMockCalls(failMap),
		NewMockPixels(failMap),
		NewMockGatewayServices())
}

func TestReverseProxyHandler(t *testing.T) {
	svr := createServer()

	r := common.BuildRequest(http.MethodGet, "/data/schema", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	w := httptest.NewRecorder()

	svr.ReverseProxyHandler(w, r)
}

func TestReverseProxyHandlerJackedConfig(t *testing.T) {
	svr := createServer()

	r := common.BuildRequest(http.MethodGet, "/corrupt/list", nil, nil)
	r.Header.Add("X-Api-Key", "user_api_key")
	w := httptest.NewRecorder()

	svr.ReverseProxyHandler(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestReverseProxyHandlerUnsupportedRoute(t *testing.T) {
	svr := createServer()

	r := common.BuildRequest(http.MethodGet, "/nope", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	w := httptest.NewRecorder()

	svr.ReverseProxyHandler(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestReverseProxyHandlerUnknownUser(t *testing.T) {
	svr := createServer()

	r := common.BuildRequest(http.MethodGet, "/data/schema", nil, nil)
	w := httptest.NewRecorder()

	svr.ReverseProxyHandler(w, r)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestReverseProxyHandlerAdminRequired(t *testing.T) {
	svr := createServer()

	r := common.BuildRequest(http.MethodGet, "/admin/list", nil, nil)
	r.Header.Add("X-Api-Key", "user_api_key")
	w := httptest.NewRecorder()

	svr.ReverseProxyHandler(w, r)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestBurnToken(t *testing.T) {
	svr := createServer()

	// valid account
	r := common.BuildRequest(http.MethodGet, "/data/schema", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	account, err := svr.authUser(r, "admin_api_key", "nods")
	assert.Nil(t, err)
	assert.Equal(t, 0, account.Tokens)

	// invalid account
	r = common.BuildRequest(http.MethodGet, "/data/schema", nil, nil)
	r.Header.Add("X-Api-Key", "nope_api_key")
	account, err = svr.authUser(r, "nope_api_key", "nods")
	assert.Equal(t, ErrApiKey, err)
	assert.Nil(t, account)
}

func TestAuditTrail(t *testing.T) {

	failMap := make(map[string]bool)
	failMap["Call"] = true

	svr := NewGatewayServer(
		NewMockReverseProxy(),
		NewMockNods(failMap),
		common.NewPersistentCache("/tmp/unwanted.db"),
		NewMockAccounts(failMap),
		NewMockCalls(failMap),
		NewMockPixels(failMap),
		NewMockGatewayServices())

	account := &Account{
		Name:        "test-name",
		Email:       "test@example.com",
		ApiKey:      "",
		Tokens:      0,
		Role:        "admin",
		LastPayment: time.Now(),
		Enabled:     true,
		StripeId:    "test-stripe-id",
	}

	call := &Call{
		AccountId: account.ID,
		RequestId: "xxx",
		IpAddr:    "1.2.3.4",
		API:       "nods",
		Page:      ""}

	svr.auditTrail(account, call)

	failMap["BurnToken"] = true
	svr.auditTrail(account, call)
}

func TestFindAccount(t *testing.T) {
	svr := createServer()

	// admin account
	r := common.BuildRequest(http.MethodGet, "/data/schema", nil, nil)
	r.Header.Add("X-Api-Key", "admin_api_key")
	account, err := svr.findAccount(r)
	assert.Nil(t, err)
	assert.Equal(t, "admin", account.Role)

	// user account
	r = common.BuildRequest(http.MethodGet, "/data/schema", nil, nil)
	r.Header.Add("X-Api-Key", "user_api_key")
	account, err = svr.findAccount(r)
	assert.Nil(t, err)
	assert.Equal(t, "user", account.Role)

	// invalid account
	r = common.BuildRequest(http.MethodGet, "/data/schema", nil, nil)
	r.Header.Add("X-Api-Key", "nope_api_key")
	account, err = svr.findAccount(r)
	assert.Equal(t, ErrApiKey, err)
	assert.Nil(t, account)

	// missing API key
	r = common.BuildRequest(http.MethodGet, "/data/schema", nil, nil)
	r.Header.Add("X-Api-Key", "")
	account, err = svr.findAccount(r)
	assert.Equal(t, ErrApiKey, err)
	assert.Nil(t, account)
}

func TestFindAccountBadKey(t *testing.T) {
	svr := createServer()
	r := common.BuildRequest(http.MethodGet, "/data/schema", nil, nil)
	r.Header.Add("X-Api-Key", "nope_api_key")
	account, err := svr.findAccount(r)
	assert.NotNil(t, err)
	assert.Nil(t, account)
}
