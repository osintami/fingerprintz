// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"encoding/json"
	"net"

	"github.com/osintami/fingerprintz/common"
)

type MockMaxmindReader struct {
	fileName string
	fail     bool
}

func NewMockMaxmindReader(fileName string, fail bool) common.IMaxmindReader {
	return &MockMaxmindReader{
		fileName: fileName,
		fail:     fail,
	}
}

func (x *MockMaxmindReader) Resync() {

}

func (x *MockMaxmindReader) FileName() string {
	return x.fileName
}

func (x *MockMaxmindReader) Lookup(ip net.IP) (json.RawMessage, error) {
	if x.fail {
		return nil, common.ErrNoDataPresent
	}
	// TODO:  support uhb and maxmind for other tests
	content := "{\"ipsum\":{\"blacklist\":{\"isBlacklisted\":true,\"lastReportedAt\":\"1942-01-01 00:00:00\"}"
	return []byte(content), nil
}
