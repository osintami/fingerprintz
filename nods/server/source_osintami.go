package server

import (
	"context"
	"encoding/json"

	"github.com/osintami/fingerprintz/common"
)

type InternalSource struct {
	browser IDataSource
	phone   IDataSource
	email   IDataSource
	rules   IDataSource
}

func NewInternalSource(router IDataRouter, tools *Toolbox) IDataSource {
	return NewDataInstance(tools, SOURCE_OSINTAMI_NAME, &InternalSource{
		browser: NewBrowserSource(tools),
		phone:   NewPhoneSource(tools),
		email:   NewEmailSource(tools),
		rules:   NewRuleSource(tools, router, tools.Schema)})
}

func (x *InternalSource) IsCached() bool {
	return false
}

func (x *InternalSource) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	switch categoryName {
	case CATEGORY_BROWSER:
		// NOTE:  a missing or corrupt user-agent regex file can cause this
		if x.browser != nil {
			return x.browser.CategoryInfo(ctx, categoryName, inputs)
		}
	case CATEGORY_EMAIL:
		return x.email.CategoryInfo(ctx, categoryName, inputs)
	// case CATEGORY_DOMAIN:
	// 	return x.email.CategoryInfo(ctx, categoryName, inputs)
	case CATEGORY_PHONE:
		return x.phone.CategoryInfo(ctx, categoryName, inputs)
	case CATEGORY_RULE:
		return x.rules.CategoryInfo(ctx, categoryName, inputs)
	}

	return nil, ErrNotImplemented
}
