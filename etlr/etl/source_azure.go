// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Azure struct{ writer IWriter }

func NewAzure(writer IWriter) ITransform {
	return &Azure{writer: writer}
}

func (x *Azure) Transform(job IETLJob) error {

	type Properties struct {
		ChangeNumber    int      `json:"changeNumber"`
		Region          string   `json:"region"`
		RegionID        int      `json:"regionId"`
		Platform        string   `json:"platform"`
		SystemService   string   `json:"systemService"`
		AddressPrefixes []string `json:"addressPrefixes"`
	}

	type Values struct {
		Name       string     `json:"name"`
		ID         string     `json:"id"`
		Properties Properties `json:"properties"`
	}

	type AzureEntry struct {
		ChangeNumber int      `json:"changeNumber"`
		Cloud        string   `json:"cloud"`
		Values       []Values `json:"values"`
	}

	entry := &AzureEntry{}
	common.LoadJson(job.Info().inputFile, entry)
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	for _, product := range entry.Values {
		for _, ip := range product.Properties.AddressPrefixes {
			cidr, err := job.Tools().Network.ParseCIDR(ip)
			if err != nil {
				log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", ip).Msg("parse cidr")
				continue
			} else {
				entry := mmdbtype.Map{
					"microsoft": mmdbtype.Map{
						"cloud": mmdbtype.Map{
							"isAzure":        mmdbtype.Bool(true),
							"lastReportedAt": mmdbtype.String(lastReportedAt),
							"product":        mmdbtype.String(product.Name),
						},
					},
				}

				x.writer.Insert(cidr, entry)
			}
		}
	}

	job.Tools().Items["microsoft.cloud.isAzure"] = Item{
		Item:        "ip/azure/cloud.isAzure",
		Enabled:     true,
		GJSON:       "microsoft.cloud.isAzure",
		Description: "IP belongs to a cloud provider.",
		Type:        common.Boolean.String()}
	job.Tools().Items["microsoft.cloud.lastReportedAt"] = Item{
		Item:        "ip/azure/cloud.lastReportedAt",
		Enabled:     true,
		GJSON:       "microsoft.cloud.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}
	job.Tools().Items["microsoft.cloud.product"] = Item{
		Item:        "ip/azure/cloud.product",
		Enabled:     true,
		GJSON:       "microsoft.cloud.product",
		Description: "Azure product.",
		Type:        common.String.String()}

	return nil
}
