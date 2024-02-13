// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"
	"strings"
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Akamai struct {
	writer IWriter
}

func NewAkamai(writer IWriter) ITransform {
	return &Akamai{writer: writer}
}

func (x *Akamai) Transform(job IETLJob) error {
	if !strings.HasPrefix(job.Source().URL, "-") && job.Source().File == "" {
		os.Remove(job.Info().workingPath + "akamai_ipv4_CIDRs.txt")
		os.Remove(job.Info().workingPath + "akamai_ipv6_CIDRs.txt")
		os.RemoveAll(job.Info().workingPath + "__MACOSX")
		if err := job.Tools().FileSystem.UnzipFile(job.Info().inputFile, job.Info().workingPath); err != nil {
			return err
		}
		os.RemoveAll(job.Info().workingPath + "__MACOSX")
	}
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)

	// IPv4 file
	err := job.Tools().CSV.ProcessFile(job.Info().workingPath+"akamai_ipv4_CIDRs.txt", ' ', func(values []string) error {
		value := strings.Replace(values[0], "\t", "", -1)
		cidr, err := job.Tools().Network.ParseCIDR(value)
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"akamai": mmdbtype.Map{
					"cloud": mmdbtype.Map{
						"isAkamai":       mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
		}
		return nil
	})

	// IPv6 file
	job.Tools().CSV.ProcessFile(job.Info().workingPath+"akamai_ipv6_CIDRs.txt", ' ', func(values []string) error {
		value := strings.Replace(values[0], "\t", "", -1)
		cidr, err := job.Tools().Network.ParseCIDR(value)
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			entry := mmdbtype.Map{
				"akamai": mmdbtype.Map{
					"cloud": mmdbtype.Map{
						"isAkamai":       mmdbtype.Bool(true),
						"lastReportedAt": mmdbtype.String(lastReportedAt),
					},
				},
			}

			x.writer.Insert(cidr, entry)
		}
		return nil
	})

	job.Tools().Items["akamai.cloud.isAkamai"] = Item{
		Item:        "ip/akamai/cloud.isAkamai",
		Enabled:     true,
		GJSON:       "akamai.cloud.isAkamai",
		Description: "IP belongs to a cloud provider.",
		Type:        common.Boolean.String()}
	job.Tools().Items["akamai.cloud.lastReportedAt"] = Item{
		Item:        "ip/akamai/cloud.lastReportedAt",
		Enabled:     true,
		GJSON:       "akamai.cloud.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}
