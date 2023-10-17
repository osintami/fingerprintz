// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/whois"
)

type WhoisProvider struct {
	whois whois.IWhois
}

type WhoisInfo struct {
	Domain             string
	DomainAgeInDays    int
	DomainAgeInYears   int
	DomainAgeDate      string
	IsRegisteredDomain bool
}

func NewWhoisSource(tools *Toolbox, whois whois.IWhois) IDataSource {
	return NewDataInstance(tools, SOURCE_WHOIS_NAME, &WhoisProvider{whois: whois})
}

func (x *WhoisProvider) IsCached() bool {
	return true
}

func (x *WhoisProvider) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	switch categoryName {
	case CATEGORY_DOMAIN:
		return x.getDomainInfo(ctx, inputs[categoryName])
	}
	return nil, ErrNotImplemented
}

func (x *WhoisProvider) getDomainInfo(ctx context.Context, domain string) (json.RawMessage, error) {

	result, err := x.whois.Get(domain, "1s")
	if err != nil {
		out, _ := json.Marshal(&whois.WhoisInfo{})
		return out, common.ErrNoDataPresent
	}

	info, err := x.whois.ParseInfo(domain, result)
	if err != nil {
		out, _ := json.Marshal(info)
		return out, common.ErrNoDataPresent
	}

	return json.Marshal(info)
}
