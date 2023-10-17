// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"time"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

func (x *WhoamiServer) GetRiskHandler(w http.ResponseWriter, r *http.Request) {
	x.PostRiskHandler(w, r)
}

func (x *WhoamiServer) PostRiskHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := x.params.Keys(r)
	if len(keys) == 0 || err != nil {
		common.SendError(w, ErrMissingInputs, http.StatusBadRequest)
		return
	}
	x.RisksHandler(w, r, keys)
}

func (x *WhoamiServer) RisksHandler(w http.ResponseWriter, r *http.Request, keys map[string]string) {

	signature := keys["fingerprint"]
	if signature == "" {
		common.SendError(w, ErrMissingInputs, http.StatusBadRequest)
		return
	}

	lastFingerprint, err := x.signer.DecodeJWT(signature)
	if err != nil {
		log.Error().Err(err).Str("component", "fingerprint").Str("partner", keys["pid"]).Str("fingerprint", signature).Msg("invalid JWT token")
		common.SendError(w, ErrFingerprintSmudged, http.StatusInternalServerError)
		return
	}

	currFingerprint := BiometricsScan(r.Context(), x.nods, keys)
	if currFingerprint == nil {
		common.SendError(w, common.ErrItemNotFound, http.StatusInternalServerError)
		return
	}

	info := &FingerprintInfo{}

	if currFingerprint.DeviceId == lastFingerprint.DeviceId {
		info.Match.Device = true
	}
	if currFingerprint.EHash == lastFingerprint.EHash {
		info.Match.Email = true
	}
	if currFingerprint.IpAddr == lastFingerprint.IpAddr {
		info.Match.IpAddr = true
	}
	if currFingerprint.Latitude == lastFingerprint.Latitude && currFingerprint.Longitude == lastFingerprint.Longitude {
		info.Match.Location = true
	} else {
		lastSeenAt, _ := time.Parse("2006-01-02T15:04:05Z07:00", lastFingerprint.LastSeenAt)
		hours := common.TimeInHours(lastSeenAt)
		miles := int(common.DistanceInMiles(currFingerprint.Latitude, currFingerprint.Longitude, lastFingerprint.Latitude, lastFingerprint.Longitude, 'M'))
		if hours < 24 && miles > 1000 {
			info.ImpossibleTravel = true
		}
	}
	if currFingerprint.UserAgent == lastFingerprint.UserAgent {
		info.Match.UserAgent = true
	}
	if currFingerprint.NetworkId == lastFingerprint.NetworkId {
		info.Match.Network = true
	}

	jwt, err := x.signer.SignJWT(currFingerprint)
	if err != nil {
		log.Error().Err(err).Str("component", "fingerprint").Str("partner", keys["pid"]).Str("fingerprint", signature).Msg("create fingerprint JWT failed")
		common.SendError(w, ErrFingerprintSmudged, http.StatusInternalServerError)
		return
	}

	info.Fingerprint = jwt
	common.SendPrettyJSON(w, info)
}
