// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UHB struct{ writer IWriter }

func NewUHB(writer IWriter) ITransform {
	return &UHB{writer: writer}
}

func (x *UHB) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				mmdbtype.String(job.Source().Name): mmdbtype.Map{
					"blacklist": mmdbtype.Map{
						"isBlacklisted":  mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["uhb.blacklist.isBlacklisted"] = Item{
		Item:        "ip/uhb/blacklist.isBlacklisted",
		Enabled:     true,
		GJSON:       "uhb.blacklist.isBlacklisted",
		Description: "IP is blacklisted.",
		Type:        common.Boolean.String()}
	job.Tools().Items["uhb.blacklist.lastReportedAt"] = Item{
		Item:        "ip/uhb/blacklist.lastReportedAt",
		Enabled:     true,
		GJSON:       "uhb.blacklist.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
