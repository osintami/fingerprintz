// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Amazon struct {
	writer IWriter
}

func NewAmazon(writer IWriter) ITransform {
	return &Amazon{writer: writer}
}

func (x *Amazon) Transform(job IETLJob) error {

	type AmazonService struct {
		CIDR    string `json:"ip_prefix"`
		Region  string
		Service string
		Group   string `json:"network_border_group"`
	}

	type AmazonEntry struct {
		SyncToken      string
		LastReportedAt string `json:"createDate"`
		Prefixes       []AmazonService
	}

	entry := &AmazonEntry{}
	err := common.LoadJson(job.Info().inputFile, entry)
	for _, item := range entry.Prefixes {
		cidr, err := job.Tools().Network.ParseCIDR(item.CIDR)
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", item.CIDR).Msg("parse cidr")
			continue
		} else {
			entry := mmdbtype.Map{
				"aws": mmdbtype.Map{
					"cloud": mmdbtype.Map{
						"isAmazon":       mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(entry.LastReportedAt),
						"region":         mmdbtype.String(item.Region),
						"product":        mmdbtype.String(item.Service),
					},
				},
			}

			x.writer.Insert(cidr, entry)
		}
	}

	job.Tools().Items["amazon.cloud.isAmazon"] = Item{
		Item:        "ip/amazon/cloud.isAmazon",
		Enabled:     true,
		GJSON:       "aws.cloud.isAmazon",
		Description: "IP belongs to a cloud provider.",
		Type:        common.Boolean.String()}
	job.Tools().Items["amazon.cloud.lastReportedAt"] = Item{
		Item:        "ip/amazon/cloud.lastReportedAt",
		Enabled:     true,
		GJSON:       "aws.cloud.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}
	job.Tools().Items["amazon.cloud.region"] = Item{
		Item:        "ip/amazon/cloud.region",
		Enabled:     true,
		GJSON:       "aws.cloud.region",
		Description: "AWS region.",
		Type:        common.String.String()}
	job.Tools().Items["amazon.cloud.product"] = Item{
		Item:        "ip/amazon/cloud.region",
		Enabled:     true,
		GJSON:       "aws.cloud.region",
		Description: "AWS product.",
		Type:        common.String.String()}

	return err
}
