// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"strings"
)

type MockGatewayServices struct {
	services []*GatewayServiceInfo
}

func NewMockGatewayServices() IGatewayServices {
	services := make([]*GatewayServiceInfo, 3)

	service := &GatewayServiceInfo{}
	service.ACL = "user"
	service.Name = "nods"
	service.Route = "/data/"
	service.Service = "http://127.0.0.1:8082"
	service.Path = "/nods/v1/data/"
	services[0] = service

	service = &GatewayServiceInfo{}
	service.ACL = "admin"
	service.Name = "admin"
	service.Route = "/admin/"
	service.Service = "http://127.0.0.1:8082"
	service.Path = "/admin/view"
	services[1] = service

	service = &GatewayServiceInfo{}
	service.ACL = "user"
	service.Name = "corrupt"
	service.Route = "/corrupt/"
	service.Service = "!@#$%^&*()"
	service.Path = "nope"
	services[2] = service
	return &MockGatewayServices{services: services}
}

func (x *MockGatewayServices) find(path string) *GatewayServiceInfo {
	for _, service := range x.services {
		if strings.HasPrefix(path, service.Route) {
			return service
		}
	}
	return nil
}
