// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/etlr/etl"
)

type ServerConfig struct {
	DataPath   string `env:"DATA_PATH" envDefault:"/home/osintami/data/"`
	ConfigPath string `env:"CONFIG_PATH" envDefaults:"./"`
	LogPath    string `env:"LOGS_PATH" envDefault:"/home/osintami/logs"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"INFO"`
	PathPrefix string `env:"PATH_PREFIX" envDefault:"/etlr"`
	ListenAddr string `env:"LISTEN_ADDR" envDefault:"127.0.0.1:8081"`
}

type Message struct {
	Message string `json:"message"`
}

type ETLrServer struct {
	manager etl.IETLManager
}

func NewETLrServer(manager etl.IETLManager) *ETLrServer {
	return &ETLrServer{manager: manager}
}

func (x *ETLrServer) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	vendorName := common.PathParam(r, "vendor")
	// NOTE:  this is dangerous, as you might really tick off
	//   a data source/vendor by being too frequent a flyer
	if vendorName == "ALL" {
		x.manager.RefreshAll()
		common.SendJSON(w, &Message{Message: "success"})
		return
	}

	err := x.manager.Refresh(vendorName)
	if err != nil {
		common.SendError(w, err, http.StatusInternalServerError)
		return
	}

	common.SendJSON(w, &Message{Message: "success"})
}
