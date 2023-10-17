// Copyright © 2023 OSINTAMI. This is not yours.
package server

import "github.com/osintami/fingerprintz/common"

// -----------------------------------------------------------------------------
type MockSecretsManager struct {
	keys map[string]string
}

func NewMockSecretsManager() common.ISecrets {
	return &MockSecretsManager{keys: make(map[string]string)}
}
func (x *MockSecretsManager) Set(key, value string) {
	x.keys[key] = value
}
func (x *MockSecretsManager) Find(key string) string {
	return x.keys[key]
}
