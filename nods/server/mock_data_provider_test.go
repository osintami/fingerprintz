// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"

	"github.com/osintami/fingerprintz/common"
)

// -----------------------------------------------------------------------------
type MockDataProvider struct {
}

func NewMockProvider() IDataProvider {
	return &MockDataProvider{}
}

func (x *MockDataProvider) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	switch categoryName {
	case CATEGORY_IPADDR:
		return []byte("{}"), nil
	case CATEGORY_PHONE:
		return []byte("{}"), common.ErrNoDataPresent
	}
	return nil, ErrNotImplemented
}

func (x *MockDataProvider) IsCached() bool {
	return true
}
