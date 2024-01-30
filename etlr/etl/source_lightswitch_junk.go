// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"net"
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type LightswitchJunk struct {
	writer IWriter
}

func NewLightswitchJunk(writer IWriter) ITransform {
	return &LightswitchJunk{writer: writer}
}

func (x *LightswitchJunk) Transform(job IETLJob) error {
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
			cidr := &net.IPNet{
				IP:   ip,
				Mask: ip.DefaultMask(),
			}
			entry := mmdbtype.Map{
				"lightswitch": mmdbtype.Map{
					"blacklist": mmdbtype.Map{
						"isHateOrJunk":   mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
		}
		return nil
	})

	job.Tools().Items["lightswitch.blacklist.isHateOrJunk"] = Item{
		Item:        "ip/lightswitch/blacklist.isHateOrJunk",
		Enabled:     true,
		GJSON:       "lightswitch.blacklist.isHateOrJunk",
		Description: "IP is hate or junk related.",
		Type:        common.Boolean.String()}
	job.Tools().Items["lightswitch.blacklist.lastReportedAt"] = Item{
		Item:        "ip/lightswitch/blacklist.lastReportedAt",
		Enabled:     true,
		GJSON:       "lightswitch.junk.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}
	return err
}
