// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/tidwall/gjson"
)

const PASSTHROUGH = ""

type IDataSource interface {
	CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error)
	ItemValue(ctx context.Context, item Item, inputs common.DataInputs) (*common.DataOutput, error)
	IsCached() bool
}

type IDataProvider interface {
	CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error)
	IsCached() bool
}

type DataInstance struct {
	name     string
	tools    *Toolbox
	provider IDataProvider
	response *DataResponse
}

func NewDataInstance(tools *Toolbox, sourceName string, provider IDataProvider) IDataSource {

	x := &DataInstance{
		name:     sourceName,
		tools:    tools,
		provider: provider,
		response: NewDataResponse(),
	}
	if provider.IsCached() {
		x.provider = NewCachedProvider(tools, sourceName, provider)
	}
	return x
}

func (x *DataInstance) IsCached() bool {
	return x.provider.IsCached()
}

func (x *DataInstance) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	// short circuit on unsupported category
	if !x.tools.Schema.IsValidCategory(categoryName) {
		return []byte("{}"), ErrCategoryNotFound
	}

	// short circuit on missing data key
	if inputs[categoryName] == "" {
		return []byte("{}"), ErrMissingInputs
	}

	return x.provider.CategoryInfo(ctx, categoryName, inputs)
}

func (x *DataInstance) ItemValue(ctx context.Context, item Item, inputs common.DataInputs) (*common.DataOutput, error) {
	data, err := x.provider.CategoryInfo(ctx, item.CategoryName, inputs)
	if err != nil {
		log.Warn().Err(err).Str("component", "source").Str("source", x.name).Str("data", item.Path).Str("inputs", inputs.String()).Str("value", "default").Msg("category info")
		return x.response.EmptyResponse(item.Type, item.Path, inputs, err), common.ErrNoDataPresent
	}

	var out gjson.Result
	if item.Type == common.JSON && item.Gjson == PASSTHROUGH {
		// NOTE:  special case for rules
		out = gjson.Result{Type: gjson.JSON, Raw: string(data)}
	} else {
		out = gjson.GetBytes(data, item.Gjson)
	}
	if out.Exists() {
		return &common.DataOutput{
			Item:   item.Path,
			Keys:   inputs,
			Result: x.response.MarshalResult(out, item.Type),
		}, nil
	}
	return x.response.EmptyResponse(item.Type, item.Path, inputs, common.ErrNoDataPresent), common.ErrNoDataPresent
}
