// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UdgerHTTP struct {
	writer IWriter
}

func NewUdgerHTTP(writer IWriter) ITransform {
	return &UdgerHTTP{writer: writer}
}

func (x *UdgerHTTP) Transform(job IETLJob) error {
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"udger": mmdbtype.Map{
					"http": mmdbtype.Map{
						"isHttpAttack":   mmdbtype.Bool(true),
						"host":           mmdbtype.String(values[1]),
						"countryCode":    mmdbtype.String(values[3]),
						"city":           mmdbtype.String(values[4]),
						"lastReportedAt": mmdbtype.String(values[5]),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["udger.http.isHttpAttack"] = Item{
		Item:        "ip/udger.http/http.isHttpAttack",
		Enabled:     true,
		GJSON:       "udger.http.isHttpAttack",
		Description: "IP is associated with an HTTP attacker.",
		Type:        common.Boolean.String()}
	job.Tools().Items["udger.http.host"] = Item{
		Item:        "ip/udger.http/http.host",
		Enabled:     true,
		GJSON:       "udger.http.host",
		Description: "Hostname.",
		Type:        common.String.String()}
	job.Tools().Items["udger.http.city"] = Item{
		Item:        "ip/udger.http/http.city",
		Enabled:     true,
		GJSON:       "udger.http.city",
		Description: "GEO city.",
		Type:        common.String.String()}
	job.Tools().Items["udger.http.countryCode"] = Item{
		Item:        "ip/udger.http/http.countryCode",
		Enabled:     true,
		GJSON:       "udger.http.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["udger.http.lastReportedAt"] = Item{
		Item:        "ip/udger.http/http.lastReportedAt",
		Enabled:     true,
		GJSON:       "udger.http.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
