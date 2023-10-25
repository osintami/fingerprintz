// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"net"

	"github.com/osintami/fingerprintz/common"
)

type MMDBProvider struct {
	reader common.IMaxmindReader
}

func NewMMDBProvider(watcher common.IFileWatcher, reader common.IMaxmindReader) IDataProvider {
	x := &MMDBProvider{
		reader: reader}

	if watcher != nil {
		watcher.Add(reader.FileName(), x.reader.Resync)
	}
	return x
}

func (x *MMDBProvider) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	if categoryName != CATEGORY_IPADDR {
		return nil, ErrNotImplemented
	}
	return x.reader.Lookup(net.ParseIP(inputs[categoryName]))
}

func (x *MMDBProvider) IsCached() bool {
	return false
}
