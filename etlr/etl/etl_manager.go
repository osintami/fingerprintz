// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/osintami/fingerprintz/log"
	"github.com/robfig/cron/v3"
)

const (
	CRON_EVERY_MINUTE    = "* * * * *"
	CRON_EVERY_HALF_HOUR = "*/30 * * * *"
	CRON_EVERY_HOUR      = "0 * * * *"
	CRON_EVERY_DAY       = "0 0 * * *"
	CRON_EVERY_WEEK      = "0 0 * * 0"
)

type IETLManager interface {
	ScheduleCronJobs() *cron.Cron
	RefreshAll()
	Refresh(vendorName string) error
	FindJob(sourceName string) *ETLJob
	Source(sourceName string) *Source
}

type ETLManager struct {
	jobs     map[string]*ETLJob
	tools    *Toolbox
	dataPath string
}

func NewETLManager(tools *Toolbox, dataPath string, sources []Source) *ETLManager {
	x := &ETLManager{
		tools:    tools,
		dataPath: dataPath,
		jobs:     make(map[string]*ETLJob),
	}

	for _, source := range sources {
		// if source.Enabled {
		if job, err := x.createInstance(source); err == nil {
			x.jobs[source.Name] = job
		}
		// }
	}

	return x
}

func (x *ETLManager) ScheduleCronJobs() *cron.Cron {
	// scheduled data updates
	cron := cron.New()
	cron.AddFunc(
		CRON_EVERY_HOUR,
		x.refreshHourly)
	cron.AddFunc(
		CRON_EVERY_DAY,
		x.refreshDaily)
	cron.AddFunc(
		CRON_EVERY_WEEK,
		x.refreshWeekly)
	cron.Start()
	return cron
}

func (x *ETLManager) refreshHourly() {
	log.Debug().Str("component", "hourly cron").Msg("cron fired")
	x.Refresh("tormetrics")
	x.Refresh("onionoo")
	x.Refresh("danmeuk")
	x.Refresh("unwanted")
}

func (x ETLManager) refreshDaily() {
	log.Debug().Str("component", "daily cron").Msg("cron fired")
	x.Refresh("amazon")
	x.Refresh("digitalocean")
	x.Refresh("cloudflare")
	x.Refresh("google")
	x.Refresh("oracle")
	x.Refresh("abuseipdb")
	x.Refresh("avastel")
	x.Refresh("azure")
	x.Refresh("uhb")
	x.Refresh("droplist")
	x.Refresh("edroplist")
	x.Refresh("ipsum")
	x.Refresh("dshield")
	x.Refresh("x4bnet.tor")
	x.Refresh("x4bnet.vpn")
	x.Refresh("x4bnet.troll")
	x.Refresh("fakefilter")
	x.Refresh("ip1sms")
	x.Refresh("useragent")
	x.Refresh("udger.bot")
	x.Refresh("udger.vpn")
	x.Refresh("udger.proxy")
	x.Refresh("udger.cloud")
	x.Refresh("udger.crawler")
	x.Refresh("udger.http")
	x.Refresh("udger.mail")
	x.Refresh("udger.ssh")
	x.Refresh("udger.tor")
}

func (x ETLManager) refreshWeekly() {
	log.Debug().Str("component", "weekly cron").Msg("cron fired")
	x.Refresh("maxmind")
	x.Refresh("ip2location")
	x.Refresh("ip2proxy")
	x.Refresh("stripe-webhooks")
	x.Refresh("stripe-apis")
	x.Refresh("ip2location")
	x.Refresh("ip2proxy")
	x.Refresh("ayra")
}

func (x *ETLManager) RefreshAll() {
	for _, job := range x.jobs {
		job.Refresh()
	}
}

func (x *ETLManager) Refresh(vendorName string) error {
	job := x.jobs[vendorName]
	if job == nil {
		return ErrVendorNotFound
	}
	return job.Refresh()
}

func (x *ETLManager) FindJob(sourceName string) *ETLJob {
	return x.jobs[sourceName]
}

