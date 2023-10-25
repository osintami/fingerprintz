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
	if fingerprint == nil {
		// NOTE:  nods is down, this should never happen
		common.SendError(w, ErrFingerprintSmudged, http.StatusInternalServerError)
		return
	}
	jwt, err := x.signer.SignJWT(fingerprint)
	if err != nil {
		log.Error().Err(err).Str("component", "fingerprint").Str("ehash", keys["eh"]).Msg("create fingerprint JWT failed")
		common.SendError(w, ErrFingerprintSmudged, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(jwt))
}
