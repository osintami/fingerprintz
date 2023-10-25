// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"

	"github.com/osintami/fingerprintz/common"
)

type MockRuleProvider struct {
	fail bool
}

func NewMockRuleProvider(fail bool) IRuleProvider {
	return &MockRuleProvider{fail: fail}
}

func (x *MockRuleProvider) IsCached() bool {
	return false
}

func (x *MockRuleProvider) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	return []byte("{\"Value\":true}"), nil
}

func (x *MockRuleProvider) Evaluate(ctx context.Context, nods Item, inputs common.DataInputs) ([]byte, error) {
	if x.fail {
		return nil, ErrItemNotFound
	}
	return []byte("{\"Value\":true}"), nil
}