func (x *ETLManager) Source(sourceName string) *Source {
	job := x.jobs[sourceName]
	if job != nil {
		return job.Source()
	}
	return nil
}

func (x *ETLManager) createInstance(source Source) (*ETLJob, error) {

	var writer IWriter
	var extractor IExtract
	var transformer ITransform
	var loader ILoad

	// most common configuration
	switch source.OutputType {
	case "mmdb":
		writer = NewMMDBWriter()
		loader = writer
		// case "csv":
		// 	writer = NewCSVWriter()
		// 	loader = writer
	}

	if source.File != "" {
		extractor = NewFileExtractor(source.File)
	} else {
		extractor = NewHttpExtractor(x.tools, &source)
	}

	switch source.Name {
	case "lightswitch.junk":
		transformer = NewLightswitchJunk(writer)
	case "ayra":
		transformer = NewAyra(writer)
	case "test":
		transformer = NewTestSource(writer)
	case "tormetrics":
		transformer = NewTorMetrics(writer)
	case "amazon":
		transformer = NewAmazon(writer)
	case "digitalocean":
		transformer = NewDigitalOcean(writer)
	case "cloudflare":
		transformer = NewCloudFlare(writer)
	case "google":
		transformer = NewGoogle(writer)
	case "oracle":
		transformer = NewOracle(writer)
	case "abuseipdb":
		transformer = NewAbuseIPDB(writer)
	case "avastel":
		transformer = NewAvastel(writer)
	case "azure":
		transformer = NewAzure(writer)
	case "uhb":
		transformer = NewUHB(writer)
	case "danmeuk":
		transformer = NewDanMeUk(writer)
	case "droplist":
		transformer = NewDroplist(writer)
	case "edroplist":
		transformer = NewEdroplist(writer)
	case "onionoo":
		transformer = NewOnionOO(writer)
	case "ipsum":
		transformer = NewIpSUM(writer)
	case "dshield":
		transformer = NewDShield(writer)
	case "x4bnet.tor":
		transformer = NewX4BNetTOR(writer)
	case "x4bnet.vpn":
		transformer = NewX4BNetVPN(writer)
	case "x4bnet.troll":
		transformer = NewX4BNetTroll(writer)
	case "udger.bot":
		transformer = NewUdgerBot(writer)
	case "udger.cloud":
		transformer = NewUdgerCloud(writer)
	case "udger.proxy":
		transformer = NewUdgerProxy(writer)
	case "udger.crawler":
		transformer = NewUdgerCrawler(writer)
	case "udger.tor":
		transformer = NewUdgerTor(writer)
	case "udger.vpn":
		transformer = NewUdgerVPN(writer)
	case "udger.http":
		transformer = NewUdgerHTTP(writer)
	case "udger.ssh":
		transformer = NewUdgerSSH(writer)
	case "udger.mail":
		transformer = NewUdgerSMTP(writer)
	case "ip2location":
		transformer = NewIp2Location(writer)
	case "ip2proxy":
		transformer = NewIp2Proxy(writer)
	case "stripe.webhooks":
		transformer = NewStripeWebhooks(writer)
	case "stripe.apis":
		transformer = NewStripeAPIs(writer)
	case "fakefilter":
		fastdb := NewFastDBWriter()
		writer = fastdb
		transformer = NewFakefilter(writer)
		loader = fastdb
	case "ip1sms":
		fastdb := NewFastDBWriter()
		writer = fastdb
		transformer = NewIP1SMS(writer)
		loader = fastdb
	case "useragent":
		writer = NewFileDBWriter("yaml")
		transformer = NewUserAgent()
		loader = writer
	case "maxmind":
		maxmind := NewMaxmind()
		writer = maxmind
		extractor = maxmind
		transformer = maxmind
		loader = maxmind
	case "unwanted":
		transformer = NewUnwanted(writer)
	default:
		log.Error().Str("component", "etlr").Str("vendor", source.Name).Msg("ETLr not found")
		return nil, ErrVendorNotFound
	}

	return NewETLJob(
		x.tools,
		&source,
		x.dataPath,
		writer,
		extractor,
		transformer,
		loader), nil
}
