// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/ua-parser/uap-go/uaparser"
)

func buildETLManager() *ETLManager {

	sources := []Source{}
	// load test data sources
	err := common.LoadJson("./test/config.json", &sources)
	if err != nil {
		return nil
	}

	// pre-clean
	for _, source := range sources {
		cleanETLFragments(source.Name)
	}

	return NewETLManager(fillToolbox(nil), "./test/data/", sources)
}

func cleanETLFragments(sourceName string) {
	os.RemoveAll(fmt.Sprintf("/tmp/%s/", sourceName))
	os.Remove(fmt.Sprintf("./test/data/%s.mmdb", sourceName))
	os.Remove(fmt.Sprintf("./test/data/%s.json", sourceName))
}

func TestMaxmindSource(t *testing.T) {

	source := &Source{
		Name:       "maxmind",
		Enabled:    true,
		InputType:  "mmdb",
		OutputType: "mmdb"}

	writer := NewMaxmind()
	extract := writer
	transform := writer
	load := writer

	client := resty.New()
	tools := fillToolbox(client)
	tools.Secrets.Set("MAXMIND_API_KEY", "test-api-key")

	job := NewETLJob(tools, source, "/tmp/", writer, extract, transform, load)
	assert.Equal(t, "mmdb", job.writer.Type())

	// EXTRACT

	// read maxmind download simulated tiny data
	content, err := os.ReadFile("./test/GeoLite2-City.tar.gz")
	assert.Nil(t, err)

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	url := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=test-api-key&suffix=tar.gz"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewBytesResponder(http.StatusOK, content))

	err = job.Extract()
	assert.Nil(t, err)

	// simulate download fail
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewBytesResponder(http.StatusNotFound, content))

	err = job.Extract()
	assert.NotNil(t, err)

	// fail on GZIP -d with bad content
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewBytesResponder(http.StatusOK, []byte("nope")))

	err = job.Extract()
	assert.NotNil(t, err)

	// TODO:  fail on TAR -xvf with bad content
	content, err = os.ReadFile("./test/nope.tar.gz")
	assert.Nil(t, err)

	httpmock.RegisterResponder(
		"GET", url, httpmock.NewBytesResponder(http.StatusOK, content))

	err = job.Extract()
	assert.NotNil(t, err)

	// TRANSFORM (copy creates schema)
	err = job.Transform()
	assert.Nil(t, err)
	assert.Equal(t, 10, len(job.tools.Items))
	assert.Equal(t, "ip/maxmind/location", job.tools.Items["ip.maxmind.location"].Item)

	// NOTE: these are transform no-ops as there isn't a real transform with this data set, only schema creation
	job.writer.Create("any")
	job.writer.Insert(nil, nil)

	// LOAD
	err = job.Load()
	assert.Nil(t, err)

	// force glob to fail
	job.info.workingPath = "/tmp/maxmind/nope/"
	err = job.Load()
	assert.Equal(t, ErrFileMissing, err)

	// force copy to fail
	job.info.workingPath = "/tmp/maxmind/"
	job.info.snapshotFile = "./..///\\/tmp/maxmind/nope.mmdb"
	err = job.Load()
	assert.Equal(t, ErrFileMissing, err)

	// cleanup
	os.RemoveAll("/tmp/maxmind")
}

func TestSources(t *testing.T) {
	// NOTE:  as sample data sets are collected and published to ./test/source/{vendor}.{type}
	//   add them here and they will automatically be tested
	sources := []string{
		"test",
		"amazon",
		"abuseipdb",
		"avastel",
		"cloudflare",
		"azure",
		"danmeuk",
		"digitalocean",
		"droplist",
		"edroplist",
		"dshield",
		"fakefilter",
		"google",
		"ip1sms",
		"ipsum",
		"uhb",
		"oracle",
		"onionoo",
		"x4bnet.tor",
		"x4bnet.vpn",
		"x4bnet.troll",
		"stripe.apis",
		"stripe.webhooks",
		"tormetrics",
		"ip2proxy",
		"ip2location",
		"udger.bot",
		"udger.cloud",
		"udger.proxy",
		"udger.crawler",
		"udger.tor",
		"udger.vpn",
		"udger.http",
		"udger.ssh",
		"udger.mail",
		"unwanted",
		"useragent",
		"ayra",
		// TODO:  can't pass the 1.2.3.4 lookup test because these are domains that are resolved
		"lightswitch.junk",
		"lightswitch.aggressive",
		// TODO:  can't pass because it parses two files instead of one
		//"akamai",
	}

	manager := buildETLManager()
	if manager == nil {
		assert.FailNow(t, "ETL manager failed to instantiate")
	}

	for _, sourceName := range sources {
		// ETL this data source
		job := manager.FindJob(sourceName)
		assert.NotNil(t, job)
		if !job.Source().Enabled {
			continue
		}
		err := job.Refresh()
		assert.Nil(t, err)

		// read a data point
		var data json.RawMessage

		// create the proper reader based on database type
		switch manager.Source(sourceName).OutputType {
		case "mmdb":
			mmdb, err := common.NewMaxmindReader("./test/data/" + sourceName + ".mmdb")
			assert.Nil(t, err)
			data, err = mmdb.Lookup(net.ParseIP("1.2.3.4"))
			assert.Nil(t, err)

		case "fast":
			fast := common.NewFastCache()
			fast.LoadFile("./test/data/" + sourceName + ".fast")
			if sourceName == "fakefilter" {
				obj, found := fast.Get("nope.com")
				assert.True(t, found)
				data = obj.([]uint8)
			} else if sourceName == "ip1sms" {
				obj, found := fast.Get("15121234567")
				assert.True(t, found)
				data = obj.([]uint8)
			}

		case "yaml":
			ua, err := uaparser.New("./test/data/" + sourceName + ".yaml")
			assert.Nil(t, err)
			uaInfo := ua.Parse(strings.ReplaceAll("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36", "%20", " "))
			data, err = json.Marshal(uaInfo)
			assert.Nil(t, err)
		case "csv":
			fmt.Println("OUTPUT CSV", sourceName)
		}

		// test each schema item for this ETL job
		for _, item := range manager.tools.Items {
			result := gjson.GetBytes(data, item.GJSON)
			if !result.Exists() {
				log.Error().Str("component", "test").Str("source", sourceName).Str("schema", item.GJSON).Str("data", string(data)).Msg("schema or data is banged up")
			}
			assert.True(t, result.Exists())
		}

		if false {
			cleanETLFragments(sourceName)
		}
	}
}
