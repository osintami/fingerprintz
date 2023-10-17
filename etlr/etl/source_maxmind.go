// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"errors"
	"path/filepath"

	"github.com/osintami/fingerprintz/log"
)

type Maxmind struct{}

func NewMaxmind() *Maxmind {
	x := &Maxmind{}
	return x
}

func (x *Maxmind) Transform(job IETLJob) error {
	job.Tools().Items["ip.maxmind.location"] = Item{
		Item:        "ip/maxmind/location",
		Enabled:     true,
		GJSON:       "{\"city\":city.names.en,\"continent\":continent.code,\"country\":country.iso_code,\"latitude\":location.latitude,\"longitude\":location.longitude}",
		Description: "GEO",
		Type:        "JSON"}

	job.Tools().Items["ip.maxmind.country"] = Item{
		Item:        "ip/maxmind/country",
		Enabled:     true,
		GJSON:       "country.iso_code",
		Description: "ISO country code.",
		Type:        "String"}

	job.Tools().Items["ip.maxmind.longitude"] = Item{
		Item:        "ip/maxmind/longitude",
		Enabled:     true,
		GJSON:       "location.longitude",
		Description: "Longitude",
		Type:        "Float"}

	job.Tools().Items["ip.maxmind.timezone"] = Item{
		Item:        "ip/maxmind/timezone",
		Enabled:     true,
		GJSON:       "location.time_zone",
		Description: "Timezone",
		Type:        "String"}

	job.Tools().Items["ip.maxmind.isEU"] = Item{
		Item:        "ip/maxmind/isEU",
		Enabled:     true,
		GJSON:       "country.is_in_european_union",
		Description: "European Union",
		Type:        "Boolean"}

	job.Tools().Items["ip.maxmind.latitude"] = Item{
		Item:        "ip/maxmind/latitude",
		Enabled:     true,
		GJSON:       "location.latitude",
		Description: "Latitude",
		Type:        "Float"}

	job.Tools().Items["ip.maxmind.weatherCode"] = Item{
		Item:        "ip/maxmind/weatherCode",
		Enabled:     true,
		GJSON:       "location.weather_code",
		Description: "Weather Code",
		Type:        "String"}

	job.Tools().Items["ip.maxmind.city"] = Item{
		Item:        "ip/maxmind/city",
		Enabled:     true,
		GJSON:       "city.names.en",
		Description: "City name.",
		Type:        "String"}

	job.Tools().Items["ip.maxmind.continent"] = Item{
		Item:        "ip/maxmind/continent",
		Enabled:     true,
		GJSON:       "continent.code",
		Description: "Country name.",
		Type:        "String"}

	job.Tools().Items["ip.maxmind.subdivisionCod"] = Item{
		Item:        "ip/maxmind/subdivisionCode",
		Enabled:     true,
		GJSON:       "subdivisions.0.iso_code",
		Description: "ISO subdivision code.",
		Type:        "String"}

	return nil
}

func (x *Maxmind) Extract(job IETLJob) error {

	gzFile := "GeoLite2-City.tar.gz"
	url := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + job.Tools().Secrets.Find("MAXMIND_API_KEY") + "&suffix=tar.gz"
	if err := job.Tools().Network.DownloadFile(url, job.Info().workingPath+gzFile); err != nil {
		log.Error().Err(err).Str("component", "maxmind").Str("file", gzFile).Msg("download")
		return err
	}
	if err := job.Tools().FileSystem.UnGzipFile(job.Info().workingPath, gzFile, job.Info().workingPath); err != nil {
		log.Error().Err(err).Str("component", "maxmind").Str("file", gzFile).Msg("gzip -d")
		return err
	}

	tarFile := job.Info().workingPath + "GeoLite2-City.tar"
	if err := job.Tools().FileSystem.UnTarFile(tarFile, job.Info().workingPath); err != nil {
		// file already exists, nothing to see here, move along
		log.Error().Err(err).Str("component", "maxmind").Str("file", tarFile).Msg("remove")
		return err
	}

	return nil
}

func (x *Maxmind) Type() string {
	return "mmdb"
}

func (x *Maxmind) Create(string) error {
	return nil
}

func (x *Maxmind) Insert(key, value interface{}) error {
	return nil
}

var ErrFileMissing = errors.New("file missing")

func (x *Maxmind) Load(job IETLJob) error {
	mmdbFile := job.Info().workingPath + "GeoLite2-City_*/GeoLite2-City.mmdb"
	files, err := filepath.Glob(mmdbFile)
	if files == nil {
		log.Error().Err(err).Str("component", "maxmind").Str("file", mmdbFile).Msg("glob")
		return ErrFileMissing
	}

	err = job.Tools().FileSystem.Copy(files[len(files)-1], job.Info().snapshotFile)
	if err != nil {
		log.Error().Err(err).Str("component", "maxmind").Str("file", files[0]).Msg("copy")
		return ErrFileMissing
	}
	return nil
}
