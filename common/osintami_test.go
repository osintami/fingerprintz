// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestOsintamiItem(t *testing.T) {
	hc := resty.New()
	httpmock.ActivateNonDefault(hc.GetClient())

	oc := NewOSINTAMIClient(hc, "https://api.osintami.com", "")

	// login OSINT data API
	content := "{\"Item\": \"rule/osintami/isBlacklisted\",\"Result\": {\"Type\": 1,\"Bool\": true},\"Keys\": {\"ip\": \"1.2.3.4\"},\"Error\": \"\"}"
	url := "https://api.osintami.com/data/rule/osintami/isBlacklisted?ip=1.2.3.4"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	keys := make(map[string]string)
	keys["ip"] = "1.2.3.4"
	out, err := oc.Item("rule/osintami/isBlacklisted", keys)

	assert.Nil(t, err)
	assert.True(t, out.Result.Boolean())

	// test HTTP error
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(http.StatusInternalServerError, content))
	defer httpmock.DeactivateAndReset()

	out, err = oc.Item("rule/osintami/isBlacklisted", keys)
	assert.NotNil(t, err)
	assert.False(t, out.Result.Boolean())
}

func TestOsintamiWhoami(t *testing.T) {
	hc := resty.New()
	httpmock.ActivateNonDefault(hc.GetClient())

	oc := NewOSINTAMIClient(hc, "https://api.osintami.com", "")

	// whoami OSINT data API
	content := "{\"IpAddr\": \"1.2.3.4\",\"Tor\": true,\"CloudNode\": false,\"Proxy\": false,\"VPN\": false,\"Bot\": false,\"Blacklist\": true,\"Latitude\": 0,\"Longitude\": 0,\"City\": \"\",\"Country\": \"\",\"Continent\": \"\",\"LastSeen\": \"2023-20-12 12:33:13\"}"
	url := "https://api.osintami.com/v1/data/whoami?ip=1.2.3.4"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(http.StatusOK, content))
	defer httpmock.DeactivateAndReset()

	out, err := oc.Whoami("1.2.3.4")
	assert.Nil(t, err)
	assert.True(t, out.Tor)
	assert.True(t, out.Blacklist)

	// test HTTP error
	url = "https://api.osintami.com/v1/data/whoami?ip=4.3.2.1"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(http.StatusInternalServerError, content))

	out, err = oc.Whoami("4.3.2.1")
	assert.NotNil(t, err)
	assert.False(t, out.Tor)
	assert.False(t, out.Blacklist)

	// test invalid JSON returned
	url = "https://api.osintami.com/v1/data/whoami?ip=8.8.8.8"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(http.StatusOK, "."))

	out, err = oc.Whoami("8.8.8.8")
	assert.NotNil(t, err)
	assert.False(t, out.Tor)
	assert.False(t, out.Blacklist)
}

func TestOsintamiFingerprint(t *testing.T) {
	hc := resty.New()
	httpmock.ActivateNonDefault(hc.GetClient())

	oc := NewOSINTAMIClient(hc, "", "http://localhost:8083/whoami")

	// fingerprint OSINT data API
	content := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0aW1lIjoiMjAyMy0xMC0wM1QyMjozNzo1Ni43NTc4NTgyMi0wNjowMCIsImVoYXNoIjoiNjcxYThjOTQwNGFiNjM0Mjk2MTM2NWFhMDBjODIwYTQ0YzUwNGNhNWYxOWNhNGU4MTZmY2NmZmZjYjI2Y2I5YSIsImxhdGl0dWRlIjo0MS4yNTkxLCJsb25naXR1ZGUiOi05NS44NTE3LCJjaXR5IjoiQ291bmNpbCBCbHVmZnMiLCJjb3VudHJ5IjoiVVMiLCJpcCI6IjM0LjMxLjE3MS4yMSIsInVhIjoiZTJhNDkzNTRhNDQ0YzA1YmFlYzlmYjU0NjkwYTVjZWUyNmY4ZjM4YzI5MzIyYTgwNDE4YjgwMmViZGEwY2YyNSIsImh3aWQiOiIwMTIzNDU2Ny04OUFCQ0RFRi0wMTIzNDU2Ny04OUFCQ0RFRiIsIm5pZCI6IjBYMDAxMCIsInBpZCI6Im9zaW50YW1pIiwidmVyc2lvbiI6IjEuMC4wIiwiaXNzIjoib3NpbnRhbWkiLCJzdWIiOiJmaW5nZXJwcmludCIsImV4cCI6MTY5NjQ4MDY3NiwibmJmIjoxNjk2Mzk0Mjc2LCJpYXQiOjE2OTYzOTQyNzZ9.2g_qiH_Y-QtR4YmxfCZBFdEX3n8e40iTkMWkInHZG8s"
	url := "http://localhost:8083/whoami/v1/fingerprint/scan"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(http.StatusOK, content))
	defer httpmock.DeactivateAndReset()

	keys := make(map[string]string)
	keys["ip"] = "1.2.3.4"
	keys["email"] = "1@2.com"
	keys["ua"] = "test-user-agent"
	keys["device"] = "test-device-id"

	out, err := oc.Fingerprint(keys)
	assert.Nil(t, err)
	assert.Equal(t, content, out)

	// test HTTP error
	url = "http://localhost:8083/whoami/v1/fingerprint/scan"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(http.StatusInternalServerError, content))

	out, err = oc.Fingerprint(keys)
	assert.NotNil(t, err)
	assert.Equal(t, "", out)
}

func TestDataType(t *testing.T) {
	dataType := Null
	assert.Equal(t, dataType.String(), "Null")
	assert.Equal(t, Null, dataType.ToDataType("Null"))

	dataType = -1
	assert.Equal(t, dataType.String(), "")
	assert.Equal(t, Null, dataType.ToDataType("Nope"))

	dataType = Boolean
	assert.Equal(t, dataType.String(), "Boolean")
	assert.Equal(t, Boolean, dataType.ToDataType("Boolean"))

	dataType = Float
	assert.Equal(t, dataType.String(), "Float")
	assert.Equal(t, Float, dataType.ToDataType("Float"))

	dataType = Integer
	assert.Equal(t, dataType.String(), "Integer")
	assert.Equal(t, Integer, dataType.ToDataType("Integer"))

	dataType = String
	assert.Equal(t, dataType.String(), "String")
	assert.Equal(t, String, dataType.ToDataType("String"))

	dataType = Date
	assert.Equal(t, dataType.String(), "Date")
	assert.Equal(t, Date, dataType.ToDataType("Date"))

	dataType = JSON
	assert.Equal(t, dataType.String(), "JSON")
	assert.Equal(t, JSON, dataType.ToDataType("JSON"))
}

func TestDataResult(t *testing.T) {
	dataResult := &DataResult{Type: String, Raw: ""}
	assert.True(t, dataResult.IsEmpty())

	dataResult = &DataResult{Type: Boolean, Raw: "true"}
	assert.False(t, dataResult.IsEmpty())
	assert.True(t, dataResult.Boolean())
}

func TestDataInputs(t *testing.T) {
	dataInputs := DataInputs{}
	dataInputs["key"] = "value"
	// special case
	dataInputs["rule"] = "rule/osintami/isBlacklisted"
	params := dataInputs.Params()

	assert.Equal(t, 2, len(dataInputs))
	assert.Equal(t, "?key=value", params)
	assert.Equal(t, dataInputs.String(), params)

	dataInputs["ip"] = "1.2.3.4"
	params = dataInputs.Params()

	assert.Equal(t, "?ip=1.2.3.4&key=value", params)
	assert.Equal(t, dataInputs.String(), params)
}
