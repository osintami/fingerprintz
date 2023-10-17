// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Oracle struct {
	writer IWriter
}

func NewOracle(writer IWriter) ITransform {
	return &Oracle{writer: writer}
}

func (x *Oracle) Transform(job IETLJob) error {

	type Cidrs struct {
		Cidr string   `json:"cidr"`
		Tags []string `json:"tags"`
	}
	type Regions struct {
		Region string  `json:"region"`
		Cidrs  []Cidrs `json:"cidrs"`
	}
	type OracleEntry struct {
		LastUpdatedTimestamp string    `json:"last_updated_timestamp"`
		Regions              []Regions `json:"regions"`
	}

	entry := &OracleEntry{}
	err := common.LoadJson(job.Info().inputFile, entry)
	for _, region := range entry.Regions {
		for _, ip := range region.Cidrs {
			cidr, err := job.Tools().Network.ParseCIDR(ip.Cidr)
			if err != nil {
				log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", ip.Cidr).Msg("parse cidr")
				continue
			} else {
				entry := mmdbtype.Map{
					"oracle": mmdbtype.Map{
						"cloud": mmdbtype.Map{
							"isOracle":       mmdbtype.Bool(true),
							"lastReportedAt": mmdbtype.String(entry.LastUpdatedTimestamp),
							"region":         mmdbtype.String(region.Region),
							"product":        mmdbtype.String(ip.Tags[0]),
						},
					},
				}

				x.writer.Insert(cidr, entry)
			}
		}
	}

	job.Tools().Items["oracle.cloud.isOracle"] = Item{
		Item:        "ip/oracle/cloud.isOracle",
		Enabled:     true,
		GJSON:       "oracle.cloud.isOracle",
		Description: "IP belongs to a cloud provider.",
		Type:        common.Boolean.String()}
	job.Tools().Items["oracle.cloud.region"] = Item{
		Item:        "ip/oracle/cloud.region",
		Enabled:     true,
		GJSON:       "oracle.cloud.region",
		Description: "Oracle region.",
		Type:        common.String.String()}
	job.Tools().Items["oracle.cloud.product"] = Item{
		Item:        "ip/oracle/cloud.product",
		Enabled:     true,
		GJSON:       "oracle.cloud.product",
		Description: "Oracle product.",
		Type:        common.String.String()}
	job.Tools().Items["oracle.cloud.lastReportedAt"] = Item{
		Item:        "ip/oracle/cloud.lastReportedAt",
		Enabled:     true,
		GJSON:       "oracle.cloud.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
