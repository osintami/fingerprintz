// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBiometricsScan(t *testing.T) {

	keys := make(map[string]string)
	keys["ip"] = "1.2.3.4"
	keys["email"] = "1@2.com"
	keys["ua"] = "test-user-agent"
	keys["hw"] = "test-hardware-id"
	keys["pid"] = "test-partner-id"
	fingerprint := BiometricsScan(context.TODO(), NewMockNods(false), keys)
	assert.Equal(t, "1.2.3.4", fingerprint.IpAddr)

	fingerprint = BiometricsScan(context.TODO(), NewMockNods(true), keys)
	assert.Nil(t, fingerprint)
}
