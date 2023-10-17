// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Google struct {
	writer IWriter
}

func NewGoogle(writer IWriter) ITransform {
	return &Google{writer: writer}
}

func (x *Google) Transform(job IETLJob) error {

	type GoogleService struct {
		Ipv4Prefix string
		Ipv6Prefix string
		Service    string // always "Google Cloud"
		Region     string `json:"scope"`
	}

	type GoogleEntry struct {
		SyncToken      string
		LastReportedAt string `json:"creationTime"`
		Prefixes       []GoogleService
	}

	entry := &GoogleEntry{}
	err := common.LoadJson(job.Info().inputFile, entry)
	for _, item := range entry.Prefixes {
		var ip = item.Ipv4Prefix
		if ip == "" {
			ip = item.Ipv6Prefix
		}
		cidr, err := job.Tools().Network.ParseCIDR(ip)
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", ip).Msg("parse cidr")
			continue
		} else {
			entry := mmdbtype.Map{
				"google": mmdbtype.Map{
					"cloud": mmdbtype.Map{
						"isGoogle":       mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(entry.LastReportedAt),
						"region":         mmdbtype.String(item.Region),
						"product":        mmdbtype.String(item.Service),
					},
				},
			}

			x.writer.Insert(cidr, entry)
		}
	}

	job.Tools().Items["google.cloud.isGoogle"] = Item{
		Item:        "ip/google/cloud.isGoogle",
		Enabled:     true,
		GJSON:       "google.cloud.isGoogle",
		Description: "IP belongs to a cloud provider.",
		Type:        common.Boolean.String()}
	job.Tools().Items["google.cloud.lastReportedAt"] = Item{
		Item:        "ip/google/cloud.lastReportedAt",
		Enabled:     true,
		GJSON:       "google.cloud.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}
	job.Tools().Items["google.cloud.region"] = Item{
		Item:        "ip/google/cloud.region",
		Enabled:     true,
		GJSON:       "google.cloud.region",
		Description: "GCP region.",
		Type:        common.Date.String()}
	job.Tools().Items["google.cloud.product"] = Item{
		Item:        "ip/google/cloud.product",
		Enabled:     true,
		GJSON:       "google.cloud.product",
		Description: "GCP product.",
		Type:        common.Date.String()}

	return err

}
