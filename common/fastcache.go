// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type IFastCache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, duration time.Duration)
	Clear()
	Persist()
}

type FastCache struct {
	persistFile string
	cache       *cache.Cache
}

func NewFastCache() *FastCache {
	return &FastCache{cache: cache.New(24*time.Hour, 60*time.Minute)}
}

func NewPersistentCache(persistFile string) *FastCache {
	x := &FastCache{
		persistFile: persistFile,
		cache:       cache.New(24*time.Hour, 60*time.Minute)}
	x.cache.LoadFile(persistFile)
	return x
}

func (x *FastCache) Get(key string) (interface{}, bool) {
	return x.cache.Get(key)
}

func (x *FastCache) Set(key string, value interface{}, duration time.Duration) {
	x.cache.Set(key, value, duration)
}

func (x *FastCache) LoadFile(fileName string) *FastCache {
	x.cache.LoadFile(fileName)
	return x
}

func (x *FastCache) Save(fileName string) error {
	return x.cache.SaveFile(fileName)
}

func (x *FastCache) Persist() {
	x.cache.SaveFile(x.persistFile)
}

func (x *FastCache) Clear() {
	for k := range x.cache.Items() {
		x.cache.Delete(k)
	}
}

func (x *FastCache) Delete(pattern string) {
	for k := range x.cache.Items() {
		if strings.Contains(k, pattern) {
			x.cache.Delete(k)
		}
	}
}
