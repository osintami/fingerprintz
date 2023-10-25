// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVWriter(t *testing.T) {
	writer := NewCSVWriter()
	assert.Equal(t, "csv", writer.Type())

	testFile := "/tmp/csv-writer-test.csv"
	err := writer.Create(testFile)
	assert.Nil(t, err)

	row := []string{"1.2.3.4", "isBlacklisted"}

	err = writer.Insert(nil, row)
	assert.Nil(t, err)

	err = writer.Load(nil)
	assert.Nil(t, err)

	data, err := os.ReadFile(testFile)
	assert.Nil(t, err)
	assert.Equal(t, "1.2.3.4,isBlacklisted\n", string(data))

	// error path
	err = writer.Create("./././nope././.")
	assert.NotNil(t, err)

	// cleanup
	os.Remove(testFile)
}
