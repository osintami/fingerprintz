// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type X4BNetTrollBot struct {
	writer IWriter
}

func NewX4BNetTroll(writer IWriter) ITransform {
	return &X4BNetTrollBot{writer: writer}
}

func (x *X4BNetTrollBot) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"x4bnet": mmdbtype.Map{
					"bot": mmdbtype.Map{
						"isForumBot":     mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["x4bnet.troll.isForumBot"] = Item{
		Item:        "ip/x4bnet.troll/bot.isForumBot",
		Enabled:     true,
		GJSON:       "x4bnet.bot.isForumBot",
		Description: "IP has a bot associated with it.",
		Type:        common.Boolean.String()}
	job.Tools().Items["x4bnet.troll.lastReportedAt"] = Item{
		Item:        "ip/x4bnet.troll/bot.lastReportedAt",
		Enabled:     true,
		GJSON:       "x4bnet.bot.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
