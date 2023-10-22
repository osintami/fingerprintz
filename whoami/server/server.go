// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"github.com/osintami/fingerprintz/common"
)

type ServerConfig struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
	LogPath     string `env:"LOG_PATH" envDefault:"/home/osintami/logs"`
	PathPrefix  string `env:"PATH_PREFIX" envDefault:"/"`
	ListenAddr  string `env:"LISTEN_ADDR" envDefault:"127.0.0.1:8083"`
	NodsURL     string `env:"OSINTAMI,required"`
	RedirectURL string `env:"LATENCY_REDIRECT,required"`
}

type WhoamiServer struct {
	nods            common.INods
	latencyRedirect string
	signer          IJWTSigner
	params          common.IParameterHelper
}

func NewWhoamiServer(signer IJWTSigner, nods common.INods, latencyRedirect string) *WhoamiServer {
	return &WhoamiServer{
		signer:          signer,
		nods:            nods,
		latencyRedirect: latencyRedirect,
		params:          common.NewParameterHelper(),
	}
}
