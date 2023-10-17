// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type UdgerCrawler struct {
	writer IWriter
}

func NewUdgerCrawler(writer IWriter) ITransform {
	return &UdgerCrawler{writer: writer}
}

func (x *UdgerCrawler) Transform(job IETLJob) error {
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().inputFile, rune(job.Source().Separator[0]), func(values []string) error {
		cidr, err := job.Tools().Network.ParseCIDR(values[0])
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"udger": mmdbtype.Map{
					"crawler": mmdbtype.Map{
						"isCrawler":      mmdbtype.Bool(true),
						"useragent":      mmdbtype.String(values[0]),
						"host":           mmdbtype.String(values[2]),
						"city":           mmdbtype.String(values[4]),
						"countryCode":    mmdbtype.String(values[3]),
						"family":         mmdbtype.String(values[5]),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
			return nil
		}
	})

	job.Tools().Items["udger.crawler.isCrawler"] = Item{
		Item:        "ip/udger/crawler.isCrawler",
		Enabled:     true,
		GJSON:       "udger.crawler.isCrawler",
		Description: "IP is associated with a crawler.",
		Type:        common.Boolean.String()}
	job.Tools().Items["udger.crawler.useragent"] = Item{
		Item:        "ip/udger/crawler.useragent",
		Enabled:     true,
		GJSON:       "udger.crawler.useragent",
		Description: "User-Agent for known crawler.",
		Type:        common.String.String()}
	job.Tools().Items["udger.crawler.host"] = Item{
		Item:        "ip/udger/crawler.host",
		Enabled:     true,
		GJSON:       "udger.crawler.host",
		Description: "Hostname.",
		Type:        common.String.String()}
	job.Tools().Items["udger.crawler.city"] = Item{
		Item:        "ip/udger/crawler.city",
		Enabled:     true,
		GJSON:       "udger.crawler.city",
		Description: "GEO city.",
		Type:        common.String.String()}
	job.Tools().Items["udger.crawler.countryCode"] = Item{
		Item:        "ip/udger/crawler.countryCode",
		Enabled:     true,
		GJSON:       "udger.crawler.countryCode",
		Description: "GEO country code.",
		Type:        common.String.String()}
	job.Tools().Items["udger.crawler.lastReportedAt"] = Item{
		Item:        "ip/udger/crawler.lastReportedAt",
		Enabled:     true,
		GJSON:       "udger.crawler.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
