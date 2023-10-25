// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCachedProvider(t *testing.T) {
	tools := mockToolbox()

	provider := NewCachedProvider(tools, "nope", NewMockProvider())
	assert.True(t, provider.IsCached())

	inputs := make(map[string]string)
	inputs[CATEGORY_IPADDR] = "1.2.3.4"

	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_IPADDR, inputs)
	assert.Nil(t, err)
	assert.Equal(t, "{}", string(data))

}
