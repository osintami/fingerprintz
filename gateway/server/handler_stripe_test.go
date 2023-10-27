// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go"
)

func TestStripeInWhitelist(t *testing.T) {
	svr := createServer(nil)
	qParams := make(map[string]string)
	qParams["stripe"] = "{\"payload\":true}"
	r := common.BuildRequest(http.MethodPost, "/v1/stripe", nil, qParams)
	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "admin_api_key")
	stripe := svr.StripeWhitelistCheck(r)
	assert.True(t, stripe)
}

func TestStripeNotInWhitelist(t *testing.T) {
	svr := createServer(nil)
	qParams := make(map[string]string)
	qParams["stripe"] = "{\"payload\":true}"
	r := common.BuildRequest(http.MethodPost, "/v1/stripe", nil, qParams)
	r.RemoteAddr = "1.2.3.4"
	r.Header.Add("X-Api-Key", "admin_api_key")
	stripe := svr.StripeWhitelistCheck(r)
	assert.False(t, stripe)
}

func TestStripeWhitelistNodsErrorInHardcodedWhitelist(t *testing.T) {
	svr := createServer(nil)
	table := make(map[string]bool)
	table["Item"] = true
	svr.oc = NewMockNods(table)
	qParams := make(map[string]string)
	qParams["stripe"] = "{\"payload\":true}"
	r := common.BuildRequest(http.MethodPost, "/v1/stripe", nil, qParams)
	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "admin_api_key")
	stripe := svr.StripeWhitelistCheck(r)
	assert.True(t, stripe)
}

func TestStripeWhitelistNodsErrorNotInHardcodedWhitelist(t *testing.T) {
	svr := createServer(nil)
	table := make(map[string]bool)
	table["Item"] = true
	svr.oc = NewMockNods(table)
	qParams := make(map[string]string)
	qParams["stripe"] = "{\"payload\":true}"
	r := common.BuildRequest(http.MethodPost, "/v1/stripe", nil, qParams)
	r.RemoteAddr = "1.2.3.4"
	r.Header.Add("X-Api-Key", "admin_api_key")
	stripe := svr.StripeWhitelistCheck(r)
	assert.False(t, stripe)
}

