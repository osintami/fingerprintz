// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UdgerCloud struct {
	writer IWriter
}

func NewUdgerCloud(writer IWriter) ITransform {
	return &UdgerCloud{writer: writer}
}

func (x *UdgerCloud) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.CalculateCIDR(net.ParseIP(values[2]), net.ParseIP(values[3]))
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[2]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			url, err := url.Parse(values[1])
			if err != nil || url.Scheme == "" {
				log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[1]).Msg("parse website")
			}
			domain := strings.TrimPrefix(url.Hostname(), "www.")

			entry := mmdbtype.Map{
				"udger": mmdbtype.Map{
					"cloud": mmdbtype.Map{
						"isCloudNode":    mmdbtype.Bool(true),
						"provider":       mmdbtype.String(values[0]),
						"website":        mmdbtype.String(values[1]),
						"domain":         mmdbtype.String(domain),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["udger.cloud.isCloudNode"] = Item{
		Item:        "ip/udger.cloud/cloud.isCloudNode",
		Enabled:     true,
		GJSON:       "udger.cloud.isCloudNode",
		Description: "IP belongs to a cloud provider.",
		Type:        common.Boolean.String()}
	job.Tools().Items["udger.cloud.provider"] = Item{
		Item:        "ip/udger.cloud/cloud.provider",
		Enabled:     true,
		GJSON:       "udger.cloud.provider",
		Description: "Cloud provider name.",
		Type:        common.String.String()}
	job.Tools().Items["udger.cloud.website"] = Item{
		Item:        "ip/udger.cloud/cloud.website",
		Enabled:     true,
		GJSON:       "udger.cloud.website",
		Description: "Cloud provider website.",
		Type:        common.String.String()}
	job.Tools().Items["udger.cloud.domain"] = Item{
		Item:        "ip/udger.cloud/cloud.domain",
		Enabled:     true,
		GJSON:       "udger.cloud.domain",
		Description: "Cloud provider domain.",
		Type:        common.String.String()}
	job.Tools().Items["udger.cloud.lastReportedAt"] = Item{
		Item:        "ip/udger.cloud/cloud.lastReportedAt",
		Enabled:     true,
		GJSON:       "udger.cloud.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
