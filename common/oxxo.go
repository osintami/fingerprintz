// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/osintami/fingerprintz/log"
	"github.com/rs/zerolog"
)

func InitZeroLog(level string) error {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.Error().Err(err).Msg("zero log init")
		return err
	}
	zerolog.SetGlobalLevel(parsedLevel)
	return nil
}

func InitFileLog(path, fileName, logLevel string, stdError bool) {
	log.InitLogger(path, fileName, logLevel, stdError)
}

func LoadEnv(clearEnv bool, useEnv bool, output interface{}) error {
	if clearEnv {
		os.Clearenv()
	}
	if useEnv && godotenv.Load(".env.personal") != nil {
		if godotenv.Load(".env") != nil {
			fmt.Println("[fingerprintz] .env.personal and .env are both missing, using default values for config")
		}
	}
	return env.Parse(output)
}

func LoadJson(fileName string, cfg interface{}) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Error().Err(err).Str("component", "load json").Str("file", fileName).Msg("read file")
		return err
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		log.Error().Err(err).Str("component", "load json").Str("file", fileName).Msg("unmarshal")
		return err
	}
	return nil
}

func ListenAndServe(ListenAddr, SSLCertFile, SSLKeyFile string, router http.Handler) error {
	server := &http.Server{Addr: ListenAddr, Handler: router}
	var err error
	if SSLCertFile != "" && SSLKeyFile != "" {
		err = server.ListenAndServeTLS(SSLCertFile, SSLKeyFile)
	} else {
		err = server.ListenAndServe()
	}
	return err
}

func PrintEnvironment() {
	env := os.Environ()
	for _, variable := range env {
		pair := strings.Split(variable, "=")
		value := pair[1]
		if strings.Contains(pair[0], "API_KEY") || strings.Contains(pair[0], "PASS") || strings.Contains(pair[0], "USER") {
			value = "<masked>"
		}
		log.Info().Str(pair[0], value).Msg("environment")
	}
}

func PathParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func QueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func SendError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
}

func SendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func SendPrettyJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	// NOTE:  this makes it very easy to read and impossible to test
	encoder.SetIndent("", "    ")
	encoder.Encode(data)
}
func IpAddr(r *http.Request) string {
	addr := r.RemoteAddr
	if strings.HasPrefix(addr, "[::1]") {
		addr = "127.0.0.1:0"
	}
	parts := strings.Split(addr, ":")
	return parts[0]
}
