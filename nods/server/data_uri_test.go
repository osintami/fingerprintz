// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataUri(t *testing.T) {
	uri := &DataURI{CategoryName: CATEGORY_IPADDR, SourceName: SOURCE_OSINTAMI_NAME, ItemName: "IsBlacklisted", URI: "ip/osintami/IsBlacklisted"}
	assert.Equal(t, uri.Key(), "ip/osintami/IsBlacklisted")

	uri = NewDataURI(CATEGORY_IPADDR, SOURCE_OSINTAMI_NAME, "IsBlacklisted")
	assert.Equal(t, uri.CategoryName, CATEGORY_IPADDR)
	assert.Equal(t, uri.SourceName, SOURCE_OSINTAMI_NAME)
	assert.Equal(t, uri.ItemName, "IsBlacklisted")

	assert.Equal(t, uri.Key(), uri.URI)
	assert.Equal(t, uri.Key(), "ip/osintami/IsBlacklisted")

	uri = NewItemSplitter("/nope/blacklist.isBlacklisted")
	assert.Equal(t, "unknown/unknown/unknown", uri.Key())
}
