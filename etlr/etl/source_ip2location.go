// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import "github.com/osintami/fingerprintz/common"

type Ip2Location struct{ writer IWriter }

func NewIp2Location(writer IWriter) ITransform {
	return &Ip2Location{writer: writer}
}

func (x *Ip2Location) Transform(job IETLJob) error {
	Ip2LocationOrProxyTransform(job, x.writer, "IP2LOCATION-LITE-DB1.CSV", "LICENSE-CC-BY-SA-4.0.TXT", false)

	job.Tools().Items["ip2location.country.countryCode"] = Item{
		Item:        "ip/ip2location/country.countryCode",
		Enabled:     true,
		GJSON:       "ip2location.country.countryCode",
		Description: "GEO country code.",
		Type:        common.Boolean.String()}

	return nil
}
