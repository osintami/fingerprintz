// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UdgerSSH struct {
	writer IWriter
}

func NewUdgerSSH(writer IWriter) ITransform {
	return &UdgerSSH{writer: writer}
}

func (x *UdgerSSH) Transform(job IETLJob) error {
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"udger": mmdbtype.Map{
					"ssh": mmdbtype.Map{
						"isSshAttack":    mmdbtype.Bool(true),
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

	job.Tools().Items["udger.ssh.isSshAttack"] = Item{
		Item:        "ip/udger/ssh.isSshAttack",
		Enabled:     true,
		GJSON:       "udger.ssh.isSshAttack",
		Description: "IP is associated with an HTTP attacker.",
		Type:        common.Boolean.String()}
	job.Tools().Items["udger.ssh.host"] = Item{
		Item:        "ip/udger/ssh.host",
		Enabled:     true,
		GJSON:       "udger.ssh.host",
		Description: "Hostname.",
		Type:        common.String.String()}
	job.Tools().Items["udger.ssh.city"] = Item{
		Item:        "ip/udger/ssh.city",
		Enabled:     true,
		GJSON:       "udger.ssh.city",
		Description: "GEO city.",
		Type:        common.String.String()}
	job.Tools().Items["udger.ssh.countryCode"] = Item{
		Item:        "ip/udger/ssh.countryCode",
		Enabled:     true,
		GJSON:       "udger.ssh.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["udger.ssh.lastReportedAt"] = Item{
		Item:        "ip/udger/ssh.lastReportedAt",
		Enabled:     true,
		GJSON:       "udger.ssh.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
