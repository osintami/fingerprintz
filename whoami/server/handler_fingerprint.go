// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

const (
	COOKIE_ID        = "osintami"
	COOKIE_DOMAIN    = "osintami.com"
	ONE_YEAR_SECONDS = 60 * 60 * 24 * 365
)

type FingerprintInfo struct {
	Match struct {
		Email     bool
		Location  bool
		IpAddr    bool
		UserAgent bool
		Device    bool
		Network   bool
	}
	ImpossibleTravel bool
	Fingerprint      string
}

const (
	BIT_ANONYMOUS    = 0x0001
	BIT_BLACKLIST    = 0x0002
	BIT_SUSPECT_BOT  = 0x0004
	BIT_KNOWN_BOT    = 0x0008
	BIT_CLOUD_NODE   = 0x0010
	BIT_PROXY_IN_USE = 0x0100
	BIT_TOR_IN_USE   = 0x0200
	BIT_VPN_IN_USE   = 0x0400
	// NOTE:  add critical data points as needed
)

func (x *WhoamiServer) GetFingerprintHandler(w http.ResponseWriter, r *http.Request) {
	x.PostFingerprintHandler(w, r)
}

func (x *WhoamiServer) PostFingerprintHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := x.params.Keys(r)
	if len(keys) == 0 || err != nil {
		common.SendError(w, ErrMissingInputs, http.StatusBadRequest)
		return
	}
	x.FingerprintHandler(w, r, keys)
}

func (x *WhoamiServer) FingerprintHandler(w http.ResponseWriter, r *http.Request, keys map[string]string) {

	keys["ip"] = common.IpAddr(r)
	keys["ua"] = r.UserAgent()

	fingerprint := BiometricsScan(r.Context(), x.nods, keys)

	if fingerprint == nil {
		log.Error().Str("component", "fingerprint").Str("ehash", keys["eh"]).Msg("create fingerprint failed")
		common.SendError(w, ErrFingerprintSmudged, http.StatusInternalServerError)
		return
	}

	fmt.Println("****** LATENCY START:", fingerprint.Latency)

	jwt, err := x.signer.SignJWT(fingerprint)
	if err != nil {
		log.Error().Err(err).Str("component", "fingerprint").Str("ehash", keys["eh"]).Msg("create fingerprint JWT failed")
		common.SendError(w, ErrFingerprintSmudged, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "")
	http.Redirect(w, r, x.latencyRedirect+"?fingerprint="+jwt, http.StatusFound)
}

func (x *WhoamiServer) GetLatencyHandler(w http.ResponseWriter, r *http.Request) {
	x.PostLatencyHandler(w, r)
}

func (x *WhoamiServer) PostLatencyHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := x.params.Keys(r)
	if len(keys) == 0 || err != nil {
		common.SendError(w, ErrMissingInputs, http.StatusBadRequest)
		return
	}
	x.LatencyHandler(w, r, keys)
}

func (x *WhoamiServer) LatencyHandler(w http.ResponseWriter, r *http.Request, keys map[string]string) {

	token := common.QueryParam(r, "fingerprint")
	if token == "" {
		common.SendError(w, errors.New("missing fingerprint"), http.StatusBadRequest)
		return
	}

	fingerprint, err := x.signer.DecodeJWT(keys["fingerprint"])
	if err != nil {
		common.SendError(w, err, http.StatusBadRequest)
		return
	}

	latency := time.Now().UnixMilli() - fingerprint.Latency

	fmt.Println("****** LATENCY:", latency)
	fingerprint.Latency = latency

	if fingerprint.IpAddr != common.IpAddr(r) || fingerprint.UserAgent != r.UserAgent() {
		common.SendError(w, errors.New("man in the middle"), http.StatusBadRequest)
		return
	}

	jwt, err := x.signer.SignJWT(fingerprint)
	if err != nil {
		common.SendError(w, err, http.StatusBadRequest)
		return
	}

	w.Write([]byte(jwt))
}
