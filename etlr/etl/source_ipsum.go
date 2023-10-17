// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"strconv"
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type IpSUM struct {
	writer IWriter
}

func NewIpSUM(writer IWriter) ITransform {
	return &IpSUM{writer: writer}
}

func (x *IpSUM) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		count, err := strconv.Atoi(string(values[1]))
		if err != nil {
			// junk data
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", string(values[1])).Msg("junk data")
			return ErrBadSourceData
		}
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"ipsum": mmdbtype.Map{
					"blacklist": mmdbtype.Map{
						"isBlacklisted":  mmdbtype.Bool(true),
						"blacklistCount": mmdbtype.Int32(count),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["ipsum.blacklist.isBlacklisted"] = Item{
		Item:        "ip/ipsum/blacklist.isBlacklisted",
		Enabled:     true,
		GJSON:       "ipsum.blacklist.isBlacklisted",
		Description: "IP is blacklisted.",
		Type:        common.Boolean.String()}
	job.Tools().Items["ipsum.blacklist.blacklistCount"] = Item{
		Item:        "ip/ipsum/blacklist.blacklistCount",
		Enabled:     true,
		GJSON:       "ipsum.blacklist.isBlacklisted",
		Description: "IP is on x blacklists.",
		Type:        common.Integer.String()}
	job.Tools().Items["ipsum.blacklist.lastReportedAt"] = Item{
		Item:        "ip/ipsum/blacklist.lastReportedAt",
		Enabled:     true,
		GJSON:       "ipsum.blacklist.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
