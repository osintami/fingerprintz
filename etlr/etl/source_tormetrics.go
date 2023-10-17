// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type TorMetrics struct {
	writer IWriter
}

func NewTorMetrics(writer IWriter) ITransform {
	return &TorMetrics{writer: writer}
}

func (x *TorMetrics) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {

			switch x.writer.Type() {
			case "mmdb":
				entry := mmdbtype.Map{
					"tormetrics": mmdbtype.Map{
						"tor": mmdbtype.Map{
							"isTorExitNode":  mmdbtype.Bool(true),
							"lastReportedAt": mmdbtype.String(lastReportedAt),
						},
					},
				}
				x.writer.Insert(cidr, entry)
			}
			return nil
		}
	})

	job.Tools().Items["tormetrics.tor.isTorExitNode"] = Item{
		Item:        "ip/tormetrics/tor.isTorExitNode",
		Enabled:     true,
		GJSON:       "tormetrics.tor.isTorExitNode",
		Description: "IP belongs to a Tor Exit node.",
		Type:        common.Boolean.String()}

	job.Tools().Items["tormetrics.tor.lastReportedAt"] = Item{
		Item:        "ip/tormetrics/tor.lastReportedAt",
		Enabled:     true,
		GJSON:       "tormetrics.tor.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
