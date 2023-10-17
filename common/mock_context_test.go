// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"context"
	"time"
)

type MockContext struct {
}

func NewMockContext() context.Context {
	return &MockContext{}
}

func (x *MockContext) Deadline() (deadline time.Time, ok bool) {
	return time.Now(), true
}

func (x *MockContext) Done() <-chan struct{} {
	return nil
}
func (x *MockContext) Err() error {
	return nil
}
func (x *MockContext) Value(key any) any {
	return nil
}
