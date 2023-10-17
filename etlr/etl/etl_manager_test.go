// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestETLManager(t *testing.T) {
	tools := fillToolbox(nil)
	dataPath := "./test/data/"
	source := Source{
		Name:       "test",
		Enabled:    true,
		File:       "./test/source/test.csv",
		InputType:  "csv",
		OutputType: "mmdb",
		Separator:  ",",
	}
	sources := []Source{source}
	manager := NewETLManager(tools, dataPath, sources)

	cron := manager.ScheduleCronJobs()
	assert.Equal(t, 3, len(cron.Entries()))

	manager.refreshDaily()
	manager.refreshHourly()
	manager.refreshWeekly()

	assert.Equal(t, "test", manager.Source("test").Name)
	assert.Nil(t, manager.Source("nope"))

	err := manager.Refresh("test")
	assert.Nil(t, err)

	err = manager.Refresh("nope")
	assert.Equal(t, ErrVendorNotFound, err)

	source.Name = "nope"
	job, err := manager.createInstance(source)
	assert.Equal(t, ErrVendorNotFound, err)
	assert.Nil(t, job)

	manager.RefreshAll()
}

func TestETLManagerInstances(t *testing.T) {
	client := resty.New()
	tools := fillToolbox(client)
	dataPath := "/tmp/dbs"
	// load all production sources for this test
	sources := []Source{}
	err := common.LoadJson("../config.json", &sources)
	assert.Nil(t, err)

	manager := NewETLManager(tools, dataPath, sources)
	assert.NotNil(t, manager.FindJob("amazon"))
}
