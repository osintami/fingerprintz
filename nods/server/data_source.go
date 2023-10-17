package server

import (
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"golang.org/x/exp/maps"
)

type Item struct {
	//URI          *DataURI        `json:"-"`
	Path         string          `json:"Item"`
	CategoryName string          `json:"-"`
	SourceName   string          `json:"-"`
	Enabled      bool            `json:"Enabled"`
	Gjson        string          `json:"GJSON"`
	Description  string          `json:"Description"`
	TypeName     string          `json:"Type"`
	Type         common.DataType `json:"-"`
	Query        string          `json:"Query,omitempty"`
}

type Source struct {
	items    map[ItemKey]Item
	dataFile string
	source   SourceInfo
}

type SourceKey string

func NewSource(dataFile string, source SourceInfo) *Source {
	return &Source{
		items:    make(map[ItemKey]Item),
		dataFile: dataFile,
		source:   source}
}

func (x *Source) Key() SourceKey {
	return SourceKey(x.source.Name)
}

func (x *Source) Refresh() {
	items := []Item{}
	common.LoadJson(x.dataFile, &items)
	for _, item := range items {
		uri := NewItemSplitter(item.Path)
		item.CategoryName = uri.CategoryName
		item.SourceName = uri.SourceName
		item.Type = item.Type.ToDataType(item.TypeName)
		if item.Enabled && !x.source.Enabled {
			item.Enabled = false
		}
		log.Debug().Str("component", "schema").Str("item", item.Path).Msg("item load")
		x.items[item.Key()] = item
	}
}

func (x *Source) Item(dataURI *DataURI) (Item, error) {
	if val, ok := x.items[ItemKey(dataURI.Key())]; ok {
		return val, nil
	}
	return Item{}, ErrItemNotFound
}

func (x *Source) ListItems() []Item {
	return maps.Values(x.items)
}

func (x *Source) Source() SourceInfo {
	return x.source
}

type ItemKey string

func (x Item) Key() ItemKey {
	return ItemKey(x.Path)
}
