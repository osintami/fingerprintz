// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"time"

	"github.com/osintami/fingerprintz/common"
)

type MockNods struct {
	fail bool
}

func NewMockNods(fail bool) common.IOsintami {
	return &MockNods{fail: fail}
}

func (x *MockNods) Whoami(ip string) (*common.WhoamiInfo, error) {

	info := &common.WhoamiInfo{}

	if x.fail {
		return info, common.ErrItemNotFound
	}

	if ip == "1.2.3.4" {
		info.Blacklist = true
		info.Bot = true
		info.CloudNode = true
		info.Proxy = true
		info.Tor = true
		info.VPN = true
		info.Latitude = 37.511200
		info.Longitude = 126.974100
		info.City = "Seoul"
		info.Country = "South Korea"
	}
	if ip == "4.3.2.1" {
		info.Latitude = 30.633263
		info.Longitude = -97.677986
		info.City = "Georgetown"
		info.Country = "United States"
	}
	info.LastSeen = time.Now().Format(common.GO_DEFAULT_DATE)
	return info, nil
}

func (x *MockNods) Item(name string, keys map[string]string) (*common.DataOutput, error) {
	return nil, nil
}

func (x *MockNods) Fingerprint(keys map[string]string) (string, error) {
	return "", nil
}
