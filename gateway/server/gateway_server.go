// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type GatewayServer struct {
	proxy    IProxy
	accounts IAccounts
	calls    ICalls
	pixels   IPixels
	cache    common.IFastCache
	nods     common.INods
	services IGatewayServices
}

func NewGatewayServer(proxy IProxy, nods common.INods, cache common.IFastCache, accounts IAccounts, calls ICalls, pixels IPixels, gws IGatewayServices) *GatewayServer {

	return &GatewayServer{
		proxy:    proxy,
		nods:     nods,
		cache:    cache,
		accounts: accounts,
		calls:    calls,
		pixels:   pixels,
		services: gws}
}

func (x *GatewayServer) ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	// identify the service from the URL
	service := x.services.find(r.URL.Path)
	if service == nil {
		// TODO:  log page so we can learn who's being bad and how (addr, user-agent, page, count)
		common.SendError(w, ErrNotFound, http.StatusNotFound)
		return
	}

	// generate a request identifier for logging calls
	requestId := uuid.NewString()

	// find the user from the API key
	if !service.IsPublic() {

		user, err := x.authUser(r, requestId, service.Name)
		if err != nil {
			log.Error().Err(err).Str("component", "gateway").Str("service", service.Name).Msg("token burn")
			common.SendError(w, ErrNotAuthorized, http.StatusForbidden)
			return
		}

		// TODO:  when and if we start counting tokens this is the spot

		// basic user vs admin ACL checking for this API
		if service.WantsAdmin() && !user.IsAdmin() {
			log.Error().Err(err).Str("component", "gateway").Str("service", service.Name).Msg("admin required")
			common.SendError(w, ErrNotAuthorized, http.StatusForbidden)
			return
		}

		values := r.URL.Query()
		// strip out security concerns before forwarding
		values.Del("key")
		values.Del("role")

		role := service.ACL
		// bump ACL's up for the service based on user ACL vs service minimum ACL
		if user.IsAdmin() {
			role = ROLE_ADMIN
		}
		values.Add("role", role)
		r.URL.RawQuery = values.Encode()
	}

	path := fmt.Sprintf("%s%s", service.Service, strings.ReplaceAll(r.URL.Path, service.Route, service.Path))
	if url, err := url.Parse(path); err != nil {
		log.Error().Err(err).Str("component", "gateway").Str("service", service.Name).Msg("parse path")
		common.SendError(w, ErrBadRequest, http.StatusInternalServerError)
		return
	} else {
		x.proxy.ServeHTTP(w, r, url)
	}
}

func (x *GatewayServer) authUser(r *http.Request, requestId string, svc string) (*Account, error) {
	if account, err := x.findAccount(r); err == nil {

		// don't log the API key, it's already been verified to match the account
		values := r.URL.Query()
		values.Del("key")
		call := &Call{
			AccountId: account.ID,
			RequestId: requestId,
			IpAddr:    common.IpAddr(r),
			API:       svc,
			Page:      r.URL.Path + "?" + values.Encode()}

		go x.auditTrail(account, call)
		return account, nil
	}
	return nil, ErrApiKey
}

func (x *GatewayServer) auditTrail(account *Account, call *Call) error {

	err := x.accounts.BurnToken(context.Background(), account)
	if err != nil {
		log.Error().Err(err).Str("audit", call.API).Str("key", account.ApiKey).Msg("burn token")
	}
	err = x.calls.Call(context.Background(), call)
	if err != nil {
		log.Error().Err(err).Str("audit", call.API).Str("key", account.ApiKey).Msg("log call")
	}
	return err
}

func (x *GatewayServer) findAccount(r *http.Request) (*Account, error) {

	// grab the auth token from the header
	key := r.URL.Query().Get("key")

	// maybe they stuck it in the query params, we support both to make things easy as possible
	if key == "" {
		key = r.Header.Get("X-Api-Key")
	}

	if key == "" {
		x.logUnwantedGuests(r)
		return nil, ErrApiKey
	} else {
		account, err := x.accounts.FindByApiKey(r.Context(), key)
		if err != nil || key != account.ApiKey || !account.Enabled {
			log.Error().Str("component", "gateway").Str("api-key", key).Msg("account not found")
			return nil, ErrApiKey
		}
		return account, nil
	}
}

// func (x *GatewayServer) checkWhitelist(h http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		incomingIP := common.IpAddr(r)
// 		account, err := x.findAccount(r)
// 		if err != nil {
// 			common.SendError(w, ErrNotAuthorized, http.StatusTeapot)
// 			return
// 		}
// 		// TODO:  put whitelisted IPs in the account object
// 		if account.isWhitelistedIP(incomingIP) {
// 			err := errors.New("address not whitelisted")
// 			log.Error().Err(err).Str("component", "whitelist").Msg("IP address not whitelisted")
// 			common.SendError(w, err, http.StatusForbidden)
// 			return
// 		}
// 		h.ServeHTTP(w, r)
// 	}

// 	return http.HandlerFunc(fn)
// }
