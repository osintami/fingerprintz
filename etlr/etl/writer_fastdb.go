// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/patrickmn/go-cache"
)

type FastDBWriter struct {
	writer   *cache.Cache
	fileName string
}

func NewFastDBWriter() IWriter {
	return &FastDBWriter{}
}

func (x *FastDBWriter) Type() string {
	return "fast"
}

func (x *FastDBWriter) Create(fileName string) error {
	x.fileName = fileName
	x.writer = cache.New(-1, -1)
	return nil
}

func (x *FastDBWriter) Load(job IETLJob) error {
	return x.writer.SaveFile(x.fileName)
}

func (x *FastDBWriter) Insert(key interface{}, value interface{}) error {
	x.writer.Set(key.(string), value, -1)
	return nil
}
