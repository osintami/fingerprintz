// Copyright © 2023 OSINTAMI. This is not yours.
package server

import "github.com/osintami/fingerprintz/common"

type MockNods struct {
	table map[string]bool
}

func NewMockNods(table map[string]bool) *MockNods {
	return &MockNods{table: table}
}

func (x *MockNods) Whoami(ip string) (*common.WhoamiInfo, error) {
	return nil, common.ErrItemNotFound
}

func (x *MockNods) Item(name string, keys map[string]string) (*common.DataOutput, error) {
	if x.table["Item"] == true {
		return &common.DataOutput{}, common.ErrItemNotFound
	}
	var value bool
	if keys["ip"] == "1.2.3.4" {
		value = false
	} else {
		value = true
	}
	return &common.DataOutput{Result: common.DataResult{Bool: &value}}, nil
}

func (x *MockNods) Fingerprint(keys map[string]string) (string, error) {
	return "", nil
}
