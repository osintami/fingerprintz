// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"strings"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UdgerProxy struct {
	writer IWriter
}

func NewUdgerProxy(writer IWriter) ITransform {
	return &UdgerProxy{writer: writer}
}

func (x *UdgerProxy) Transform(job IETLJob) error {
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"udger": mmdbtype.Map{
					"proxy": mmdbtype.Map{
						"isProxy":        mmdbtype.Bool(true),
						"host":           mmdbtype.String(values[2]),
						"city":           mmdbtype.String(strings.TrimPrefix(values[4], " ")),
						"countryCode":    mmdbtype.String(values[3]),
						"lastReportedAt": mmdbtype.String(values[5]),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["udger.proxy.isProxy"] = Item{
		Item:        "ip/udger.proxy/proxy.isProxy",
		Enabled:     true,
		GJSON:       "udger.proxy.isProxy",
		Description: "IP is associated with a proxy service.",
		Type:        common.Boolean.String()}
	job.Tools().Items["udger.proxy.host"] = Item{
		Item:        "ip/udger.proxy/proxy.host",
		Enabled:     true,
		GJSON:       "udger.proxy.host",
		Description: "Hostname.",
		Type:        common.String.String()}
	job.Tools().Items["udger.proxy.city"] = Item{
		Item:        "ip/udger.proxy/proxy.city",
		Enabled:     true,
		GJSON:       "udger.proxy.city",
		Description: "GEO city.",
		Type:        common.String.String()}
	job.Tools().Items["udger.proxy.countryCode"] = Item{
		Item:        "ip/udger.proxy/proxy.countryCode",
		Enabled:     true,
		GJSON:       "udger.proxy.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["udger.proxy.lastReportedAt"] = Item{
		Item:        "ip/udger.proxy/proxy.lastReportedAt",
		Enabled:     true,
		GJSON:       "udger.proxy.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
