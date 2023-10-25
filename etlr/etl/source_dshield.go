// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"strconv"
	"strings"
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type DShield struct {
	writer IWriter
}

func NewDShield(writer IWriter) ITransform {
	return &DShield{writer: writer}
}

func (x *DShield) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		if strings.HasPrefix(values[0], "Start") {
			return ErrBadSourceData
		}
		ip := values[0] + "/" + values[2]

		cidr, err := job.Tools().Network.ParseCIDR(ip)
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			attacks, _ := strconv.ParseInt(values[3], 10, 32)
			entry := mmdbtype.Map{
				"dshield": mmdbtype.Map{
					"blacklist": mmdbtype.Map{
						"isBlacklisted":  mmdbtype.Bool(true),
						"attackCount":    mmdbtype.Int32(attacks),
						"ASN":            mmdbtype.String(values[4]),
						"countryCode":    mmdbtype.String(values[5]),
						"email":          mmdbtype.String(values[6]),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["dshield.blacklist.isBlacklisted"] = Item{
		Item:        "ip/dshield/blacklist.isBlacklisted",
		Enabled:     true,
		GJSON:       "dshield.blacklist.isBlacklisted",
		Description: "IP is blacklisted.",
		Type:        common.Boolean.String()}

	job.Tools().Items["dshield.blacklist.attackCount"] = Item{
		Item:        "ip/dshield/blacklist.attackCount",
		Enabled:     true,
		GJSON:       "dshield.blacklist.attackCount",
		Description: "IP has been reported in x attacks.",
		Type:        common.Integer.String()}
	job.Tools().Items["dshield.blacklist.ASN"] = Item{
		Item:        "ip/dshield/blacklist.ASN",
		Enabled:     true,
		GJSON:       "dshield.blacklist.ASN",
		Description: "ASN.",
		Type:        common.String.String()}
	job.Tools().Items["dshield.blacklist.countryCode"] = Item{
		Item:        "ip/dshield/blacklist.countryCode",
		Enabled:     true,
		GJSON:       "dshield.blacklist.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["dshield.blacklist.email"] = Item{
		Item:        "ip/dshield/blacklist.email",
		Enabled:     false,
		GJSON:       "dshield.blacklist.email",
		Description: "Contact email.",
		Type:        common.String.String()}
	job.Tools().Items["dshield.blacklist.lastReportedAt"] = Item{
		Item:        "ip/dshield/blacklist.lastReportedAt",
		Enabled:     true,
		GJSON:       "dshield.blacklist.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
