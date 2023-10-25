// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestFileWriter(t *testing.T) {
	writer := NewFileDBWriter("csv")
	assert.Equal(t, "csv", writer.Type())

	testFile := "/tmp/test.csv"

	// no-op
	err := writer.Create("/tmp/nope.csv")
	assert.Nil(t, err)
	// file should not exist
	_, err = os.Stat("/tmp/nope.csv")
	assert.NotNil(t, err)

	// no-op
	row := []string{"0.0.0.0", "isBlacklisted"}
	err = writer.Insert(nil, row)
	assert.Nil(t, err)

	job := createETLJob(resty.New(), "csv", "csv")
	// normally done in the etlr job setup
	os.Mkdir("/tmp/test", 01777)
	// normally done in the etlr job extract
	data, err := os.ReadFile("/tmp/test.csv")
	assert.Nil(t, err)
	err = os.WriteFile("/tmp/test/test.csv", data, 0644)
	assert.Nil(t, err)

	err = writer.Load(job)
	assert.Nil(t, err)

	data, err = os.ReadFile(job.Info().snapshotFile)
	assert.Nil(t, err)
	assert.Equal(t, TEST_CONTENT, string(data))

	// cleanup
	os.Remove(testFile)
}
