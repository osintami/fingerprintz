// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/stripe/stripe-go"
)

var WEBHOOK_WHITELIST = []string{
	"3.18.12.63",
	"3.130.192.231",
	"13.235.14.237",
	"13.235.122.149",
	"18.211.135.69",
	"35.154.171.200",
	"52.15.183.38",
	"54.88.130.119",
	"54.88.130.237",
	"54.187.174.169",
	"54.187.205.235",
	"54.187.216.72"}

func (x *GatewayServer) StripeWhitelistCheck(r *http.Request) bool {

	keys := make(map[string]string)
	incomingIP := common.IpAddr(r)
	keys["ip"] = incomingIP
	// OSINTAMI updates Stripe webhook IPs once per week
	data, err := x.oc.Item("ip/stripe.webhooks/isStripe", keys)
	if err != nil {
		// on error from OSINTAMI, use the static list of webhook IPs
		log.Error().Err(err).Str("component", "stripe").Str("ip", incomingIP).Msg("OSINTAMI lookup failed")
		for _, ip := range WEBHOOK_WHITELIST {
			if ip == incomingIP {
				return true
			}
		}
		return false
	}
	return data.Result.Boolean()
}

func (x *GatewayServer) StripeHandler(w http.ResponseWriter, r *http.Request) {

	if !x.StripeWhitelistCheck(r) {
		common.SendError(w, ErrNotAuthorized, http.StatusForbidden)
		return
	}

	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		log.Error().Err(err).Str("component", "stripe").Str("event", "unknown").Msg("read POST body failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Error().Err(err).Str("component", "stripe").Str("event", "unknown").Msg("unmarshal Stripe event failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ------------------------------------------------ AUDIT/DEBUG ---------------------------------------------
	// Stripe audit/debug /home/osintami/gateway/stripe-{event name}-{timestamp}.json
	data, _ := json.MarshalIndent(event, "", "  ")
	date := time.Now().Format("2006-01-02 15:04:05")
	os.WriteFile(fmt.Sprintf("/home/osintami/logs/stripe-event-%s-%s.json", event.Type, date), data, 0644)
	//------------------------------------------------- AUDIT/DEBUG ---------------------------------------------

	switch event.Type {

	case "customer.subscription.deleted":
		subscription := &stripe.Subscription{}
		// NOTE: entire event has already been parsed
		json.Unmarshal(event.Data.Raw, subscription)
		log.Info().Str("component", "stripe").Str("customer", subscription.Customer.ID).Msg("subscription cancelled")

		stripeId := subscription.Customer.ID
		account, err := x.accounts.FindByStripeId(r.Context(), stripeId)
		if err == nil {
			account.Enabled = false
			x.accounts.EnableOrDisableAccount(r.Context(), account)
			w.WriteHeader(http.StatusOK)
			return
		}

	case "charge.succeeded":
		charge := &stripe.Charge{}
		// NOTE: entire event has already been parsed
		json.Unmarshal(event.Data.Raw, charge)
		log.Info().Str("component", "stripe").Str("customer", charge.Customer.ID).Msg("payment received")

		email := charge.BillingDetails.Email
		account, err := x.accounts.FindByEmail(r.Context(), email)
		if err == nil {
			account.LastPayment = time.UnixMilli(charge.Created)
			x.accounts.UpdateLastPayment(r.Context(), account)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			account := &Account{}
			account.Name = charge.BillingDetails.Name
			account.Email = email
			account.ApiKey = ""
			account.Role = ROLE_USER
			account.Tokens = 0
			account.LastPayment = time.UnixMilli(charge.Created)
			account.StripeId = charge.Customer.ID

			err := x.accounts.CreateAccount(r.Context(), account)
			if err != nil {
				log.Error().Err(err).Str("component", "stripe").Msg("create user failed")
			}
			go x.accounts.WelcomeEmail(account)
			w.WriteHeader(http.StatusOK)
			return

		}
	default:
		log.Error().Str("component", "stripe").Str("event", event.Type).Msg("unhandled event")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
}
