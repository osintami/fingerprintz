// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"

	"github.com/osintami/fingerprintz/common"
)

type SpamhausInfo struct {
	IPAddr        string
	Domain        string
	IsBlacklisted bool
}

type SpamhausSource struct {
	rbl common.IRealtimeBlackholeList
}

func NewSpamhausSource(tools *Toolbox, rbl common.IRealtimeBlackholeList) IDataSource {
	return NewDataInstance(
		tools,
		SOURCE_SPAMHAUS_NAME,
		&SpamhausSource{rbl: rbl})
}

func (x *SpamhausSource) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	info := &SpamhausInfo{}
	switch categoryName {
	case CATEGORY_IPADDR:
		info.IPAddr = inputs[categoryName]
		info.IsBlacklisted = x.rbl.Lookup("fingerprintz.spamhaus.org", inputs[categoryName])
		return json.Marshal(info)
	case CATEGORY_DOMAIN:
		info.Domain = inputs[categoryName]
		info.IsBlacklisted = x.rbl.Lookup("dbl.spamhaus.org", inputs[categoryName])
		return json.Marshal(info)
	}
	return nil, ErrNotImplemented
}

func (x *SpamhausSource) IsCached() bool {
	return true
}
