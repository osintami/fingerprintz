// Copyright © 2023 OSINTAMI. This is not yours.
package main

import (
	"crypto/tls"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/etlr/etl"
	"github.com/osintami/fingerprintz/etlr/server"
	"github.com/osintami/fingerprintz/etlr/utils"
	"github.com/osintami/fingerprintz/log"
)

func main() {

	svrConfig := &server.ServerConfig{}
	common.LoadEnv(true, true, svrConfig)
	log.InitLogger(svrConfig.LogPath, "etlr.log", svrConfig.LogLevel, false)
	if svrConfig.LogLevel == "TRACE" {
		common.PrintEnvironment()
	}

	resty := resty.New()
	// HACK: UHB's mirror server has an expired SSL certificate
	resty = resty.SetTransport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}})

	tools := etl.Toolbox{
		Network:    utils.NewNetworkingHelper(resty),
		FileSystem: utils.NewFSHelper(),
		Secrets:    common.NewSecrets([]string{"ABUSEIPDB_API_KEY", "IP2LOCATION_API_KEY", "MAXMIND_API_KEY", "MAXMIND_ACCOUNT", "UDGER_API_KEY"}),
		CSV:        utils.NewCSVReader(),
		Items:      make(map[string]etl.Item)}

	// load the ETL instructions
	sources := []etl.Source{}
	err := common.LoadJson("config.json", &sources)
	if err != nil {
		log.Fatal().Err(err).Str("component", "etlr").Msg("config.json missing or corrupt")
		return
	}

	manager := etl.NewETLManager(&tools, svrConfig.DataPath, sources)

	handlers := server.NewETLrServer(manager)
	router := chi.NewMux()
	router.Route(svrConfig.PathPrefix, func(r chi.Router) {
		r.Get("/v1/refresh/{vendor}", handlers.RefreshHandler)
	})

	// schedule source data refreshes
	manager.ScheduleCronJobs()

	err = common.ListenAndServe(svrConfig.ListenAddr, "", "", router)
	if err != nil {
		log.Fatal().Err(err).Str("component", "etlr").Str("state", "stopped").Msg("orchestration")
	} else {
		log.Info().Str("component", "etlr").Str("state", "stopped").Msg("orchestration")
	}
}
