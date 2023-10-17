// Copyright © 2023 OSINTAMI. This is not yours.
package server

import "github.com/osintami/fingerprintz/common"

// -----------------------------------------------------------------------------
type MockRealtimeBlackholeList struct {
}

func NewMockRealtimeBlackholeList() common.IRealtimeBlackholeList {
	return &MockRealtimeBlackholeList{}
}
func (x MockRealtimeBlackholeList) Lookup(host, addr string) bool {
	return true
}
