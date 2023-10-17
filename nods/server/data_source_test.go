// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func createSource(sourceName string, sourceType string) SourceInfo {
	return SourceInfo{
		Name:     sourceName,
		Database: sourceType,
		Enabled:  true}
}

func TestDataSource(t *testing.T) {
	// sourceInfo comes from config.json
	sourceInfo := createSource("ipsum", "mmdb")

	// create a source for the data dictionary
	source := NewSource("./test/ipsum.json", sourceInfo)
	assert.Equal(t, SourceKey("ipsum"), source.Key())

	// load the source configuration file
	source.Refresh()

	// check for item
	uri := NewDataURI(CATEGORY_IPADDR, "ipsum", "blacklist.isBlacklisted")
	item, err := source.Item(uri)
	assert.Nil(t, err)

	assert.Equal(t, CATEGORY_IPADDR, item.CategoryName)
	assert.Equal(t, "ipsum", item.SourceName)
	assert.Equal(t, "ip/ipsum/blacklist.isBlacklisted", item.Path)
	assert.Equal(t, common.Boolean, item.Type)
	assert.Equal(t, "ipsum.blacklist.isBlacklisted", item.Gjson)
	assert.Equal(t, ItemKey("ip/ipsum/blacklist.isBlacklisted"), item.Key())
}

func TestDataSourceDoesNotExist(t *testing.T) {
	sourceInfo := createSource("ipsum", "mmdb")
	source := NewSource("./test/ipsum.json", sourceInfo)

	// check for item
	uri := NewDataURI(CATEGORY_IPADDR, "ipsum", "nope.isBlacklisted")
	item, err := source.Item(uri)
	assert.Equal(t, ErrItemNotFound, err)
	assert.NotNil(t, item)
}

func TestDataSourceDisabledWithEnabledItem(t *testing.T) {
	sourceInfo := createSource("nope", "mmdb")
	sourceInfo.Enabled = false
	source := NewSource("./test/nope.json", sourceInfo)
	source.Refresh()

	// check for item
	uri := NewDataURI(CATEGORY_IPADDR, "nope", "isDisabled")
	item, err := source.Item(uri)
	assert.Nil(t, err)

	assert.Equal(t, CATEGORY_IPADDR, item.CategoryName)
	assert.Equal(t, "nope", item.SourceName)
	assert.Equal(t, "ip/nope/isDisabled", item.Path)
	assert.Equal(t, common.Boolean, item.Type)
	assert.Equal(t, "isDisabled", item.Gjson)
	assert.Equal(t, ItemKey("ip/nope/isDisabled"), item.Key())
}
