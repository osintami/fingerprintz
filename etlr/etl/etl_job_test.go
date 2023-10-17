// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/etlr/utils"
	"github.com/stretchr/testify/assert"
)

func fillToolbox(client *resty.Client) *Toolbox {
	tools := &Toolbox{
		Network:    utils.NewNetworkingHelper(client),
		FileSystem: utils.NewFSHelper(),
		Secrets:    common.NewSecrets([]string{"TEST_API_KEY"}),
		CSV:        utils.NewCSVReader(),
		Items:      make(map[string]Item)}
	tools.Secrets.Set("TEST_API_KEY", "GHSAT0AAAAAACE25RHZPF6EBDUUB7JY4F7EZIS6QSA")
	return tools
}

var TEST_CONTENT = `1.2.3.4\nnope\n34.31.171.21\n`

func createETLJob(client *resty.Client, inputType, outputType string) *ETLJob {
	tools := fillToolbox(client)

	dataPath := "/tmp/dbs/"
	var writer IWriter
	var load ILoad

	switch outputType {
	case "mmdb":
		writer = NewMMDBWriter()
	case "csv":
		writer = NewCSVWriter()
	case "fast":
		writer = NewFastDBWriter()
	}

	// TODO:  create the Name: test-{timestamp millis}
	os.Remove("/tmp/test.csv")
	os.RemoveAll("/tmp/dbs/")
	os.Mkdir("/tmp/dbs", 01777)
	os.WriteFile("/tmp/test.csv", []byte(TEST_CONTENT), 0644)

	extract := NewFileExtractor("/tmp/test.csv")
	transform := NewTestSource(writer)

	load = writer

	source := &Source{
		Name:       "test",
		Enabled:    true,
		URL:        "https://raw.githubusercontent.com/osintami/fingerprintz/master/etlr/etl/test/test.csv?token={key}",
		ApiKey:     "TEST_API_KEY",
		InputType:  inputType,
		OutputType: outputType,
		Separator:  ",",
	}

	return NewETLJob(tools, source, dataPath, writer, extract, transform, load)
}

func cleanup() {
	os.Remove("/tmp/test.csv")
	os.RemoveAll("/tmp/dbs/")
}

func TestETLJob(t *testing.T) {
	job := createETLJob(resty.New(), "csv", "mmdb")

	// 1) create working directory /tmp/test
	// 2) create the writer (mmdb, csv, or fast)
	err := job.prepareETL()
	assert.Nil(t, err)
	// test for working directory
	_, err = os.Stat("/tmp/test")
	assert.Nil(t, err)

	// copy data from source and write to /tmp/test/test.csv
	err = job.extract.Extract(job)
	assert.Nil(t, err)
	_, err = os.Stat("/tmp/test/test.csv")
	assert.Nil(t, err)

	// transform from csv to mmdb in memory
	err = job.transform.Transform(job)
	assert.Nil(t, err)

	// load the mmdb file to disk (test_YYYY_MM_DD.mmdb)
	err = job.load.Load(job)
	assert.Nil(t, err)
	_, err = os.Stat("/tmp/test/" + job.Info().snapshotName)
	assert.Nil(t, err)

	// publish the mmdb file to final destination (test.mmdb)
	err = job.publishDatums()
	assert.Nil(t, err)
	_, err = os.Stat(job.Info().outputFile)
	assert.Nil(t, err)

	// publish the schema file to final destination (test.json)
	err = job.publishSchema()
	assert.Nil(t, err)
	_, err = os.Stat(job.Info().schemaFile)
	assert.Nil(t, err)

	// publish an empty schema
	for k := range job.tools.Items {
		delete(job.tools.Items, k)
	}
	assert.Equal(t, ErrEmptySchema, job.publishSchema())

	// clean up working directory and schema cache
	err = job.cleanupETL()
	assert.Nil(t, err)
	_, err = os.Stat("/tmp/test")
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(job.Tools().Items))

	cleanup()
}

func TestETLJobRefreshMMDB(t *testing.T) {
	os.Remove("/temp/test.csv")

	// create source file
	err := os.WriteFile("/tmp/test.csv", []byte(TEST_CONTENT), 0644)
	assert.Nil(t, err)

	job := createETLJob(resty.New(), "csv", "mmdb")
	err = job.Refresh()
	assert.Nil(t, err)

	_, err = os.Stat("/tmp/test.csv")
	assert.Nil(t, err)

	cleanup()
}

func TestETLJobPaths(t *testing.T) {
	job := createETLJob(resty.New(), "csv", "mmdb")
	// test file names
	assert.Equal(t, "/tmp/test/", job.Info().workingPath)
	assert.Equal(t, "/tmp/test/test.csv", job.Info().inputFile)
	assert.Equal(t, "/tmp/dbs/test.mmdb", job.Info().outputFile)
	assert.Equal(t, "/tmp/dbs/test.json", job.Info().schemaFile)
}
