// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

// NOTE:  these tests rely on the files in ./test/ to work

func TestSchema(t *testing.T) {
	schema := NewDataSchema(
		NewMockWatcher(),
		NewMockCache(),
		"./test/",
		"./test/")

	// item is valid
	assert.True(t, schema.IsValidItem(NewItemSplitter("ip/ipsum/blacklist.isBlacklisted")))

	// item is enabled
	assert.True(t, schema.IsEnabled("ipsum"))

	// item has a valid source
	source, err := schema.Source("ipsum")
	assert.Nil(t, err)
	assert.Equal(t, source.Name, "ipsum")
	assert.True(t, source.Enabled)

	// list all available items
	items := schema.ListItems()
	assert.Equal(t, 4, len(items))

	// list all available sources
	sources := schema.ListSources()
	assert.Equal(t, 3, len(sources))

	// list all available categories
	categories := schema.ListCategories()
	assert.Equal(t, 7, len(categories))

	// list all items by category
	items = schema.ListItemsByCategory(CATEGORY_BROWSER)
	assert.Equal(t, 1, len(items))

	// list all rules
	rules := schema.ListRulesItems()
	assert.Equal(t, 1, len(maps.Keys(rules)))
}

func TestSchemaValidCategory(t *testing.T) {
	schema := NewDataSchema(
		NewMockWatcher(),
		NewMockCache(),
		"./test/",
		"./test/")

	assert.True(t, schema.IsValidCategory(CATEGORY_IPADDR))
	assert.False(t, schema.IsValidCategory("nope"))
}

func TestDataItem(t *testing.T) {
	schema := NewDataSchema(
		NewMockWatcher(),
		NewMockCache(),
		"./test/",
		"./test/")
	dataURI := NewItemSplitter("ip/ipsum/blacklist.isBlacklisted")
	item, err := schema.Item(dataURI)
	assert.Nil(t, err)
	assert.Equal(t, item.CategoryName, CATEGORY_IPADDR)
	assert.Equal(t, item.SourceName, "ipsum")
	assert.Equal(t, item.Path, "ip/ipsum/blacklist.isBlacklisted")
}

func TestSchemaInvalidItem(t *testing.T) {
	schema := NewDataSchema(
		NewMockWatcher(),
		NewMockCache(),
		"./test/",
		"./test/")
	dataURI := NewItemSplitter("nope")
	item, err := schema.Item(dataURI)
	assert.Equal(t, ErrItemNotFound, err)
	assert.Equal(t, "", item.Path)
}

func TestSchemaInvalidSource(t *testing.T) {
	schema := NewDataSchema(
		NewMockWatcher(),
		NewMockCache(),
		"./test/",
		"./test/")
	source, err := schema.Source("nope")
	assert.Equal(t, ErrSourceNotFound, err)
	assert.Equal(t, "", source.Name)
}
