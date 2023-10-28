// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFastCache(t *testing.T) {
	cache := NewFastCache()
	cache.Set("test key", "test value", -1)
	value, found := cache.Get("test key")
	assert.True(t, found)
	assert.Equal(t, "test value", value)

	err := cache.Save("/tmp/test.db")
	assert.Nil(t, err)

	cache.Delete("test key")
	value, found = cache.Get("test key")
	assert.False(t, found)
	assert.Nil(t, value)

	cache.LoadFile("/tmp/test.db")
	value, found = cache.Get("test key")
	assert.True(t, found)
	assert.Equal(t, "test value", value)

	cache.Clear()
	value, found = cache.Get("test key")
	assert.False(t, found)
	assert.Nil(t, value)

	os.Remove("/tmp/test.db")
}

func TestPersistentCache(t *testing.T) {
	cache := NewPersistentCache("/tmp/test.db")
	cache.Set("test key", "test value", -1)
	value, found := cache.Get("test key")
	assert.True(t, found)
	assert.Equal(t, "test value", value)

	cache.Persist()

	cache.Delete("test key")
	value, found = cache.Get("test key")
	assert.False(t, found)
	assert.Nil(t, value)

	cache.LoadFile("/tmp/test.db")
	value, found = cache.Get("test key")
	assert.True(t, found)
	assert.Equal(t, "test value", value)

	os.Remove("/tmp/test.db")
}

func TestFastCacheToJSON(t *testing.T) {
	cache := NewFastCache()
	cache.Set("test key", "test value", -1)
	err := cache.ToJSON("/tmp/test.json")
	assert.Nil(t, err)
	os.Remove("/tmp/test.json")
}
