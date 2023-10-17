// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecrets(t *testing.T) {
	keys := []string{"TEST_API_KEY"}
	os.Setenv("TEST_API_KEY", "xxx")
	secrets := NewSecrets(keys)
	assert.Equal(t, "xxx", secrets.Find("TEST_API_KEY"))
	secrets.Set("key", "value")
	assert.Equal(t, "value", secrets.Find("key"))
	os.Clearenv()
}
