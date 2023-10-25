// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/mcnijman/go-emailaddress"
	"github.com/osintami/fingerprintz/common"
)

var ErrMissingOrInvalidEmail = errors.New("missing or invalid email")
var ErrMissingOrInvalidName = errors.New("missing or invalid name")
var ErrNotHuman = errors.New("captcha failed")
var ErrFakeDomain = errors.New("disposable email")
var ErrFingerprintSmudged = errors.New("fingerprint smudged")
var ErrIpAddrInUse = errors.New("ip address in use, please contact support@osintami.com")

func (x *GatewayServer) SignupHandler(w http.ResponseWriter, r *http.Request) {

	// check for AI
	ipAddr := common.IpAddr(r)
	keys := make(map[string]string)
	keys["ip"] = ipAddr
	value, err := x.oc.Item("rule/osintami/isHuman", keys)
	if err == nil && value.Result.Boolean() {
		common.SendError(w, ErrNotHuman, http.StatusTeapot)
		return
	}

	// check for valid email
	email := common.QueryParam(r, "email")
	if email == "" {
		common.SendError(w, ErrMissingOrInvalidEmail, http.StatusBadRequest)
		return
	}
	em, err := emailaddress.Parse(email)
	if err != nil {
		common.SendError(w, ErrMissingOrInvalidEmail, http.StatusBadRequest)
		return
	}

	// check for disposable email
	keys["domain"] = em.Domain
	value, err = x.oc.Item("domain/fakefilter/isFake", keys)
	if err == nil && value.Result.Boolean() {
		common.SendError(w, ErrFakeDomain, http.StatusBadRequest)
		return
	}

	// check for valid/sane name
	name := common.QueryParam(r, "name")
	if name == "" || len(name) > 32 {
		common.SendError(w, ErrMissingOrInvalidName, http.StatusBadRequest)
		return
	}

	// check for possible account fraud
	account, _ := x.accounts.FindByIpAddr(r.Context(), ipAddr)
	if account != nil {
		common.SendError(w, ErrIpAddrInUse, http.StatusBadRequest)
		return
	}

	// check for existing account
	account, _ = x.accounts.FindByEmail(r.Context(), email)
	if account != nil {
		common.SendError(w, ErrUserExists, http.StatusBadRequest)
		return
	} else {
		account := &Account{}
		account.Name = name
		account.Email = email
		account.Role = ROLE_USER
		account.Tokens = 0
		account.LastPayment = time.Now()
		account.StripeId = "free"
		account.Enabled = true
		account.IpAddr = ipAddr

		// fingerprint account
		// keys := make(map[string]string)
		// keys["eh"] = common.SHA256ToBase16Lowercase(email)
		// keys["ip"] = ipAddr
		// keys["ua"] = r.UserAgent()
		// fingerprint, err := x.oc.Fingerprint(keys)
		// if err != nil {
		// 	common.SendError(w, ErrFingerprintSmudged, http.StatusInternalServerError)
		// 	return
		// }
		// account.Fingerprint = fingerprint

		// create account
		err = x.accounts.CreateAccount(r.Context(), account)
		if err != nil {
			common.SendError(w, err, http.StatusInternalServerError)
			return
		}

		// send API key and instructions for free account
		go x.accounts.WelcomeEmail(account)

		// return eye candy
		whoami, _ := x.oc.Whoami(ipAddr)
		common.SendPrettyJSON(w, whoami)
	}
}
