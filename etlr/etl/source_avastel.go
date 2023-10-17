// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Avastel struct{ writer IWriter }

func NewAvastel(writer IWriter) ITransform {
	return &Avastel{writer: writer}
}

func (x *Avastel) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		if values[0] == "ip_address" || values[0] == "" {
			return ErrBadSourceData
		}
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"avastel": mmdbtype.Map{
					"bot": mmdbtype.Map{
						"isBot":          mmdbtype.Bool(true),
						"ASN":            mmdbtype.String(values[1]),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["avastel.bot.isBot"] = Item{
		Item:        "ip/avastel/bot.isBot",
		Enabled:     true,
		GJSON:       "avastel.bot.isBot",
		Description: "IP has a bot associated with it.",
		Type:        common.Boolean.String()}
	job.Tools().Items["avastel.bot.ASN"] = Item{
		Item:        "ip/avastel/bot.ASN",
		Enabled:     true,
		GJSON:       "avastel.bot.ASN",
		Description: "ASN.",
		Type:        common.String.String()}
	job.Tools().Items["avastel.bot.lastReportedAt"] = Item{
		Item:        "ip/avastel/bot.lastReportedAt",
		Enabled:     true,
		GJSON:       "avastel.bot.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
