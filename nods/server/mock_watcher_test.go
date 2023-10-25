// Copyright © 2023 OSINTAMI. This is not yours.
package server

import "github.com/osintami/fingerprintz/common"

// -----------------------------------------------------------------------------
type MockWatcher struct {
}

func NewMockWatcher() common.IFileWatcher {
	return &MockWatcher{}
}
func (x *MockWatcher) Add(file string, refresh func()) error {
	return nil
}
func (x *MockWatcher) Listen() {
}
