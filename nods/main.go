// Copyright © 2023 OSINTAMI. This is not yours.
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-resty/resty/v2"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/osintami/fingerprintz/nods/server"
)

func main() {
	svrConfig := &server.ServerConfig{}
	common.LoadEnv(true, true, svrConfig)

	log.InitLogger(svrConfig.LogPath, "nods.log", svrConfig.LogLevel, false)
	if svrConfig.LogLevel == "TRACE" {
		common.PrintEnvironment()
	}

	watcher := common.NewFileWatcher()
	if watcher == nil {
		log.Fatal().Str("component", "nods").Msg("file watcher")
	}
	cache := common.NewFastCache()
	client := resty.New()
	schema := server.NewDataSchema(watcher, cache, svrConfig.ConfigPath, svrConfig.SchemaPath)
	secrets := LoadSecrets()

	tools := server.Toolbox{
		Client:   client,
		Cache:    cache,
		Watcher:  watcher,
		Schema:   schema,
		Secrets:  secrets,
		DataPath: svrConfig.DataPath,
	}

	router := server.NewDataRouter(&tools)
	router.Init()

	rules := server.NewRuleProvider(router, schema)
	handlers := server.NewNormalizedDataServer(schema, router, secrets, rules)
	mux := chi.NewMux()
	mux.Route(svrConfig.PathPrefix, func(r chi.Router) {
		r.Use(middleware.RequestID)
		// view of the data dictionary
		r.Get("/v1/data/schema", handlers.DictionaryHandler)
		// query one item
		r.Get("/v1/data/{category}/{vendor}/{item}", handlers.GetItemHandler)
		r.Post("/v1/data/{category}/{vendor}/{item}", handlers.PostItemHandler)
		// query a item across a category
		r.Get("/v1/data/category/{category}", handlers.GetCategoryHandler)
		r.Post("/v1/data/category/{category}", handlers.PostCategoryHandler)
		// test a custom rule
		r.Get("/v1/data/rule", handlers.GetEvaluateHandler)
		r.Post("/v1/data/rule", handlers.PostEvaluateHandler)
		// whoami aggregate information
		r.Get("/v1/data/whoami", handlers.GetWhoamiHandler)
		r.Post("/v1/data/whoami", handlers.PostWhoamiHandler)
	})

	watcher.Listen()

	err := common.ListenAndServe(svrConfig.ListenAddr, "", "", mux)
	if err != nil {
		log.Fatal().Err(err).Str("component", "nods").Str("state", "stopped").Msg("orchestration")
	} else {
		log.Info().Str("component", "nods").Str("state", "stopped").Msg("orchestration")
	}
}

func LoadSecrets() *common.Secrets {
	keys := []string{
		"PWNED_API_KEY",
		"IPINFO_API_KEY",
		"FINGERPRINT_JWT_SECRET"}
	return common.NewSecrets(keys)
}
