// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
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
	fingerprint := BiometricsScan(r.Context(), x.nods, keys)
	jwt, err := x.signer.SignJWT(fingerprint)
	if err != nil {
		log.Error().Err(err).Str("component", "fingerprint").Str("ehash", keys["eh"]).Msg("create fingerprint JWT failed")
		common.SendError(w, ErrFingerprintSmudged, http.StatusInternalServerError)
	}
	w.Write([]byte(jwt))
}
