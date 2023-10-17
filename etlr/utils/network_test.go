// Copyright © 2023 OSINTAMI. This is not yours.
package utils

import (
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCIDR(t *testing.T) {
	network := buildNetwork()

	// range of IP's, calculate the network mask
	cidr, err := network.CalculateCIDR(net.ParseIP("1.2.3.4"), net.ParseIP("1.2.3.0"))
	assert.Nil(t, err)
	bits, _ := cidr.Mask.Size()
	assert.Equal(t, 29, bits)
	assert.Equal(t, "1.2.3.0", cidr.IP.String())

}

func TestParseCIDR(t *testing.T) {
	network := buildNetwork()

	// normal IPv4
	cidr, err := network.ParseCIDR("1.2.3.4/24")
	assert.Nil(t, err)
	assert.Equal(t, "1.2.3.0", cidr.IP.String())
	bits, _ := cidr.Mask.Size()
	assert.Equal(t, 24, bits)

	// private IPv4
	cidr, err = network.ParseCIDR("10.0.1.100/24")
	assert.Equal(t, ErrPrivateNetworkAddress, err)
	assert.Nil(t, cidr)

	// missing mask and private Ipv4
	cidr, err = network.ParseCIDR("10.0.1.100")
	assert.Equal(t, ErrPrivateNetworkAddress, err)
	assert.Nil(t, cidr)

	// missing mask normal IPv4
	cidr, err = network.ParseCIDR("1.2.3.4")
	assert.Nil(t, err)
	bits, _ = cidr.Mask.Size()
	assert.Equal(t, 32, bits)

	// invalid IPv4
	cidr, err = network.ParseCIDR("nope")
	assert.NotNil(t, err)
	assert.Nil(t, cidr)

	// Ipv6
	cidr, err = network.ParseCIDR("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	assert.Nil(t, err)
	bits, _ = cidr.Mask.Size()
	assert.Equal(t, 128, bits)
}

func TestIp4ToUint(t *testing.T) {
	network := buildNetwork()

	// 16 byte IPv4
	assert.Equal(t, uint(16909060), network.Ipv4ToUint(net.ParseIP("1.2.3.4")))

	// big number IPv4
	num := network.Ipv4ToUint(net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))
	assert.Equal(t, uint(57701172), num)

	// unspecified IPv4
	assert.Equal(t, uint(0x0), network.Ipv4ToUint(net.ParseIP("0.0.0.0")))

	// 4 byte IPv4
	assert.Equal(t, uint(16909060), network.Ipv4ToUint(net.ParseIP("1.2.3.4").To4()))
}

func TestUint2IPv4(t *testing.T) {
	network := buildNetwork()
	ip := network.Uint2IPv4(16909060)
	assert.Equal(t, "1.2.3.4", ip.To4().String())
}

func TestContent(t *testing.T) {
	network := buildNetwork()
	httpmock.ActivateNonDefault(network.client.GetClient())
	content := `"{"message":"the earth is flat"}`
	url := "http://localhost:8082/v1/data"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	data, code, err := network.Content(url)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, content, string(data))
}

func TestContentErrorFromRemoteHost(t *testing.T) {
	network := buildNetwork()
	httpmock.ActivateNonDefault(network.client.GetClient())
	content := `"{"message":"the earth is flat"}`
	url := "nope://localhost:8082/v1/nope"
	httpmock.RegisterResponder(
		"*?", url, httpmock.NewStringResponder(500, content))
	defer httpmock.DeactivateAndReset()

	data, code, err := network.Content(url)
	assert.NotNil(t, err)
	assert.Equal(t, -1, code)
	assert.Equal(t, "", string(data))
}

func TestDownloadFile(t *testing.T) {
	downloadedFile := "/tmp/test-download-file.csv"
	os.Remove(downloadedFile)

	network := buildNetwork()
	httpmock.ActivateNonDefault(network.client.GetClient())
	content := `test download to file`
	url := "http://localhost:8082/v1/download"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(200, content))
	defer httpmock.DeactivateAndReset()

	err := network.DownloadFile(url, downloadedFile)
	assert.Nil(t, err)

	_, err = os.Stat(downloadedFile)
	assert.Nil(t, err)

	data, err := os.ReadFile(downloadedFile)
	assert.Nil(t, err)
	assert.Equal(t, "test download to file", string(data))

	os.Remove(downloadedFile)
}

func TestDownloadFileErrorResponse(t *testing.T) {
	downloadedFile := "/tmp/test-download-file.csv"

	network := buildNetwork()
	httpmock.ActivateNonDefault(network.client.GetClient())
	content := `test download to file`
	url := "http://localhost:8082/v1/download"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(500, content))
	defer httpmock.DeactivateAndReset()

	err := network.DownloadFile(url, downloadedFile)
	assert.NotNil(t, err)

	os.Remove(downloadedFile)
}

var NETWORK *Network

func buildNetwork() *Network {
	if NETWORK == nil {
		NETWORK = NewNetworkingHelper(resty.New())
	}
	return NETWORK
}
