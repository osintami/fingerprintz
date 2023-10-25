// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGatewayServices(t *testing.T) {
	svc := NewGatewayServices("../config.json")
	assert.NotNil(t, svc)
	svcInfo := &GatewayServiceInfo{Name: "nods", ACL: "user", Route: "/data/", Service: "http://localhost:8082", Path: "/nods/v1/data/"}
	assert.True(t, svcInfo.IsPrivate())
	assert.False(t, svcInfo.IsPublic())
	assert.True(t, svcInfo.WantsUser())
	assert.False(t, svcInfo.WantsAdmin())

	svcInfo = svc.find("/data/")
	assert.Equal(t, "nods", svcInfo.Name)

	svcInfo = svc.find("nope")
	assert.Nil(t, svcInfo)
}

func TestGatewayServicesNoConfig(t *testing.T) {
	svc := NewGatewayServices("nope.json")
	assert.Nil(t, svc)
}
