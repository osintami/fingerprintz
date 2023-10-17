// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"net"
	"strings"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type OnionOOEntry struct {
	Relays []struct {
		Fingerprint string   `json:"fingerprint"`
		OrAddresses []string `json:"or_addresses"`
		LastSeen    string   `json:"last_seen"`
		FirstSeen   string   `json:"first_seen"`
		Running     bool     `json:"running"`
		Flags       []string `json:"flags"`
		Country     string   `json:"country"`
		As          string   `json:"as"`
		AsName      string   `json:"as_name"`
		Contact     string   `json:"contact"`
	} `json:"relays"`
}

type OnionOO struct {
	writer IWriter
}

func NewOnionOO(writer IWriter) ITransform {
	return &OnionOO{writer: writer}
}

func (x *OnionOO) Transform(job IETLJob) error {
	entry := &OnionOOEntry{}
	err := common.LoadJson(job.Info().inputFile, entry)
	for _, item := range entry.Relays {
		for _, host := range item.OrAddresses {
			var cidr *net.IPNet
			var err error
			if strings.ContainsAny(host, ".") {
				// IPv4
				ipv4 := strings.Split(host, ":")[0]
				cidr, err = job.Tools().Network.ParseCIDR(ipv4)
				if err != nil {
					log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", ipv4).Msg("parse cidr")
					continue
				}
			} else {
				// IPv6
				ipv6 := strings.Split(host, "]:")[0]
				ipv6 = strings.TrimLeft(ipv6, "[")
				cidr, err = job.Tools().Network.ParseCIDR(ipv6)
				if err != nil {
					log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", ipv6).Msg("parse cidr")
					continue
				}
			}

			// TODO:  we need history to set a wasTorNode flag, very useful
			relayMap := mmdbtype.Map{
				"firstReportedAt": mmdbtype.String(item.FirstSeen),
				"lastReportedAt":  mmdbtype.String(item.LastSeen),
				"countryCode":     mmdbtype.String(strings.ToUpper(item.Country)),
				"ASN":             mmdbtype.String(item.As),
				"ASO":             mmdbtype.String(item.AsName),
				"Contact":         mmdbtype.String(item.Contact),
				"isTorNode":       mmdbtype.Bool(item.Running),
			}

			for _, flag := range item.Flags {
				relayMap[mmdbtype.String("is"+strings.Title(flag))] = mmdbtype.Bool(true)
			}

			entry := mmdbtype.Map{
				"onionoo": mmdbtype.Map{
					"tor": relayMap,
				},
			}
			x.writer.Insert(cidr, entry)
		}
	}

	job.Tools().Items["onionoo.tor.isTorNode"] = Item{
		Item:        "ip/onionoo/tor.isTorNode",
		Enabled:     true,
		GJSON:       "onionoo.tor.isTorNode",
		Description: "IP is an active Tor relay node.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.firstReportedAt"] = Item{
		Item:        "ip/onionoo/tor.firstReportedAt",
		Enabled:     true,
		GJSON:       "onionoo.tor.firstReportedAt",
		Description: "First update date/time.",
		Type:        common.Date.String()}
	job.Tools().Items["onionoo.tor.lastReportedAt"] = Item{
		Item:        "ip/onionoo/tor.lastReportedAt",
		Enabled:     true,
		GJSON:       "onionoo.tor.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}
	job.Tools().Items["onionoo.tor.countryCode"] = Item{
		Item:        "ip/onionoo/tor.countryCode",
		Enabled:     true,
		GJSON:       "onionoo.tor.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["onionoo.tor.ASN"] = Item{
		Item:        "ip/onionoo/tor.ASN",
		Enabled:     true,
		GJSON:       "onionoo.tor.ASN",
		Description: "ASN.",
		Type:        common.String.String()}
	job.Tools().Items["onionoo.tor.ASO"] = Item{
		Item:        "ip/onionoo/tor.ASO",
		Enabled:     true,
		GJSON:       "onionoo.tor.ASO",
		Description: "ASO.",
		Type:        common.String.String()}
	job.Tools().Items["onionoo.tor.isExit"] = Item{
		Item:        "ip/onionoo/tor.isExit",
		Enabled:     true,
		GJSON:       "onionoo.tor.isExit",
		Description: "IP has Tor flag Exit.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.isValid"] = Item{
		Item:        "ip/onionoo/tor.isValid",
		Enabled:     true,
		GJSON:       "onionoo.tor.isValid",
		Description: "IP has Tor flag Valid.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.isV2Dir"] = Item{
		Item:        "ip/onionoo/tor.isV2Dir",
		Enabled:     true,
		GJSON:       "onionoo.tor.isV2Dir",
		Description: "IP has Tor flag V2Dir.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.isStaleDesc"] = Item{
		Item:        "ip/onionoo/tor.isStaleDesc",
		Enabled:     true,
		GJSON:       "onionoo.tor.isStaleDesc",
		Description: "IP has Tor flag StaleDesc.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.isStable"] = Item{
		Item:        "ip/onionoo/tor.isStable",
		Enabled:     true,
		GJSON:       "onionoo.tor.isStable",
		Description: "IP has Tor flag Stable.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.isHsDir"] = Item{
		Item:        "ip/onionoo/tor.isHsDir",
		Enabled:     true,
		GJSON:       "onionoo.tor.isHsDir",
		Description: "IP has Tor flag HSDir.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.isGuard"] = Item{
		Item:        "ip/onionoo/tor.isGuard",
		Enabled:     true,
		GJSON:       "onionoo.tor.isGuard",
		Description: "IP is a Tor Guard node.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.isFast"] = Item{
		Item:        "ip/onionoo/tor.isFast",
		Enabled:     true,
		GJSON:       "onionoo.tor.isFast",
		Description: "IP is a Tor Fast node.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.isAuthority"] = Item{
		Item:        "ip/onionoo/tor.isAuthority",
		Enabled:     true,
		GJSON:       "onionoo.tor.isAuthority",
		Description: "IP is a Tor Authority node.",
		Type:        common.Boolean.String()}
	job.Tools().Items["onionoo.tor.isBadExit"] = Item{
		Item:        "ip/onionoo/tor.isBadExit",
		Enabled:     true,
		GJSON:       "onionoo.tor.isBadExit",
		Description: "IP has Tor flag BadExit.",
		Type:        common.Boolean.String()}

	return err
}
