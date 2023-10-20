// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UdgerVPN struct {
	writer IWriter
}

func NewUdgerVPN(writer IWriter) ITransform {
	return &UdgerVPN{writer: writer}
}

func (x *UdgerVPN) Transform(job IETLJob) error {
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"udger": mmdbtype.Map{
					"vpn": mmdbtype.Map{
						"isVPN":          mmdbtype.Bool(true),
						"host":           mmdbtype.String(values[1]),
						"countryCode":    mmdbtype.String(values[2]),
						"lastReportedAt": mmdbtype.String(values[5]),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["udger.vpn.isVPN"] = Item{
		Item:        "ip/udger.vpn/vpn.isVPN",
		Enabled:     true,
		GJSON:       "udger.vpn.isVPN",
		Description: "IP is associated with a VPN service.",
		Type:        common.Boolean.String()}
	job.Tools().Items["udger.vpn.host"] = Item{
		Item:        "ip/udger.vpn/vpn.host",
		Enabled:     true,
		GJSON:       "udger.vpn.host",
		Description: "Hostname.",
		Type:        common.String.String()}
	job.Tools().Items["udger.vpn.countryCode"] = Item{
		Item:        "ip/udger.vpn/vpn.countryCode",
		Enabled:     true,
		GJSON:       "udger.vpn.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["udger.vpn.lastReportedAt"] = Item{
		Item:        "ip/udger.vpn/vpn.lastReportedAt",
		Enabled:     true,
		GJSON:       "udger.vpn.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
