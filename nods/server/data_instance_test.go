// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func createDataInstance(nodata bool) IDataSource {
	tools := mockToolbox()
	reader := NewMockMaxmindReader("./test/ipsum.mmdb", nodata)
	provider := NewMMDBProvider(tools.Watcher, reader)
	return NewDataInstance(tools, "ipsum", provider)
}

func createItemInstance() Item {
	return Item{
		Path:         "ip/ipsum/blacklist.isBlacklisted",
		CategoryName: "ip",
		SourceName:   "ipsum",
		Enabled:      true,
		Gjson:        "ipsum.blacklist.isBlacklisted",
		Description:  "",
		TypeName:     "Boolean",
		Type:         common.Boolean}
}

func TestDataInstance(t *testing.T) {
	instance := createDataInstance(false)
	assert.False(t, instance.IsCached())
	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"
	item := createItemInstance()
	data, err := instance.ItemValue(context.TODO(), item, inputs)
	assert.Nil(t, err)
	assert.True(t, data.Result.Boolean())
}

func TestDataInstanceInvalidCategory(t *testing.T) {
	instance := createDataInstance(false)
	assert.False(t, instance.IsCached())
	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"
	data, err := instance.CategoryInfo(context.TODO(), "nope", inputs)
	assert.Equal(t, ErrCategoryNotFound, err)
	assert.Equal(t, "{}", string(data))
}

func TestDataInstanceRules(t *testing.T) {
	instance := createDataInstance(false)
	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"
	item := createItemInstance()
	item.Type = common.JSON
	item.Gjson = PASSTHROUGH
	data, err := instance.ItemValue(context.TODO(), item, inputs)
	assert.Nil(t, err)
	assert.False(t, data.Result.Boolean())

}

func TestDataInstanceNoDataFound(t *testing.T) {
	instance := createDataInstance(true)
	inputs := common.DataInputs{}
	inputs["ip"] = "0.0.0.0"
	item := createItemInstance()
	data, err := instance.ItemValue(context.TODO(), item, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.False(t, data.Result.Boolean())
}

func TestDataInstanceBadGjson(t *testing.T) {
	instance := createDataInstance(false)
	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"
	item := createItemInstance()
	item.Gjson = "nope"
	data, err := instance.ItemValue(context.TODO(), item, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.False(t, data.Result.Boolean())
}

func TestDataInstanceItemNotFound(t *testing.T) {
	instance := createDataInstance(true)
	inputs := common.DataInputs{}
	inputs["ip"] = "0.0.0.0"
	item := createItemInstance()
	item.SourceName = "nope"
	data, err := instance.ItemValue(context.TODO(), item, inputs)
	assert.Equal(t, common.ErrNoDataPresent, err)
	assert.False(t, data.Result.Boolean())
}
