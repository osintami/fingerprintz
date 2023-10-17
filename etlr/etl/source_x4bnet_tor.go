// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type X4BNetTOR struct {
	writer IWriter
}

func NewX4BNetTOR(writer IWriter) ITransform {
	return &X4BNetTOR{writer: writer}
}

func (x *X4BNetTOR) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	// if job.Source().OutputType == "csv" {
	// 	x.writer.Insert(nil, []string{"# cidr", "lastReportedAt"})
	// }
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			switch job.Source().OutputType {
			case "mmdb":
				entry := mmdbtype.Map{
					"x4bnet": mmdbtype.Map{
						"tor": mmdbtype.Map{
							"isTorExitNode":  mmdbtype.Bool(true),
							"lastReportedAt": mmdbtype.String(lastReportedAt),
						},
					},
				}

				x.writer.Insert(cidr, entry)
				// case "csv":
				// 	row := []string{cidr.String(), lastReportedAt}
				// 	x.writer.Insert(nil, row)
			}
			return nil
		}
	})

	job.Tools().Items["x4bnet.tor.isTorExitNode"] = Item{
		Item:        "ip/x4bnet.tor/tor.isTorExitNode",
		Enabled:     true,
		GJSON:       "x4bnet.tor.isTorExitNode",
		Description: "IP belongs to a Tor Exit node.",
		Type:        common.Boolean.String()}
	job.Tools().Items["x4bnet.tor.lastReportedAt"] = Item{
		Item:        "ip/x4bnet.tor/tor.lastReportedAt",
		Enabled:     true,
		GJSON:       "x4bnet.tor.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
