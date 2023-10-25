// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type StripeWebhooks struct {
	writer IWriter
}

func NewStripeWebhooks(writer IWriter) ITransform {
	return &StripeWebhooks{writer: writer}
}

func (x *StripeWebhooks) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"stripe": mmdbtype.Map{
					"webhooks": mmdbtype.Map{
						"isStripe":       mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}
			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["stripe.webhooks.isStripe"] = Item{
		Item:        "ip/stripe.webhooks/isStripe",
		Enabled:     true,
		GJSON:       "stripe.webhooks.isStripe",
		Description: "IP belongs to Stripe.",
		Type:        common.Boolean.String()}
	job.Tools().Items["stripe.webhooks.lastReportedAt"] = Item{
		Item:        "ip/stripe.webhooks/lastReportedAt",
		Enabled:     true,
		GJSON:       "stripe.webhooks.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
