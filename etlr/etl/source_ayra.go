// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Ayra struct {
	writer IWriter
	tiers  map[string]*ItemMapper
}

type ItemMapper struct {
	Layer1     string
	Layer2     string `json:",omitempty"`
	Layer3     string `json:",omitempty"`
	IsItemName string `json:",omitempty"`
}

func NewAyra(writer IWriter) ITransform {

	x := &Ayra{writer: writer, tiers: make(map[string]*ItemMapper)}

	// NOTE:  country, amazon, azure, and cloudflare are disabled because
	//   the GEO data is taken directly from the free IP2Location data set
	//   and the cloud data is collected from the cloud providers directly
	// NOTE:  tor data is from tormetrics, so is ASN
	x.tiers["asn"] = &ItemMapper{Layer1: "asn", Layer2: "ASO", Layer3: "ASN"}

	return x
}

func (x *Ayra) Transform(job IETLJob) error {
	data := make(map[string]map[string]map[string][]string)
	common.LoadJson(job.Info().inputFile, &data)

	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	for category, tier1Map := range data {
		tier := x.tiers[category]
		for tier1, tier2Map := range tier1Map {
			for tier2, ipList := range tier2Map {
				for _, ip := range ipList {

					cidr, err := job.Tools().Network.ParseCIDR(ip)
					if err != nil {
						log.Error().Err(err).Str("component", job.Source().Name).Str("ip", ip).Msg("parse cidr")
						continue
					} else {
						itemMap := mmdbtype.Map{}
						if tier.Layer2 != "" {
							itemMap[mmdbtype.String(tier.Layer2)] = mmdbtype.String(tier1)
						}
						if tier.Layer3 != "" {
							itemMap[mmdbtype.String(tier.Layer3)] = mmdbtype.String(tier2)
						}
						if tier1 == "__" && tier2 == "__" {
							itemMap[mmdbtype.String("isUnassigned")] = mmdbtype.Bool(true)
						}
						itemMap[mmdbtype.String("lastReportedAt")] = mmdbtype.String(lastReportedAt)
						entry := mmdbtype.Map{
							mmdbtype.String(job.Source().Name): mmdbtype.Map{
								mmdbtype.String(tier.Layer1): itemMap}}

						x.writer.Insert(cidr, entry)
					}
				}
			}
		}
	}

	job.Tools().Items["ayra.asn.ASN"] = Item{
		Item:        "ip/ayra/asn.ASN",
		Enabled:     true,
		GJSON:       "ayra.asn.ASN",
		Description: "ASN.",
		Type:        common.String.String()}
	job.Tools().Items["ayra.asn.isUnassigned"] = Item{
		Item:        "ip/ayra/asn.iUnassigned",
		Enabled:     true,
		GJSON:       "ayra.asn.isUnassigned",
		Description: "Unassigned ASN/ASO.",
		Type:        common.Boolean.String()}
	job.Tools().Items["ayra.asn.ASO"] = Item{
		Item:        "ip/ayra/asn.ASO",
		Enabled:     true,
		GJSON:       "ayra.asn.ASO",
		Description: "ASO.",
		Type:        common.String.String()}
	job.Tools().Items["ayra.asn.lastReportedAt"] = Item{
		Item:        "ip/ayra/asn.lastReportedAt",
		Enabled:     true,
		GJSON:       "ayra.asn.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return nil
}
