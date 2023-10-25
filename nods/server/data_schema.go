// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"sort"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type SourceInfo struct {
	Name     string
	Database string // mmdb, yaml, fast, api
	Enabled  bool
	API      *API `json:",omitempty"`
}

type CategoryInfo struct {
	Name string
}

type HeaderInfo struct {
	Key   string
	Value string
}

type BasicAuthInfo struct {
	Name string
	Pass string
}

type API struct {
	Name    string
	URL     string
	Method  string
	ApiKey  *string        `json:",omitempty"`
	Auth    *BasicAuthInfo `json:",omitempty"`
	Headers []HeaderInfo   `json:",omitempty"`
}

type ServerInfo struct {
	Host string
}

type SchemaInfo struct {
	Sources []SourceInfo
}

type IDataSchema interface {
	Item(dataURI *DataURI) (Item, error)
	IsValidItem(dataURI *DataURI) bool
	IsValidCategory(categoryName string) bool
	IsEnabled(sourceName string) bool
	Source(sourceName string) (SourceInfo, error)
	ListItems() []Item
	ListSources() []SourceInfo
	ListCategories() []string
	ListItemsByCategory(categoryName string) []Item
	ListRulesItems() map[string]Item
}

type DataSchema struct {
	cache    common.IFastCache
	schema   SchemaInfo
	dataFile string
	sources  map[SourceKey]*Source
}

const ONLY_LOAD_ENABLED = true

func NewDataSchema(watcher common.IFileWatcher, cache common.IFastCache, configPath, dataPath string) *DataSchema {
	x := &DataSchema{
		cache:    cache,
		dataFile: dataPath,
		sources:  make(map[SourceKey]*Source),
	}

	schema := SchemaInfo{}
	common.LoadJson(configPath+"config.json", &schema)
	for _, source := range schema.Sources {

		if ONLY_LOAD_ENABLED && !source.Enabled {
			continue
		}

		log.Debug().Str("component", "schema").Str("source", source.Name).Msg("source load")
		sourceFile := x.dataFile + source.Name + ".json"
		dataSource := NewSource(sourceFile, source)
		dataSource.Refresh()
		x.sources[dataSource.Key()] = dataSource

		if watcher != nil {
			watcher.Add(sourceFile, dataSource.Refresh)
		}
	}

	x.schema = schema
	return x
}

func (x *DataSchema) Item(dataURI *DataURI) (Item, error) {
	dataSource := x.sources[SourceKey(dataURI.SourceName)]
	if dataSource == nil {
		return Item{}, ErrItemNotFound
	}
	return dataSource.Item(dataURI)
}

func (x *DataSchema) IsValidItem(dataURI *DataURI) bool {
	// if dataURI.CategoryName == CATEGORY_RULE && dataURI.SourceName == SOURCE_OSINTAMI_NAME {
	// 	return true
	// }
	_, err := x.Item(dataURI)
	return err == nil
}

func (x *DataSchema) IsValidCategory(categoryName string) bool {
	for _, name := range CATEGORIES {
		if name == categoryName {
			return true
		}
	}
	return false
}

func (x *DataSchema) IsEnabled(sourceName string) bool {
	return x.sources[SourceKey(sourceName)].Source().Enabled
}

func (x *DataSchema) Source(sourceName string) (SourceInfo, error) {
	source := x.sources[SourceKey(sourceName)]
	if source == nil {
		return SourceInfo{}, ErrSourceNotFound
	}
	return source.Source(), nil
}

func (x *DataSchema) ListItems() []Item {
	list, found := x.cache.Get("ListItems")
	if !found {
		items := []Item{}
		for _, source := range x.sources {
			items = append(items, source.ListItems()...)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Path < items[j].Path
		})
		x.cache.Set("ListItems", items, -1)
		return items
	}
	return list.([]Item)
}

func (x *DataSchema) ListSources() []SourceInfo {
	out := []SourceInfo{}
	for _, source := range x.sources {
		out = append(out, source.Source())
	}
	return out
}

func (x *DataSchema) ListCategories() []string {
	return CATEGORIES
}

func (x *DataSchema) ListItemsByCategory(categoryName string) []Item {
	out := []Item{}
	for _, item := range x.ListItems() {
		if categoryName == item.CategoryName {
			out = append(out, item)
		}
	}
	return out
}

func (x *DataSchema) ListRulesItems() map[string]Item {
	rules := make(map[string]Item)
	for _, item := range x.ListItems() {
		if item.CategoryName == CATEGORY_RULE || item.Query != "" {
			rules[item.Path] = item
		}
	}
	return rules
}
