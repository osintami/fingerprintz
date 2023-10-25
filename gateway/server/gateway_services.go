// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"strings"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type IGatewayServices interface {
	find(path string) *GatewayServiceInfo
}

type GatewayServiceInfo struct {
	Name    string
	ACL     string
	Route   string
	Service string
	Path    string
}

func (x GatewayServiceInfo) IsPublic() bool {
	return x.ACL == ""
}

func (x GatewayServiceInfo) IsPrivate() bool {
	return !x.IsPublic()
}

func (x GatewayServiceInfo) WantsUser() bool {
	return x.ACL == "user"
}

func (x GatewayServiceInfo) WantsAdmin() bool {
	return x.ACL == "admin"
}

type GatewayServices struct {
	services []*GatewayServiceInfo
}

func NewGatewayServices(fileName string) *GatewayServices {
	services := []*GatewayServiceInfo{}
	err := common.LoadJson(fileName, &services)
	if err != nil {
		log.Error().Err(err).Str("component", "gateway services").Str("file", fileName).Msg("gateway services config missing")
		return nil
	}
	return &GatewayServices{services: services}
}

func (x *GatewayServices) find(path string) *GatewayServiceInfo {
	for _, service := range x.services {
		if strings.HasPrefix(path, service.Route) {
			return service
		}
	}
	return nil
}
