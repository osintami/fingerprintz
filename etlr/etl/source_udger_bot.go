// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"strings"

	"github.com/biter777/countries"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UdgerBot struct {
	writer IWriter
}

func NewUdgerBot(writer IWriter) ITransform {
	return &UdgerBot{writer: writer}
}

func (x *UdgerBot) Transform(job IETLJob) error {
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[1])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[1]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"udger": mmdbtype.Map{
					"bot": mmdbtype.Map{
						"isBot":          mmdbtype.Bool(true),
						"host":           mmdbtype.String(values[2]),
						"city":           mmdbtype.String(strings.TrimPrefix(values[4], " ")),
						"countryCode":    mmdbtype.String(countries.ByName(values[3]).Alpha2()),
						"lastReportedAt": mmdbtype.String(values[5]),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["udger.bot.isBot"] = Item{
		Item:        "ip/udger/bot.isBot",
		Enabled:     true,
		GJSON:       "udger.bot.isBot",
		Description: "IP is associated with a bot.",
		Type:        common.Boolean.String()}
	job.Tools().Items["udger.bot.host"] = Item{
		Item:        "ip/udger/bot.host",
		Enabled:     true,
		GJSON:       "udger.bot.host",
		Description: "Hostname.",
		Type:        common.String.String()}
	job.Tools().Items["udger.bot.city"] = Item{
		Item:        "ip/udger/bot.city",
		Enabled:     true,
		GJSON:       "udger.bot.city",
		Description: "GEO city.",
		Type:        common.String.String()}
	job.Tools().Items["udger.bot.countryCode"] = Item{
		Item:        "ip/udger/bot.countryCode",
		Enabled:     true,
		GJSON:       "udger.bot.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["udger.bot.lastReportedAt"] = Item{
		Item:        "ip/udger/bot.lastReportedAt",
		Enabled:     true,
		GJSON:       "udger.bot.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
