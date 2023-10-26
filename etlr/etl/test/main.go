// Copyright © 2023 OSINTAMI. This is not yours.
package main

import (
	"crypto/tls"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/etlr/etl"
	"github.com/osintami/fingerprintz/etlr/server"
	"github.com/osintami/fingerprintz/etlr/utils"
	"github.com/osintami/fingerprintz/log"
)

func main() {

	// if len(os.Args) < 2 {
	// 	fmt.Println("USAGE: go run . source")
	// 	return
	// }
	// source := os.Args[1]

	source := "unwanted"

	svrConfig := &server.ServerConfig{}
	common.LoadEnv(true, true, svrConfig)
	log.InitLogger(svrConfig.LogPath, "jobs.log", svrConfig.LogLevel, false)
	if svrConfig.LogLevel == "TRACE" {
		common.PrintEnvironment()
	}

	// create tools object used by everyone, done this way to make writing tests easier
	resty := resty.New()
	// HACK:  ultimate hosts blacklist's mirror server has an expired certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	resty = resty.SetTransport(tr)
	tools := etl.Toolbox{
		Network:    utils.NewNetworkingHelper(resty),
		FileSystem: utils.NewFSHelper(),
		Secrets:    common.NewSecrets([]string{"ABUSEIPDB_API_KEY", "IP2LOCATION_API_KEY", "MAXMIND_API_KEY", "MAXMIND_ACCOUNT", "UDGER_API_KEY"}),
		CSV:        utils.NewCSVReader(),
		Items:      make(map[string]etl.Item)}

	// load all the ETL instructions organized by vendor
	sources := []etl.Source{}
	// load production data sources
	err := common.LoadJson("./config.json", &sources)
	//err := common.LoadJson("udger.json", &sources)
	if err != nil {
		log.Fatal().Err(err).Str("component", "etlr").Msg("config.json missing or corrupt")
		return
	}
	manager := etl.NewETLManager(&tools, svrConfig.DataPath, sources)

	job := manager.FindJob(source)
	job.Refresh()
}
