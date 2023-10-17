// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"strings"
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Droplist struct {
	writer IWriter
}

func NewDroplist(writer IWriter) ITransform {
	return &Droplist{writer: writer}
}

func (x *Droplist) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		values = strings.Split(values[0], " ")
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"droplist": mmdbtype.Map{
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

	job.Tools().Items["droplist.blacklist.isBlacklisted"] = Item{
		Item:        "ip/droplist/blacklist.isBlacklisted",
		Enabled:     true,
		GJSON:       "droplist.blacklist.isBlacklisted",
		Description: "IP is blacklisted.",
		Type:        common.Boolean.String()}
	job.Tools().Items["droplist.blacklist.lastReportedAt"] = Item{
		Item:        "ip/droplist/blacklist.lastReportedAt",
		Enabled:     true,
		GJSON:       "droplist.blacklist.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
