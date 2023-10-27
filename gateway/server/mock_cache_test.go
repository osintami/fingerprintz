// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"os"
	"time"

	"github.com/osintami/fingerprintz/common"
)

// -----------------------------------------------------------------------------
type MockCache struct {
	cache map[string]interface{}
	fails map[string]bool
}

func NewMockCache(fails map[string]bool) common.IFastCache {
	return &MockCache{cache: make(map[string]interface{}), fails: fails}
}
func (x *MockCache) Get(key string) (interface{}, bool) {
	item := x.cache[key]
	if item != nil {
		return item, true
	}
	return nil, false
}
func (x *MockCache) Set(key string, value interface{}, duration time.Duration) {
	x.cache[key] = value
}
func (x *MockCache) Clear() {
	for k := range x.cache {
		delete(x.cache, k)
	}
}
func (x *MockCache) Persist() {
}
func (x *MockCache) ToJSON(fileName string) error {
	if x.fails != nil && x.fails["ToJSON"] {
		return os.ErrPermission
	}
	return nil
}
