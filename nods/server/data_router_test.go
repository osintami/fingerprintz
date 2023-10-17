// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestDataRouterDataValue(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = "./test/"
	router := NewDataRouter(tools)
	router.Init()

	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"

	// success
	uri := NewDataURI(CATEGORY_IPADDR, "ipsum", "blacklist.isBlacklisted")
	value, err := router.DataValue(context.TODO(), uri, inputs)
	assert.Nil(t, err)
	assert.True(t, value.Result.Boolean())
}

func TestDataRouterDataValueInvalidCategory(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = "./test/"
	router := NewDataRouter(tools)
	router.Init()

	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"

	// invalid category
	uri := NewDataURI("nope", "ipsum", "blacklist.isBlacklisted")
	value, err := router.DataValue(context.TODO(), uri, inputs)
	assert.Equal(t, ErrCategoryNotFound, err)
	assert.False(t, value.Result.Boolean())
}

func TestDataRouterSourceNotFound(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = "./test/"
	router := NewDataRouter(tools)
	router.Init()

	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"

	// source not found
	uri := NewDataURI(CATEGORY_IPADDR, "nope", "blacklist.isBlacklisted")
	value, err := router.DataValue(context.TODO(), uri, inputs)
	assert.Equal(t, ErrSourceNotFound, err)
	assert.False(t, value.Result.Boolean())
}

func TestDataRouterItemNotFound(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = "./test/"
	router := NewDataRouter(tools)
	router.Init()

	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"

	// item not found
	uri := NewDataURI(CATEGORY_IPADDR, "ipsum", "blklst.isBlacklisted")
	value, err := router.DataValue(context.TODO(), uri, inputs)
	assert.Equal(t, ErrItemNotFound, err)
	assert.False(t, value.Result.Boolean())
}

func TestDataRouterCategoryMismatch(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = "./test/"
	router := NewDataRouter(tools)
	router.Init()

	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"

	// category mismatch for item
	uri := NewDataURI(CATEGORY_PHONE, "ipsum", "blacklist.isBlacklisted")
	value, err := router.DataValue(context.TODO(), uri, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.False(t, value.Result.Boolean())
}

func TestDataRouterDataValueNotFound(t *testing.T) {
	tools := mockToolbox()
	tools.DataPath = "./test/"
	router := NewDataRouter(tools)
	router.Init()

	// data not present
	inputs := common.DataInputs{}
	inputs["ip"] = "0.0.0.0"
	uri := NewDataURI(CATEGORY_IPADDR, "ipsum", "blacklist.isBlacklisted")
	value, err := router.DataValue(context.TODO(), uri, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.False(t, value.Result.Boolean())
}

func TestDataRouterCategoryValue(t *testing.T) {
	router := NewDataRouter(mockToolbox())
	router.Init()

	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"

	// success
	values, err := router.CategoryValues(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(values))
}

func TestDataRouterCategoryValueNotFound(t *testing.T) {
	router := NewDataRouter(mockToolbox())
	router.Init()

	inputs := common.DataInputs{}
	inputs["ip"] = "0.0.0.0"

	// success
	values, err := router.CategoryValues(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, ErrItemNotFound, err)
	assert.Equal(t, 0, len(values))
}

func TestDataRouterCodeInstances(t *testing.T) {
	router := NewDataRouter(mockToolbox())
	router.Init()

	osintami := router.newInstance(SOURCE_OSINTAMI_NAME)
	pwned := router.newInstance(SOURCE_PWNED_NAME)
	whois := router.newInstance(SOURCE_WHOIS_NAME)
	spamhaus := router.newInstance(SOURCE_SPAMHAUS_NAME)
	ipinfo := router.newInstance(SOURCE_IPINFO_NAME)

	assert.False(t, osintami.IsCached())
	assert.True(t, pwned.IsCached())
	assert.True(t, whois.IsCached())
	assert.True(t, spamhaus.IsCached())
	assert.True(t, ipinfo.IsCached())

	nope := router.newInstance("nope")
	assert.Nil(t, nope)

	// TODO:  test a data point from osintami
}

func TestDataRouterAPI(t *testing.T) {
	// TODO:
}

func TestDataRouterFastDB(t *testing.T) {
	// TODO:
}
