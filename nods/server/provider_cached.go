// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/osintami/fingerprintz/common"
)

type CachedProvider struct {
	sourceName string
	provider   IDataProvider
	cache      common.IFastCache
}

const DEFAULT_CACHE_TTL = 24 * time.Hour

func NewCachedProvider(tools *Toolbox, sourceName string, provider IDataProvider) IDataProvider {
	return &CachedProvider{
		sourceName: sourceName,
		provider:   provider,
		cache:      tools.Cache,
	}
}

func (x *CachedProvider) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	key := categoryName + "/" + x.sourceName + "/" + inputs[categoryName]
	var err error
	out, found := x.cache.Get(key)
	if !found {
		out, err = x.provider.CategoryInfo(ctx, categoryName, inputs)
		if err == nil {
			x.cache.Set(key, out, DEFAULT_CACHE_TTL)
		}
	}
	return out.(json.RawMessage), err
}

func (x *CachedProvider) IsCached() bool {
	return true
}
