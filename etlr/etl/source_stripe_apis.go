// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type StripeAPIs struct {
	writer IWriter
}

func NewStripeAPIs(writer IWriter) ITransform {
	return &StripeAPIs{writer: writer}
}

func (x *StripeAPIs) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"stripe": mmdbtype.Map{
					"apis": mmdbtype.Map{
						"isStripe":       mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}
			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["stripe.apis.isStripe"] = Item{
		Item:        "ip/stripe.apis/isStripe",
		Enabled:     true,
		GJSON:       "stripe.apis.isStripe",
		Description: "IP belongs to Stripe.",
		Type:        common.Boolean.String()}
	job.Tools().Items["stripe.apis.lastReportedAt"] = Item{
		Item:        "ip/stripe.apis/lastReportedAt",
		Enabled:     true,
		GJSON:       "stripe.apis.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
