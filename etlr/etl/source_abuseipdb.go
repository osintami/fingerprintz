// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type AbuseIPDB struct {
	writer IWriter
}

func NewAbuseIPDB(writer IWriter) ITransform {
	return &AbuseIPDB{writer: writer}
}

func (x *AbuseIPDB) Transform(job IETLJob) error {
	type AbuseIPDBEntry struct {
		Meta struct {
			GeneratedAt time.Time `json:"generatedAt"`
		} `json:"meta"`
		Data []struct {
			IPAddress            string    `json:"ipAddress"`
			CountryCode          string    `json:"countryCode"`
			AbuseConfidenceScore int       `json:"abuseConfidenceScore"`
			LastReportedAt       time.Time `json:"lastReportedAt"`
		} `json:"data"`
	}
	entry := &AbuseIPDBEntry{}
	err := common.LoadJson(job.Info().inputFile, entry)
	for _, item := range entry.Data {
		cidr, err := job.Tools().Network.ParseCIDR(item.IPAddress)
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", item.IPAddress).Msg("parse cidr")
			continue
		} else {
			entry := mmdbtype.Map{
				"abuseipdb": mmdbtype.Map{
					"blacklist": mmdbtype.Map{
						"isBlacklisted":   mmdbtype.Bool(true),
						"lastReportedAt":  mmdbtype.String(item.LastReportedAt.String()),
						"confidenceScore": mmdbtype.Int32(item.AbuseConfidenceScore),
						"countryCode":     mmdbtype.String(item.CountryCode),
					},
				},
			}
			x.writer.Insert(cidr, entry)
		}
	}

	job.Tools().Items["abuseipdb.blacklist.isBlacklisted"] = Item{
		Item:        "ip/abuseipdb/blacklist.isBlacklisted",
		Enabled:     true,
		GJSON:       "abuseipdb.blacklist.isBlacklisted",
		Description: "IP is blacklisted.",
		Type:        common.Boolean.String()}
	job.Tools().Items["abuseipdb.blacklist.lastReportedAt"] = Item{
		Item:        "ip/abuseipdb/blacklist.lastReportedAt",
		Enabled:     true,
		GJSON:       "abuseipdb.blacklist.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}
	job.Tools().Items["abuseipdb.blacklist.confidenceScore"] = Item{
		Item:        "ip/abuseipdb/blacklist.confidenceScore",
		Enabled:     true,
		GJSON:       "abuseipdb.blacklist.confidenceScore",
		Description: "Blacklist confidence s",
		Type:        common.Integer.String()}
	job.Tools().Items["abuseipdb.blacklist.countryCode"] = Item{
		Item:        "ip/abuseipdb/blacklist.countryCode",
		Enabled:     true,
		GJSON:       "abuseipdb.blacklist.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}

	return err
}
