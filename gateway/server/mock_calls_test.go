// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"

	"gorm.io/gorm"
)

type MockCalls struct {
	table map[string]bool
}

func NewMockCalls(table map[string]bool) ICalls {
	return &MockCalls{table: table}
}

func (x *MockCalls) Call(ctx context.Context, call *Call) error {
	if x.table["Call"] == true {
		return gorm.ErrDuplicatedKey
	}
	return nil
}
