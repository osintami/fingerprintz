// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type DigitalOcean struct {
	writer IWriter
}

func NewDigitalOcean(writer IWriter) ITransform {
	return &DigitalOcean{writer: writer}
}

func (x *DigitalOcean) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"digitalocean": mmdbtype.Map{
					"cloud": mmdbtype.Map{
						"isDigitalOcean": mmdbtype.Bool(true),
						"countryCode":    mmdbtype.String(values[1]),
						"region":         mmdbtype.String(values[2]),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}
			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["digitalocean.cloud.isDigitalOcean"] = Item{
		Item:        "ip/digitalocean/cloud.isDigitalOcean",
		Enabled:     true,
		GJSON:       "digitalocean.cloud.isDigitalOcean",
		Description: "IP belongs to a cloud provider.",
		Type:        common.Boolean.String()}
	job.Tools().Items["digitalocean.cloud.countryCode"] = Item{
		Item:        "ip/digitalocean/cloud.countryCode",
		Enabled:     true,
		GJSON:       "digitalocean.cloud.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["digitalocean.cloud.region"] = Item{
		Item:        "ip/digitalocean/cloud.region",
		Enabled:     true,
		GJSON:       "digitalocean.cloud.region",
		Description: "GEO region.",
		Type:        common.String.String()}
	job.Tools().Items["digitalocean.cloud.lastReportedAt"] = Item{
		Item:        "ip/digitalocean/cloud.lastReportedAt",
		Enabled:     true,
		GJSON:       "digitalocean.cloud.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
