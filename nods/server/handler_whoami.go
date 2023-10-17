// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Location struct {
	Latitude  float64
	Longitude float64
	City      string
	Country   string
	Continent string
}

func (x *NormalizedDataServer) GetWhoamiHandler(w http.ResponseWriter, r *http.Request) {
	x.PostWhoamiHandler(w, r)
}
func (x *NormalizedDataServer) PostWhoamiHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := x.params.Keys(r)
	if len(keys) == 0 || err != nil {
		common.SendError(w, ErrMissingInputs, http.StatusBadRequest)
		return
	}
	x.WhoamiHandler(w, r, keys)
}

func (x *NormalizedDataServer) WhoamiHandler(w http.ResponseWriter, r *http.Request, keys common.DataInputs) {
	out := x.WhoamiInfo(r.Context(), keys)
	w.Header().Set("Last-Modified", out.LastSeen)
	common.SendJSON(w, out)
}

func (x *NormalizedDataServer) WhoamiInfo(ctx context.Context, keys map[string]string) *common.WhoamiInfo {

	out := &common.WhoamiInfo{IpAddr: keys[CATEGORY_IPADDR], LastSeen: time.Now().Format(common.GO_DEFAULT_DATE)}

	uri := NewItemSplitter("rule/osintami/isTor")
	keys[common.INPUT_RULE] = uri.Key()
	value, _ := x.router.DataValue(ctx, uri, keys)
	if !value.Result.IsEmpty() {
		out.Tor = *value.Result.Bool
	}
	uri = NewItemSplitter("rule/osintami/isCloudNode")
	keys[common.INPUT_RULE] = uri.Key()
	value, _ = x.router.DataValue(ctx, uri, keys)
	if !value.Result.IsEmpty() {
		out.CloudNode = *value.Result.Bool
	}
	uri = NewItemSplitter("rule/osintami/isProxy")
	keys[common.INPUT_RULE] = uri.Key()
	value, _ = x.router.DataValue(ctx, uri, keys)
	if !value.Result.IsEmpty() {
		out.Proxy = *value.Result.Bool
	}
	uri = NewItemSplitter("rule/osintami/isVPN")
	keys[common.INPUT_RULE] = uri.Key()
	value, _ = x.router.DataValue(ctx, uri, keys)
	if !value.Result.IsEmpty() {
		out.VPN = *value.Result.Bool
	}
	uri = NewItemSplitter("rule/osintami/isBot")
	keys[common.INPUT_RULE] = uri.Key()
	value, _ = x.router.DataValue(ctx, uri, keys)
	if !value.Result.IsEmpty() {
		out.Bot = *value.Result.Bool
	}
	uri = NewItemSplitter("rule/osintami/isBlacklisted")
	keys[common.INPUT_RULE] = uri.Key()
	value, _ = x.router.DataValue(ctx, uri, keys)
	if !value.Result.IsEmpty() {
		out.Blacklist = *value.Result.Bool
	}

	uri = NewItemSplitter("ip/maxmind/location")
	keys[common.INPUT_RULE] = uri.Key()
	value, _ = x.router.DataValue(ctx, uri, keys)
	if !value.Result.IsEmpty() {
		location := Location{}
		err := json.Unmarshal([]byte(value.Result.Raw), &location)
		if err != nil {
			log.Error().Err(err).Str("component", "login").Msg("location json parse failure")
		} else {
			// TODO:  if we have GPS coordinates and no city, run against Google maps?
			out.Latitude = location.Latitude
			out.Longitude = location.Longitude
			out.City = location.City
			out.Country = location.Country
			out.Continent = location.Continent
		}
	}
	return out
}
