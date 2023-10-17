// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"net"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/osintami/fingerprintz/common"
)

type IIpInfo interface {
	GetIPInfo(net.IP) (*ipinfo.Core, error)
}

type IpInfoSource struct {
	client IIpInfo
}

func NewIpInfoSource(tools *Toolbox, ipinfo IIpInfo) IDataSource {
	return NewDataInstance(
		tools,
		SOURCE_IPINFO_NAME,
		&IpInfoSource{client: ipinfo})
}

func (x *IpInfoSource) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	switch categoryName {
	case CATEGORY_IPADDR:
		return x.GetIpInfo(ctx, inputs[categoryName])
	}
	return nil, ErrNotImplemented
}

func (x *IpInfoSource) IsCached() bool {
	return true
}

func (x *IpInfoSource) GetIpInfo(ctx context.Context, ip string) (json.RawMessage, error) {
	info, err := x.client.GetIPInfo(net.ParseIP(ip).To4())
	if err != nil {
		data, _ := json.Marshal(ipinfo.Core{})
		return data, common.ErrNoDataPresent
	}
	return json.Marshal(info)
}
