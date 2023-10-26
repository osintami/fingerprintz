// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UnwantedEntry struct {
	IpAddr     string
	UserAgent  string
	IsUnwanted bool
	Count      int
	LastSeen   string
}

type Unwanted struct {
	writer IWriter
}

func NewUnwanted(writer IWriter) ITransform {
	return &Unwanted{writer: writer}
}

func (x *Unwanted) Transform(job IETLJob) error {
	list := []UnwantedEntry{}
	err := common.LoadJson(job.Info().inputFile, &list)
	for _, item := range list {
		cidr, err := job.Tools().Network.ParseCIDR(item.IpAddr)
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", item.IpAddr).Msg("parse cidr")
			continue
		}

		botMap := mmdbtype.Map{
			//"firstReportedAt": mmdbtype.String(item.FirstSeen),
			"userAgent":      mmdbtype.String(item.UserAgent),
			"isBot":          mmdbtype.Bool(true),
			"count":          mmdbtype.Int32(item.Count),
			"lastReportedAt": mmdbtype.String(item.LastSeen),
		}

		entry := mmdbtype.Map{
			"unwanted": mmdbtype.Map{
				"bot": botMap,
			},
		}
		x.writer.Insert(cidr, entry)
	}

	job.Tools().Items["unwanted.bot.isBot"] = Item{
		Item:        "ip/unwanted/isBot",
		Enabled:     true,
		GJSON:       "unwanted.bot.isBot",
		Description: "IP is an active bot scanning NODS.",
		Type:        common.Boolean.String()}

	job.Tools().Items["unwanted.bot.UserAgent"] = Item{
		Item:        "ip/unwanted/UserAgent",
		Enabled:     true,
		GJSON:       "unwanted.bot.UserAgent",
		Description: "The bots user-agent string.",
		Type:        common.Date.String()}

	job.Tools().Items["unwanted.bot.Count"] = Item{
		Item:        "ip/unwanted/count",
		Enabled:     true,
		GJSON:       "unwanted.bot.count",
		Description: "Number of times caught scanning.",
		Type:        common.Integer.String()}

	job.Tools().Items["unwanted.bot.lastReportedAt"] = Item{
		Item:        "ip/unwanted/lastReportedAt",
		Enabled:     true,
		GJSON:       "unwanted.bot.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}
	return err
}