func TestStripeHandlerJunkEvent(t *testing.T) {
	svr := createServer(nil)
	qParams := make(map[string]string)
	qParams["stripe"] = "{\"payload\":true}"
	r := common.BuildRequest(http.MethodPost, "/v1/stripe", nil, qParams)
	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "admin_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestStripeHandlerNotInWhitelist(t *testing.T) {
	svr := createServer(nil)
	qParams := make(map[string]string)
	qParams["stripe"] = "{\"payload\":true}"
	r := common.BuildRequest(http.MethodPost, "/v1/stripe", nil, qParams)
	r.RemoteAddr = "1.2.3.4"
	r.Header.Add("X-Api-Key", "admin_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStripeHandlerCorruptEvent(t *testing.T) {
	svr := createServer(nil)
	r := httptest.NewRequest(
		http.MethodPost,
		"/v1/stripe",
		strings.NewReader(string("...")))

	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "admin_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStripeHandlerEmptyEvent(t *testing.T) {
	svr := createServer(nil)

	r := httptest.NewRequest(
		http.MethodPost,
		"/v1/stripe",
		nil)

	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "admin_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStripeChargeSucceededEvent(t *testing.T) {
	svr := createServer(nil)
	r := httptest.NewRequest(
		http.MethodPost,
		"/v1/stripe",
		strings.NewReader(string(CHARGE_SUCEEDED_EVENT)))

	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "admin_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStripeChargeSucceededCreateAccountFail(t *testing.T) {
	svr := createServer(nil)

	table := make(map[string]bool)
	table["CreateAccount"] = true
	table["FindByEmail"] = true
	svr.accounts = NewMockAccounts(table)

	r := httptest.NewRequest(
		http.MethodPost,
		"/v1/stripe",
		strings.NewReader(string(CHARGE_SUCEEDED_EVENT)))

	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "admin_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStripeKnownUser(t *testing.T) {
	svr := createServer(nil)
	r := httptest.NewRequest(
		http.MethodPost,
		"/v1/stripe",
		strings.NewReader(string(CHARGE_SUCEEDED_EVENT)))

	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "nope_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStripeUnknownUser(t *testing.T) {
	svr := createServer(nil)
	failMap := make(map[string]bool)
	failMap["FindByEmail"] = true
	svr.accounts = NewMockAccounts(failMap)
	r := httptest.NewRequest(
		http.MethodPost,
		"/v1/stripe",
		strings.NewReader(string(CHARGE_SUCEEDED_EVENT)))

	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "nope_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStripeSubscriptionCanceledEvent(t *testing.T) {
	svr := createServer(nil)
	r := httptest.NewRequest(
		http.MethodPost,
		"/v1/stripe",
		strings.NewReader(string(SUBSCRIPTION_CANCELLED_EVENT)))

	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "admin_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStripeUnknownEvent(t *testing.T) {
	svr := createServer(nil)
	r := httptest.NewRequest(
		http.MethodPost,
		"/v1/stripe",
		strings.NewReader(string(UNKNOWN_EVENT)))

	r.RemoteAddr = "3.18.12.63"
	r.Header.Add("X-Api-Key", "admin_api_key")

	w := httptest.NewRecorder()
	svr.StripeHandler(w, r)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestStripeCustomerId(t *testing.T) {
	// customer ID from cancellation
	event := stripe.Event{}
	json.Unmarshal([]byte(SUBSCRIPTION_CANCELLED_EVENT), &event)
	subscription := &stripe.Subscription{}
	json.Unmarshal(event.Data.Raw, subscription)
	assert.Equal(t, "cus_OhoLtcMC8B1Eu6", subscription.Customer.ID)

	// customer ID from payment
	event = stripe.Event{}
	json.Unmarshal([]byte(CHARGE_SUCEEDED_EVENT), &event)
	charge := &stripe.Charge{}
	json.Unmarshal(event.Data.Raw, charge)
	assert.Equal(t, "cus_OhoLtcMC8B1Eu6", charge.Customer.ID)
}

var UNKNOWN_EVENT = `
{
	"account": "",
	"created": 1695686526,
	"data": {
	  "previous_attributes": null,
		"object": {
		"id": "sub_1NuOj3GP5WOSV9CluFBGjyiO",
		"object": "nope"
		}
	}
}
`

var SUBSCRIPTION_CANCELLED_EVENT = `
{
	"account": "",
	"created": 1695686526,
	"data": {
	  "previous_attributes": null,
		"object": {
		"id": "sub_1NuOj3GP5WOSV9CluFBGjyiO",
		"object": "subscription",
		"application": null,
		"application_fee_percent": null,
		"automatic_tax": {
			"enabled": false
		},
		"billing_cycle_anchor": 1695686525,
		"billing_thresholds": null,
		"cancel_at": null,
		"cancel_at_period_end": false,
		"canceled_at": 1695687331,
		"cancellation_details": {
			"comment": null,
			"feedback": null,
			"reason": "cancellation_requested"
		},
		"collection_method": "charge_automatically",
		"created": 1695686525,
		"currency": "mxn",
		"current_period_end": 1698278525,
		"current_period_start": 1695686525,
		"customer": "cus_OhoLtcMC8B1Eu6",
		"days_until_due": null,
		"default_payment_method": "pm_1NuOj2GP5WOSV9ClXJ6YBY4e",
		"default_source": null,
		"default_tax_rates": [
		],
		"description": null,
		"discount": null,
		"ended_at": 1695687331,
		"items": {
			"object": "list",
			"data": [
			{
				"id": "si_OhoLdSh7753s7w",
				"object": "subscription_item",
				"billing_thresholds": null,
				"created": 1695686525,
				"metadata": {
				},
				"plan": {
				"id": "price_1NoaqPGP5WOSV9ClNjVEiyXO",
				"object": "plan",
				"active": true,
				"aggregate_usage": null,
				"amount": 8888,
				"amount_decimal": "8888",
				"billing_scheme": "per_unit",
				"created": 1694303141,
				"currency": "mxn",
				"interval": "month",
				"interval_count": 1,
				"livemode": false,
				"metadata": {
				},
				"nickname": null,
				"product": "prod_OboTFC6GMzaCTr",
				"tiers_mode": null,
				"transform_usage": null,
				"trial_period_days": null,
				"usage_type": "licensed"
				},
				"price": {
				"id": "price_1NoaqPGP5WOSV9ClNjVEiyXO",
				"object": "price",
				"active": true,
				"billing_scheme": "per_unit",
				"created": 1694303141,
				"currency": "mxn",
				"custom_unit_amount": null,
				"livemode": false,
				"lookup_key": null,
				"metadata": {
				},
				"nickname": null,
				"product": "prod_OboTFC6GMzaCTr",
				"recurring": {
					"aggregate_usage": null,
					"interval": "month",
					"interval_count": 1,
					"trial_period_days": null,
					"usage_type": "licensed"
				},
				"tax_behavior": "unspecified",
				"tiers_mode": null,
				"transform_quantity": null,
				"type": "recurring",
				"unit_amount": 8888,
				"unit_amount_decimal": "8888"
				},
				"quantity": 1,
				"subscription": "sub_1NuOj3GP5WOSV9CluFBGjyiO",
				"tax_rates": [
				]
			}
			],
			"has_more": false,
			"total_count": 1,
			"url": "/v1/subscription_items?subscription=sub_1NuOj3GP5WOSV9CluFBGjyiO"
		},
		"latest_invoice": "in_1NuOj3GP5WOSV9ClfpSHEJz9",
		"livemode": false,
		"metadata": {
		},
		"next_pending_invoice_item_invoice": null,
		"on_behalf_of": null,
		"pause_collection": null,
		"payment_settings": {
			"payment_method_options": null,
			"payment_method_types": null,
			"save_default_payment_method": "off"
		},
		"pending_invoice_item_interval": null,
		"pending_setup_intent": null,
		"pending_update": null,
		"plan": {
			"id": "price_1NoaqPGP5WOSV9ClNjVEiyXO",
			"object": "plan",
			"active": true,
			"aggregate_usage": null,
			"amount": 8888,
			"amount_decimal": "8888",
			"billing_scheme": "per_unit",
			"created": 1694303141,
			"currency": "mxn",
			"interval": "month",
			"interval_count": 1,
			"livemode": false,
			"metadata": {
			},
			"nickname": null,
			"product": "prod_OboTFC6GMzaCTr",
			"tiers_mode": null,
			"transform_usage": null,
			"trial_period_days": null,
			"usage_type": "licensed"
		},
		"quantity": 1,
		"schedule": null,
		"start_date": 1695686525,
		"status": "canceled",
		"test_clock": null,
		"transfer_data": null,
		"trial_end": null,
		"trial_settings": {
			"end_behavior": {
			"missing_payment_method": "create_invoice"
			}
		},
		"trial_start": null
		}
	},
	"id": "evt_3NuOj3GP5WOSV9Cl0uMdbs2Y",
	"livemode": false,
	"pending_webhooks": 1,
	"request": {
	  "id": "req_0tViCOPOdxXtyC",
	  "idempotency_key": "1c54f60f-7cca-4bbc-9f85-2cc61aac72f2"
	},
	"type": "customer.subscription.deleted"
  }`

var CHARGE_SUCEEDED_EVENT = `
{
	"account": "",
	"created": 1695686526,
	"data": {
	  "previous_attributes": null,
	  "object": {
		"id": "ch_3NuOj3GP5WOSV9Cl0moVNjda",
		"object": "charge",
		"amount": 8888,
		"amount_captured": 8888,
		"amount_refunded": 0,
		"application": null,
		"application_fee": null,
		"application_fee_amount": null,
		"balance_transaction": "txn_3NuOj3GP5WOSV9Cl0En1E3nT",
		"billing_details": {
		  "address": {
			"city": null,
			"country": "MX",
			"line1": null,
			"line2": null,
			"postal_code": null,
			"state": null
		  },
		  "email": "test@example.com",
		  "name": "Joe Bob",
		  "phone": null
		},
		"calculated_statement_descriptor": "WWW.OSINTAMI.COM",
		"captured": true,
		"created": 1695686525,
		"currency": "mxn",
		"customer": "cus_OhoLtcMC8B1Eu6",
		"description": "Subscription creation",
		"destination": null,
		"dispute": null,
		"disputed": false,
		"failure_balance_transaction": null,
		"failure_code": null,
		"failure_message": null,
		"fraud_details": {},
		"invoice": "in_1NuOj3GP5WOSV9ClfpSHEJz9",
		"livemode": false,
		"metadata": {},
		"on_behalf_of": null,
		"order": null,
		"outcome": {
		  "network_status": "approved_by_network",
		  "reason": null,
		  "risk_level": "normal",
		  "risk_score": 20,
		  "seller_message": "Payment complete.",
		  "type": "authorized"
		},
		"paid": true,
		"payment_intent": "pi_3NuOj3GP5WOSV9Cl0w8rGXeT",
		"payment_method": "pm_1NuOj2GP5WOSV9ClXJ6YBY4e",
		"payment_method_details": {
		  "card": {
			"amount_authorized": 8888,
			"brand": "visa",
			"checks": {
			  "address_line1_check": null,
			  "address_postal_code_check": null,
			  "cvc_check": "pass"
			},
			"country": "US",
			"exp_month": 5,
			"exp_year": 2024,
			"extended_authorization": {
			  "status": "disabled"
			},
			"fingerprint": "2AwuuoYRFd8BDmpz",
			"funding": "credit",
			"incremental_authorization": {
			  "status": "unavailable"
			},
			"installments": null,
			"last4": "4242",
			"mandate": null,
			"multicapture": {
			  "status": "unavailable"
			},
			"network": "visa",
			"network_token": {
			  "used": false
			},
			"overcapture": {
			  "maximum_amount_capturable": 8888,
			  "status": "unavailable"
			},
			"three_d_secure": null,
			"wallet": null
		  },
		  "type": "card"
		},
		"receipt_email": null,
		"receipt_number": null,
		"receipt_url": "https://pay.stripe.com/receipts/invoices/CAcaFwoVYWNjdF8xTm9XeURHUDVXT1NWOUNsKP--yKgGMga-TK6YCS46LBb8MB8e8dSmDI77jR7sg0I7ZE2SoobvF_gyZJTNcjAtAFdX3Cnhazyzzo8-?s=ap",
		"refunded": false,
		"review": null,
		"shipping": null,
		"source": null,
		"source_transfer": null,
		"statement_descriptor": null,
		"statement_descriptor_suffix": null,
		"status": "succeeded",
		"transfer_data": null,
		"transfer_group": null
	  }
	},
	"id": "evt_3NuOj3GP5WOSV9Cl0uMdbs2Y",
	"livemode": false,
	"pending_webhooks": 1,
	"request": {
	  "id": "req_0tViCOPOdxXtyC",
	  "idempotency_key": "1c54f60f-7cca-4bbc-9f85-2cc61aac72f2"
	},
	"type": "charge.succeeded"
  }`
