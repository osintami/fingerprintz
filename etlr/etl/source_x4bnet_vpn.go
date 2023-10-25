// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type X4BNetVPN struct {
	writer IWriter
}

func NewX4BNetVPN(writer IWriter) ITransform {
	return &X4BNetVPN{writer: writer}
}

func (x *X4BNetVPN) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"x4bnet": mmdbtype.Map{
					"vpn": mmdbtype.Map{
						"isVPN":          mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["x4bnet.vpn.isVPN"] = Item{
		Item:        "ip/x4bnet.vpn/vpn.isVPN",
		Enabled:     true,
		GJSON:       "x4bnet.vpn.isVPN",
		Description: "IP has a VPN associated with it.",
		Type:        common.Boolean.String()}
	job.Tools().Items["x4bnet.vpn.lastReportedAt"] = Item{
		Item:        "ip/x4bnet.vpn/vpn.lastReportedAt",
		Enabled:     true,
		GJSON:       "x4bnet.vpn.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
