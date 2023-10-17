// Copyright © 2023 OSINTAMI. This is not yours.
package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/osintami/fingerprintz/log"
)

type INetworking interface {
	CalculateCIDR(ip1, ip2 net.IP) (*net.IPNet, error)
	ParseCIDR(ip string) (*net.IPNet, error)
	DownloadFile(url string, fileName string) error
	Ipv4ToUint(ip net.IP) uint
	Uint2IPv4(ip uint64) net.IP
	Content(url string) ([]byte, int, error)
}

var ErrPrivateNetworkAddress = errors.New("private network address")

type Network struct {
	client *resty.Client
}

func NewNetworkingHelper(client *resty.Client) *Network {
	return &Network{client: client}
}

func (x *Network) Content(url string) ([]byte, int, error) {
	resp, err := x.client.R().Get(url)
	if err != nil {
		log.Error().Err(err).Str("component", "network").Str("url", url).Msg("http get")
		return nil, -1, err
	}
	return resp.Body(), resp.StatusCode(), nil
}

func (x *Network) CalculateCIDR(ip1, ip2 net.IP) (*net.IPNet, error) {
	subnet := ""
	maxLen := 32
	for l := maxLen; l >= 0; l-- {
		mask := net.CIDRMask(l, maxLen)
		na := ip1.Mask(mask)
		n := net.IPNet{IP: na, Mask: mask}

		if n.Contains(ip2) {
			subnet = fmt.Sprintf("%v/%v", na, l)
			break
		}
	}
	_, cidr, err := net.ParseCIDR(subnet)
	return cidr, err
}

func (x *Network) ParseCIDR(ipStr string) (*net.IPNet, error) {
	if strings.Contains(ipStr, "/") {
		ip, cidr, err := net.ParseCIDR(ipStr)
		if ip.IsPrivate() {
			return nil, ErrPrivateNetworkAddress
		}
		return cidr, err
	} else {
		mask := "/32"
		if strings.Contains(ipStr, ":") {
			mask = "/128"
		}
		ip, cidr, err := net.ParseCIDR(ipStr + mask)
		if ip.IsPrivate() {
			return nil, ErrPrivateNetworkAddress
		}
		return cidr, err
	}
}

func (x *Network) DownloadFile(url string, fileName string) error {
	resp, err := x.client.R().SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36").SetOutput(fileName).Get(url)
	if err != nil || resp.StatusCode() != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("HTTP download status %d %s", resp.StatusCode(), url)
		}
		log.Error().Err(err).Str("component", "network").Str("url", url).Msg("download file")
		return err
	}
	return nil
}

func (x *Network) Ipv4ToUint(ip net.IP) uint {
	if ip.IsUnspecified() {
		return 0
	}
	if len(ip) == 16 {
		return uint(binary.BigEndian.Uint32(ip[12:16]))
	}
	return uint(binary.BigEndian.Uint32(ip))
}

func (x *Network) Uint2IPv4(ip uint64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ip & 0xFF)
	bytes[1] = byte((ip >> 8) & 0xFF)
	bytes[2] = byte((ip >> 16) & 0xFF)
	bytes[3] = byte((ip >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).To4()
}
