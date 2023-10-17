// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"github.com/osintami/fingerprintz/common"
)

type ServerConfig struct {
	ConfigPath string `env:"LOCAL_CONFIG_PATH" envDefault:"/home/osintami/nods/"`
	SchemaPath string `env:"LOCAL_SCHEMA_PATH" envDefault:"/home/osintami/data/"`
	DataPath   string `env:"LOCAL_DB_PATH" envDefault:"/home/osintami/data/"`
	LogPath    string `env:"LOG_PATH" envDefault:"/home/osintami/logs/"`
	PathPrefix string `env:"PATH_PREFIX" envDefault:"/"`
	ListenAddr string `env:"LISTEN_ADDR,required" envDefault:"127.0.0.1:8082"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"INFO"`
}

type NormalizedDataServer struct {
	router  IDataRouter
	schema  IDataSchema
	secrets common.ISecrets
	params  common.IParameterHelper
	rules   IDataProvider
}

func NewNormalizedDataServer(schema IDataSchema, router IDataRouter, secrets common.ISecrets, rules IDataProvider) *NormalizedDataServer {
	return &NormalizedDataServer{
		router:  router,
		schema:  schema,
		secrets: secrets,
		params:  common.NewParameterHelper(),
		rules:   rules}
}
