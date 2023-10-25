// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type FastDBProvider struct {
	dbFile string
	reader *common.FastCache
}

func NewFastDBProvider(tools *Toolbox, dbName string) (IDataProvider, error) {
	x := &FastDBProvider{
		dbFile: tools.DataPath + dbName + ".fast",
		reader: common.NewFastCache(),
	}
	x.refresh()
	if tools.Watcher != nil {
		_ = tools.Watcher.Add(x.dbFile, x.refresh)
	}
	return x, nil
}

func (x *FastDBProvider) refresh() {
	x.reader = x.reader.LoadFile(x.dbFile)
}

func (x *FastDBProvider) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	result, found := x.reader.Get(inputs[categoryName])
	if !found {
		log.Warn().Str("component", "fastdb item provider").Str("source", x.dbFile).Str("key", inputs[categoryName]).Msg("category info")
		return nil, common.ErrNoDataPresent
	}
	return result.([]uint8), nil
}

func (x *FastDBProvider) IsCached() bool {
	return false
}
