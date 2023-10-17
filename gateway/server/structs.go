// Copyright © 2023 OSINTAMI. This is not yours.
package server

type ServerConfig struct {
	PathPrefix  string `env:"PATH_PREFIX" envDefault:"/"`
	ListenAddr  string `env:"LISTEN_ADDR" envDefault:"127.0.0.1:8080"`
	PgDB        string `env:"POSTGRES_DB" envDefault:"osintami"`
	PgHost      string `env:"POSTGRES_HOST" envDefault:"127.0.0.1"`
	PgHostAuth  string `env:"POSTGRES_HOST_AUTH_METHOD" envDefault:"trust"`
	PgPassword  string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	PgPort      string `env:"POSTGRES_PORT" envDefault:"5432"`
	PgUser      string `env:"POSTGRES_USER" envDefault:"postgres"`
	SSLCertFile string `env:"CERT_FILE"`
	SSLKeyFile  string `env:"KEY_FILE"`
	LogPath     string `env:"LOG_PATH" envDefault:"/home/osintami/logs"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
	Osintami    string `env:"OSINTAMI" envDefault:"http://127.0.0.1:8082/nods/v1"`
}
