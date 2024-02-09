// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"net"
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type LightswitchAggressive struct {
	writer IWriter
}

func NewLightswitchAggressive(writer IWriter) ITransform {
	return &LightswitchAggressive{writer: writer}
}

func (x *LightswitchAggressive) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		// NOTE:  this will get resource intensive if we aren't caching DNS lookups for some period of time
		addrs, err := net.LookupIP(values[1])
		if err != nil {
			// skipping this entry, don't return an error or we bail on the whole file
			log.Error().Err(err).Str("component", job.Source().Name).Str("domain", values[1]).Msg("lookup IPs")
			return nil
		}
		for _, ip := range addrs {
			cidr, err := job.Tools().Network.ParseCIDR(ip.String())
			if err != nil {
				log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[1]).Msg("parse cidr")
				continue
			} else {
				entry := mmdbtype.Map{
					"lightswitch": mmdbtype.Map{
						"blacklist": mmdbtype.Map{
							"isTracking":     mmdbtype.Bool(true),
							"lastReportedAt": mmdbtype.String(lastReportedAt),
						},
					},
				}

				x.writer.Insert(cidr, entry)
			}
		}
		return nil
	})

	job.Tools().Items["lightswitch.blacklist.isTracking"] = Item{
		Item:        "ip/lightswitch.junk/blacklist.isTracking",
		Enabled:     true,
		GJSON:       "lightswitch.blacklist.isTracking",
		Description: "IP is an agressive tracker of identity, geo or ads.",
		Type:        common.Boolean.String()}
	job.Tools().Items["lightswitch.blacklist.lastReportedAt"] = Item{
		Item:        "ip/lightswitch.junk/blacklist.lastReportedAt",
		Enabled:     true,
		GJSON:       "lightswitch.blacklist.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}
	return err
}
