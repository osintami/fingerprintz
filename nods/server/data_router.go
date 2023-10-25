// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/osintami/fingerprintz/pwned"
	"github.com/osintami/fingerprintz/whois"
)

type IDataRouter interface {
	Init()
	DataValue(ctx context.Context, uri *DataURI, inputs common.DataInputs) (*common.DataOutput, error)
	CategoryValues(ctx context.Context, categoryName string, inputs common.DataInputs) ([]*common.DataOutput, error)
}

type DataRouter struct {
	tools     *Toolbox
	schema    IDataSchema
	secrets   common.ISecrets
	instances map[string]IDataSource
	response  *DataResponse
}

func NewDataRouter(tools *Toolbox) *DataRouter {

	x := &DataRouter{
		tools:     tools,
		secrets:   tools.Secrets,
		schema:    tools.Schema,
		instances: make(map[string]IDataSource),
		response:  NewDataResponse(),
	}
	return x
}

func (x *DataRouter) Init() {
	for _, vendor := range x.schema.ListSources() {
		if !vendor.Enabled {
			continue
		}

		switch vendor.Database {
		case "byod":
			log.Info().Str("component", "router").Str("source", vendor.Name).Msg("instantiating third party provider")
			provider := NewAPIProvider(x.tools, vendor.API)
			x.instances[vendor.Name] = NewDataInstance(x.tools, vendor.Name, provider)
		case "mmdb":
			log.Info().Str("component", "router").Str("source", vendor.Name).Msg("instantiating MMDB provider")
			reader, err := common.NewMaxmindReader(x.tools.DataPath + vendor.Name + ".mmdb")
			if err != nil {
				log.Error().Err(err).Str("component", "router").Str("source", vendor.Name).Msg("MMDB file missing")
			}
			provider := NewMMDBProvider(x.tools.Watcher, reader)
			x.instances[vendor.Name] = NewDataInstance(x.tools, vendor.Name, provider)
		case "fast":
			log.Info().Str("component", "router").Str("source", vendor.Name).Msg("instantiating FastDB provider")
			provider, _ := NewFastDBProvider(x.tools, vendor.Name)
			x.instances[vendor.Name] = NewDataInstance(x.tools, vendor.Name, provider)
		case "code":
			x.instances[vendor.Name] = x.newInstance(vendor.Name)
		}
	}
}

func (x *DataRouter) findSourceInstance(sourceName string) IDataSource {
	return x.instances[sourceName]
}

func (x *DataRouter) DataValue(ctx context.Context, dataURI *DataURI, inputs common.DataInputs) (*common.DataOutput, error) {
	// find the category
	if !x.schema.IsValidCategory(dataURI.CategoryName) {
		log.Error().Err(ErrSourceNotFound).Str("component", "router").Str("source", dataURI.SourceName).Msg("category not found")
		return x.response.EmptyResponse(common.Null, dataURI.ItemName, inputs, ErrCategoryNotFound), ErrCategoryNotFound
	}

	// make sure we have the correct input parameters for this category
	if dataURI.CategoryName != CATEGORY_RULE && inputs[dataURI.CategoryName] == "" {
		return x.response.EmptyResponse(common.Null, dataURI.Key(), inputs, ErrMissingInputs), common.ErrNoDataPresent
	}

	// find the source instance
	source := x.findSourceInstance(dataURI.SourceName)
	if source == nil {
		log.Error().Err(ErrSourceNotFound).Str("component", "router").Str("source", dataURI.SourceName).Msg("vendor not found")
		return x.response.EmptyResponse(common.Null, dataURI.ItemName, inputs, ErrSourceNotFound), ErrSourceNotFound
	}

	item, err := x.schema.Item(dataURI)
	if err != nil {
		log.Error().Err(err).Str("component", "router").Str("source", dataURI.SourceName).Str("data", dataURI.Key()).Str("inputs", inputs.String()).Str("value", "default").Msg("get data")
		return x.response.EmptyResponse(common.Null, dataURI.Key(), inputs, ErrItemNotFound), ErrItemNotFound
	}

	// query for the data value
	return source.ItemValue(ctx, item, inputs)
}

func (x *DataRouter) CategoryValues(ctx context.Context, categoryName string, inputs common.DataInputs) ([]*common.DataOutput, error) {

	itemResults := []*common.DataOutput{}

	items := x.schema.ListItemsByCategory(categoryName)
	for _, item := range items {
		if source := x.findSourceInstance(item.SourceName); source != nil {
			// strip out internal use only keys from response
			delete(inputs, common.INPUT_ROLE)
			delete(inputs, common.INPUT_KEY)
			delete(inputs, common.INPUT_TYPE)
			delete(inputs, common.INPUT_RULE)
			delete(inputs, common.INPUT_CSV)
			resp, _ := source.ItemValue(ctx, item, inputs)
			itemResults = append(itemResults, resp)
		}
	}

	if len(itemResults) == 0 {
		return itemResults, ErrItemNotFound
	}

	return itemResults, nil
}

const (
	// remote data sources
	SOURCE_PWNED_NAME    = "pwned"
	SOURCE_WHOIS_NAME    = "whois"
	SOURCE_IPINFO_NAME   = "ipinfo"
	SOURCE_SPAMHAUS_NAME = "spamhaus"
	// coded data sources
	SOURCE_OSINTAMI_NAME = "osintami"
)

func (x *DataRouter) newInstance(sourceName string) IDataSource {
	log.Info().Str("component", "router").Str("source", sourceName).Msg("instantiating CODE source")

	switch {
	case sourceName == SOURCE_PWNED_NAME:
		client := resty.New().GetClient()
		// NOTE:  basic tier subscription at haveibeenpwned is 6 requests per second
		// NOTE:  this lives here because it overrides httpmock
		client.Transport = common.NewThrottledTransport(7*time.Second, 1, http.DefaultTransport)
		return NewPwnedSource(x.tools, pwned.NewPwned(client, x.secrets.Find("PWNED_API_KEY")))
	case sourceName == SOURCE_WHOIS_NAME:
		return NewWhoisSource(x.tools, whois.NewWhois())
	case sourceName == SOURCE_IPINFO_NAME:
		client := resty.New().GetClient()
		// NOTE:  free tier subscription at ipinfo.o is 50,000 per month or about 1 per second
		// NOTE:  this lives here because it overrides httpmock
		client.Transport = common.NewThrottledTransport(1*time.Minute, 10, http.DefaultTransport)
		ipinfo := ipinfo.NewClient(
			client,
			nil,
			x.secrets.Find("IPINFO_API_KEY"))
		return NewIpInfoSource(x.tools, ipinfo)
	case sourceName == SOURCE_SPAMHAUS_NAME:
		return NewSpamhausSource(x.tools, common.NewRealtimeBlackholeList())
	case sourceName == SOURCE_OSINTAMI_NAME:
		return NewInternalSource(x, x.tools)
	}

	log.Error().Err(ErrNotImplemented).Str("component", "router").Str("source", sourceName).Msg("source no esta aqui")

	return nil
}
