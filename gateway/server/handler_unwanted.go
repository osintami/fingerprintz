// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
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

func (x *GatewayServer) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	user, err := x.findAccount(r)
	if err != nil || !user.IsAdmin() {
		common.SendError(w, ErrNotAuthorized, http.StatusForbidden)
		return
	}

	name := common.QueryParam(r, "name")
	if name != "unwanted" {
		common.SendError(w, ErrNotAuthorized, http.StatusForbidden)
		return
	}

	file, err := os.Stat("unwanted.db")
	if err != nil {
		common.SendError(w, ErrNotFound, http.StatusInternalServerError)
		return
	}

	modTime := file.ModTime().Format(common.GO_DEFAULT_DATE)
	if r.Header.Get("If-Modified-Since") == modTime {
		common.SendError(w, ErrFileNotModified, http.StatusNotModified)
		return
	}
	x.cache.Persist()

	tm := time.Now()
	fileName := fmt.Sprintf("%s_%d_%02d_%02d.%s", "unwanted", tm.Year(), tm.Month(), tm.Day(), "json")
	err = x.cache.ToJSON(fileName)
	if err == nil {
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Last-Modified", modTime)
		http.ServeFile(w, r, fileName)
		x.cache.Clear()
		return
	}
	common.SendError(w, err, http.StatusInternalServerError)
}

func (x *GatewayServer) logUnwantedGuests(r *http.Request) {
	ip := common.IpAddr(r)
	if ip != "127.0.0.1" && ip != "localhost" {
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
