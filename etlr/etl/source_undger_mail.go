// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UdgerSMTP struct {
	writer IWriter
}

func NewUdgerSMTP(writer IWriter) ITransform {
	return &UdgerSMTP{writer: writer}
}

func (x *UdgerSMTP) Transform(job IETLJob) error {
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"udger": mmdbtype.Map{
					"mail": mmdbtype.Map{
						"isMailAttack":   mmdbtype.Bool(true),
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

	job.Tools().Items["udger.mail.isMailAttack"] = Item{
		Item:        "ip/udger/mail.isMailAttack",
		Enabled:     true,
		GJSON:       "udger.mail.isMailAttack",
		Description: "IP is associated with a mail attacker.",
		Type:        common.Boolean.String()}
	job.Tools().Items["udger.mail.host"] = Item{
		Item:        "ip/udger/mail.host",
		Enabled:     true,
		GJSON:       "udger.mail.host",
		Description: "Hostname.",
		Type:        common.String.String()}
	job.Tools().Items["udger.mail.city"] = Item{
		Item:        "ip/udger/mail.city",
		Enabled:     true,
		GJSON:       "udger.mail.city",
		Description: "GEO city.",
		Type:        common.String.String()}
	job.Tools().Items["udger.mail.countryCode"] = Item{
		Item:        "ip/udger/mail.countryCode",
		Enabled:     true,
		GJSON:       "udger.mail.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["udger.mail.lastReportedAt"] = Item{
		Item:        "ip/udger/mail.lastReportedAt",
		Enabled:     true,
		GJSON:       "udger.mail.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
