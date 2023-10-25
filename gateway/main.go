// Copyright © 2023 OSINTAMI. This is not yours.
package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-resty/resty/v2"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/gateway/server"
	"github.com/osintami/fingerprintz/log"

	"github.com/go-chi/httprate"
)

func main() {
	svrConfig := &server.ServerConfig{}
	common.LoadEnv(true, true, svrConfig)
	log.InitLogger(svrConfig.LogPath, "gateway.log", svrConfig.LogLevel, false)
	if svrConfig.LogLevel == "TRACE" {
		common.PrintEnvironment()
	}

	pgConfig := &common.PostgresConfig{
		PgHost:     svrConfig.PgHost,
		PgPort:     svrConfig.PgPort,
		PgUser:     svrConfig.PgUser,
		PgPassword: svrConfig.PgPassword,
		PgDB:       svrConfig.PgDB}

	gorm, err := common.OpenDB(pgConfig, svrConfig.LogPath)
	if err != nil {
		log.Fatal().Err(err).Msg("postgres not running")
	}

	// migrations
	gorm.AutoMigrate(&server.Account{})
	gorm.AutoMigrate(&server.Call{})
	gorm.AutoMigrate(&server.Pixel{})

	// first time setup
	accounts := server.NewAccounts(gorm, server.NewSender(), "welcome.template")
	user, err := accounts.FindByEmail(context.Background(), "admin@osintami.com")
	if user == nil || err != nil {
		// create admin user, see postgres osintami/accounts table for API key
		accounts.CreateAccount(context.Background(), &server.Account{
			Name:        "OSINTAMI Admin",
			Email:       "admin@osintami.com",
			ApiKey:      "",
			Role:        "admin",
			Tokens:      0,
			LastPayment: time.Now(),
			Enabled:     true})
	}

	cache := common.NewPersistentCache("unwanted.db")
	cache.LoadFile("unwanted.db")

	shutdown := common.NewShutdownHandler()
	shutdown.AddListener(cache.Persist)
	shutdown.Listen(func() { os.Exit(1) })

	oc := common.NewOSINTAMIClient(resty.New(), svrConfig.Nods, svrConfig.Whoami)

	services := server.NewGatewayServices("config.json")
	handlers := server.NewGatewayServer(
		server.NewReverseProxy(),
		oc,
		cache,
		accounts,
		server.NewCalls(gorm),
		server.NewPixels(gorm),
		services)

	router := chi.NewMux()
	router.Group(func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(middleware.RealIP)
		// NOTE:  whitelisted IPs are off pending licensing
		// r.Use(handlers.CheckWhitelist)
		r.Use(httprate.Limit(
			10,
			1*time.Minute,
			httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
				key := r.Header.Get("X-Api-Key")
				if key == "" {
					key = r.URL.Query().Get("key")
				}
				return key, nil
			}),
			httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "{\"error\":\"rate limit exceeded.\"}", http.StatusTooManyRequests)
			}),
		))
		r.Get("/", handlers.UnwantedHandler)
		r.Get("/*", handlers.ReverseProxyHandler)
		r.Post("/*", handlers.ReverseProxyHandler)
	})
	router.Group(func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(middleware.RealIP)
		r.Use(httprate.Limit(
			1,
			1*time.Minute,
			httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
				return common.IpAddr(r), nil
			}),
			httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "{\"error\":\"please hold.\"}", http.StatusTooManyRequests)
			}),
		))
		r.Post("/v1/stripe", handlers.StripeHandler)
		r.Get("/signup", handlers.SignupHandler)
	})

	router.Group(func(r chi.Router) {
		r.Get("/images/{img}", handlers.PixelFireHandler)
	})

	log.Debug().Str("component", "main").Msg("listen and serve")

	err = common.ListenAndServe(svrConfig.ListenAddr, svrConfig.SSLCertFile, svrConfig.SSLKeyFile, router)
	if err != nil {
		log.Fatal().Err(err).Str("component", "gateway").Str("state", "stopped").Msg("orchestration")
	} else {
		log.Info().Str("component", "gateway").Str("state", "stopped").Msg("orchestration")
	}
}
