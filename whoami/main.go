// Copyright © 2023 OSINTAMI. This is not yours.
package main

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/osintami/fingerprintz/whoami/server"
)

func main() {

	svrConfig := &server.ServerConfig{}
	common.LoadEnv(true, true, svrConfig)
	log.InitLogger(svrConfig.LogPath, "whoami.log", svrConfig.LogLevel, false)
	if svrConfig.LogLevel == "TRACE" {
		common.PrintEnvironment()
	}

	handlers := server.NewWhoamiServer(
		server.NewJWTSigner(server.NewJWT(os.Getenv("FINGEPRINT_JWT_SECRET"))),
		common.NewOSINTAMIClient(resty.New(), svrConfig.NodsURL, ""))

	router := chi.NewMux()
	router.Route(svrConfig.PathPrefix, func(r chi.Router) {
		r.Get("/v1/fingerprint/scan", handlers.GetFingerprintHandler)
		r.Post("/v1/fingerprint/scan", handlers.PostFingerprintHandler)
		r.Get("/v1/fingerprint/risk", handlers.GetRiskHandler)
		r.Post("/v1/fingerprint/risk", handlers.PostRiskHandler)
	})

	err := common.ListenAndServe(svrConfig.ListenAddr, "", "", router)
	if err != nil {
		log.Fatal().Err(err).Str("component", "whoami").Str("state", "stopped").Msg("orchestration")
	} else {
		log.Info().Str("component", "whoami").Str("state", "stopped").Msg("orchestration")
	}
}
