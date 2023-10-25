// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"os"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestFastDBWriter(t *testing.T) {
	writer := NewFastDBWriter()
	assert.Equal(t, "fast", writer.Type())

	testFile := "/tmp/fast-writer-test.fast"
	err := writer.Create(testFile)
	assert.Nil(t, err)

	row := `{"isBlacklisted":true}`

	err = writer.Insert("1.2.3.4", row)
	assert.Nil(t, err)

	err = writer.Load(nil)
	assert.Nil(t, err)

	cache := common.NewFastCache()
	cache = cache.LoadFile(testFile)
	assert.Nil(t, err)
	assert.NotNil(t, cache)

	entry, found := cache.Get("1.2.3.4")
	assert.True(t, found)
	assert.Equal(t, row, entry.(string))

	// cleanup
	os.Remove(testFile)
}
