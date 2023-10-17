// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type Ip2Proxy struct {
	writer IWriter
}

func NewIp2Proxy(writer IWriter) ITransform {
	return &Ip2Proxy{writer: writer}
}

func (x *Ip2Proxy) Transform(job IETLJob) error {
	err := Ip2LocationOrProxyTransform(job, x.writer, "IP2PROXY-LITE-PX2.CSV", "LICENSE-CC-BY-4.0.TXT", true)

	job.Tools().Items["ip2proxy.proxy.isProxy"] = Item{
		Item:        "ip/ip2proxy/proxy.isProxy",
		Enabled:     true,
		GJSON:       "ip2proxy.proxy.isProxy",
		Description: "IP has a proxy associated with it.",
		Type:        common.Boolean.String()}
	job.Tools().Items["ip2proxy.proxy.lastReportedAt"] = Item{
		Item:        "ip/ip2proxy/proxy.lastReportedAt",
		Enabled:     true,
		GJSON:       "ip2proxy.proxy.lastReportedAt",
		Description: "Last seen date/time.",
		Type:        common.Date.String()}

	return err
}

func Ip2LocationOrProxyTransform(job IETLJob, writer IWriter, csvFile, licenseFile string, ip2proxy bool) error {
	// unzip fails if any of the files exist
	if !strings.HasPrefix(job.Source().URL, "-") {
		os.Remove(job.Info().workingPath + csvFile)
		os.Remove(job.Info().workingPath + licenseFile)
		os.Remove(job.Info().workingPath + "README_LITE.TXT")
		if err := job.Tools().FileSystem.UnzipFile(job.Info().inputFile, job.Info().workingPath); err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("file", job.Info().inputFile).Msg("unzip")
			if job.Source().URL != "" && strings.HasPrefix(job.Source().URL, "-") || job.Source().File != "" {
				// NOTE:  error is okay in this condition, as we are doing a run without a fresh download,
				//   copy the existing CSV file from the source location to the working directory
				err := job.Tools().FileSystem.Copy(job.Source().File, job.Info().workingPath+csvFile)
				if err != nil {
					log.Error().Err(err).Str("component", job.Source().Name).Str("file", job.Info().inputFile).Msg("unzip")
				}
			} else {
				return err
			}
		}
		os.Remove(job.Info().workingPath + licenseFile)
		os.Remove(job.Info().workingPath + "README_LITE.TXT")
		os.Chmod(job.Info().workingPath+csvFile, 0644)
	}
	lastReportedAt := time.Now().Format(common.GO_DEFAULT_DATE)
	err := job.Tools().CSV.ProcessFile(job.Info().workingPath+csvFile, rune(job.Source().Separator[0]), func(values []string) error {
		iIp1, _ := strconv.ParseUint(values[0], 10, 32)
		ip1 := job.Tools().Network.Uint2IPv4(iIp1).To4()
		if ip1 == nil {
			// some of these are IPv6
			ip1 = job.Tools().Network.Uint2IPv4(iIp1)
		}
		iIp2, _ := strconv.ParseUint(values[1], 10, 32)
		ip2 := job.Tools().Network.Uint2IPv4(iIp2).To4()
		if ip2 == nil {
			// some of these are IPv6
			ip2 = job.Tools().Network.Uint2IPv4(iIp2)
		}

		cidr, err := job.Tools().Network.CalculateCIDR(ip1, ip2)
		if err != nil {
			log.Error().Err(err).Str("component", job.Source().Name).Str("cidr", values[0]).Msg("parse cidr")
			return ErrBadSourceData
		} else {
			var entry mmdbtype.Map
			if ip2proxy {
				entry = mmdbtype.Map{
					"ip2proxy": mmdbtype.Map{
						"proxy": mmdbtype.Map{
							"isProxy":        mmdbtype.Bool(true),
							"lastReportedAt": mmdbtype.String(lastReportedAt),
						},
					},
				}
			} else {
				entry = mmdbtype.Map{
					"ip2location": mmdbtype.Map{
						"country": mmdbtype.Map{
							"countryCode": mmdbtype.String(values[2]),
						},
					},
				}
			}

			writer.Insert(cidr, entry)
		}
		return nil
	})
	return err
}
