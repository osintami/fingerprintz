// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UdgerTor struct {
	writer IWriter
}

func NewUdgerTor(writer IWriter) ITransform {
	return &UdgerTor{writer: writer}
}

func (x *UdgerTor) Transform(job IETLJob) error {
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"udger": mmdbtype.Map{
					"tor": mmdbtype.Map{
						"isTorExitNode":  mmdbtype.Bool(true),
						"host":           mmdbtype.String(values[1]),
						"city":           mmdbtype.String(values[4]),
						"countryCode":    mmdbtype.String(values[3]),
						"lastReportedAt": mmdbtype.String(values[5]),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["udger.tor.isTorExitNode"] = Item{
		Item:        "ip/udger/tor.isTorExitNode",
		Enabled:     true,
		GJSON:       "udger.tor.isTorExitNode",
		Description: "IP belongs to a Tor Exit node.",
		Type:        common.Boolean.String()}
	job.Tools().Items["udger.tor.host"] = Item{
		Item:        "ip/udger/tor.host",
		Enabled:     true,
		GJSON:       "udger.tor.host",
		Description: "Hostname.",
		Type:        common.String.String()}
	job.Tools().Items["udger.tor.city"] = Item{
		Item:        "ip/udger/tor.city",
		Enabled:     true,
		GJSON:       "udger.tor.city",
		Description: "GEO city.",
		Type:        common.String.String()}
	job.Tools().Items["udger.tor.countryCode"] = Item{
		Item:        "ip/udger/tor.countryCode",
		Enabled:     true,
		GJSON:       "udger.tor.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["udger.tor.lastReportedAt"] = Item{
		Item:        "ip/udger/tor.lastReportedAt",
		Enabled:     true,
		GJSON:       "udger.tor.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
