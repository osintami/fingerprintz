// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Fingerprint struct {
	LastSeenAt string
	EHash      string
	Latitude   float64
	Longitude  float64
	City       string
	Country    string
	IpAddr     string
	UserAgent  string
	DeviceId   string
	NetworkId  string
	PartnerId  string
	Version    string
}

type FingerprintClaims struct {
	LastSeenAt string  `json:"time"`
	EHash      string  `json:"ehash"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	IpAddr     string  `json:"ip"`
	UserAgent  string  `json:"ua"`
	DeviceId   string  `json:"hwid"`
	NetworkId  string  `json:"nid"`
	PartnerId  string  `json:"pid"`
	Version    string  `json:"version"`
	jwt.RegisteredClaims
}

func BiometricsScan(ctx context.Context, nods common.INods, keys map[string]string) *Fingerprint {

	// API consumer can send an email hashed with SHA256 16 bit lowercase OR a raw email address
	if keys["email"] != "" && keys["eh"] == "" {
		keys["eh"] = common.SHA256ToBase16Lowercase(keys["email"])
	}

	nodsInfo, err := nods.Whoami(keys["ip"])
	if err != nil {
		// NOTE:  this just can't happen
		log.Error().Err(err).Str("component", "fingerprint").Str("ehash", keys["eh"]).Msg("collect network info failed")
		return nil
	}

	var inet uint = 0x0000
	if nodsInfo.Blacklist {
		inet |= BIT_BLACKLIST
	}
	if nodsInfo.CloudNode {
		inet |= BIT_CLOUD_NODE
	}
	if nodsInfo.Bot {
		inet |= BIT_KNOWN_BOT
	}
	if nodsInfo.Tor {
		inet |= BIT_TOR_IN_USE
	}
	if nodsInfo.Proxy {
		inet |= BIT_PROXY_IN_USE
	}
	if nodsInfo.VPN {
		inet |= BIT_VPN_IN_USE
	}

	return &Fingerprint{
		LastSeenAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		EHash:      keys["eh"],
		Latitude:   nodsInfo.Latitude,
		Longitude:  nodsInfo.Longitude,
		City:       nodsInfo.City,
		Country:    nodsInfo.Country,
		IpAddr:     keys["ip"],
		UserAgent:  keys["ua"],
		DeviceId:   keys["hw"],
		NetworkId:  fmt.Sprintf("0X%04X", inet),
		PartnerId:  keys["pid"],
	}
}
