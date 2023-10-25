// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRBL(t *testing.T) {
	rbl := NewRealtimeBlackholeList()
	found := rbl.Lookup("fingerprintz.spamhaus.org", "8.8.8.2")
	assert.False(t, found)
}
