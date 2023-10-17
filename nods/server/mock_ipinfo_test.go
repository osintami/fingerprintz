// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"encoding/json"
	"net"

	"github.com/ipinfo/go/v2/ipinfo"
)

// -----------------------------------------------------------------------------
type MockIpInfo struct {
	fail bool
}

func NewMockIpInfo(fail bool) IIpInfo {
	return &MockIpInfo{fail: fail}
}

func (x *MockIpInfo) GetIPInfo(addr net.IP) (*ipinfo.Core, error) {
	if x.fail {
		return nil, ErrBadData
	}
	data := "{\"ip\":\"1.2.3.4\",\"hostname\":\"unn-89-187-175-145.cdn77.com\",\"city\":\"Dallas\",\"region\":\"Texas\",\"country\":\"US\",\"country_name\":\"United States\",\"country_flag\":{\"emoji\":\"🇺🇸\",\"unicode\":\"U+1F1FA U+1F1F8\"},\"country_flag_url\":\"https://cdn.ipinfo.io/static/images/countries-flags/US.svg\",\"country_currency\":{\"code\":\"USD\",\"symbol\":\"$\"},\"continent\":{\"code\":\"NA\",\"name\":\"North America\"},\"loc\":\"32.7813,-96.8019\",\"org\":\"AS60068 Datacamp Limited\",\"postal\":\"75270\",\"timezone\":\"America/Chicago\"}"
	v := &ipinfo.Core{}
	json.Unmarshal([]byte(data), v)
	return v, nil
}
