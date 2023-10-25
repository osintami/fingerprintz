// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/osintami/fingerprintz/common"
)

var ErrFileNotModified = errors.New("file not modified")

type UnwantedInfo struct {
	IpAddr     string
	UserAgent  string
	IsUnwanted bool
	Count      int
	LastSeen   time.Time
}

func (x *GatewayServer) UnwantedHandler(w http.ResponseWriter, r *http.Request) {
	x.logUnwantedGuests(r)
	common.SendError(w, ErrNeedToPay, http.StatusPaymentRequired)
}

/*
	func (x *GatewayServer) DownloadHandler(w http.ResponseWriter, r *http.Request) {
		user, err := x.findAccount(r)
		if err != nil {
			common.SendError(w, ErrNotAuthorized, http.StatusForbidden)
			return
		}
		if !user.IsAdmin() {
			log.Error().Err(err).Str("component", "gateway").Str("service", "download").Msg("admin required")
			common.SendError(w, ErrNotAuthorized, http.StatusForbidden)
			return
		}

		file, err := os.Stat("unwanted.db")
		if err != nil {
			common.SendError(w, ErrNotFound, http.StatusInternalServerError)

		}
		modTime := file.ModTime().Format(common.GO_DEFAULT_DATE)
		w.Header().Set("Content-Type", "application/octet-stream")
		if r.Header.Get("If-Modified-Since") == modTime {
			common.SendError(w, ErrFileNotModified, http.StatusNotModified)
			return
		}
		x.cache.Persist()
		w.Header().Set("Last-Modified", modTime)
		http.ServeFile(w, r, "unwanted.db")
	}
*/
func (x *GatewayServer) logUnwantedGuests(r *http.Request) {
	ip := common.IpAddr(r)
	if ip != "127.0.0.1" && ip != "localhost" {
		//log.Info().Str("component", "teapot").Str("addr", ip).Msg("")
		item, found := x.cache.Get(ip)
		if found {
			unwantedInfo := item.(*UnwantedInfo)
			unwantedInfo.Count += 1
			x.cache.Set(ip, unwantedInfo, -1)
		} else {
			x.cache.Set(ip, &UnwantedInfo{IpAddr: ip, UserAgent: r.UserAgent(), IsUnwanted: true, Count: 1, LastSeen: time.Now()}, -1)
		}
	}
}
