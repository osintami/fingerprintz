// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type CloudFlare struct {
	writer IWriter
}

func NewCloudFlare(writer IWriter) ITransform {
	return &CloudFlare{writer: writer}
}

func (x *CloudFlare) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"cloudflare": mmdbtype.Map{
					"cloud": mmdbtype.Map{
						"isCloudflare":   mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["cloudflare.cloud.isCloudflare"] = Item{
		Item:        "ip/cloudflare/cloud.isCloudflare",
		Enabled:     true,
		GJSON:       "cloudflare.cloud.isCloudflare",
		Description: "IP belongs to a cloud provider.",
		Type:        common.Boolean.String()}
	job.Tools().Items["azure.cloudflare.lastReportedAt"] = Item{
		Item:        "ip/cloudflare/cloud.lastReportedAt",
		Enabled:     true,
		GJSON:       "cloudflare.cloud.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
