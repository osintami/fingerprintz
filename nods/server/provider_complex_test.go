// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestComplexProvider(t *testing.T) {
	provider := NewComplexProvider(NewMockRuleProvider(false))
	assert.False(t, provider.IsCached())

	inputs := make(map[string]string)
	inputs[CATEGORY_IPADDR] = "1.2.3.4"
	inputs[common.INPUT_TYPE] = "JSON"
	inputs[common.INPUT_RULE] = "[ip/ipsum/blacklist.isBlacklisted]||[ip/uhb/blacklist.isBlacklisted]"

	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_RULE, inputs)
	assert.Nil(t, err)

	var eval EvalOutput
	json.Unmarshal(data, &eval)
	assert.Equal(t, true, eval.Result.Bool)
}

func TestComplexProviderInvalidCategory(t *testing.T) {
	provider := NewComplexProvider(NewMockRuleProvider(false))

	inputs := make(map[string]string)
	inputs[CATEGORY_PHONE] = "18001234567"

	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Equal(t, "{}", string(data))
}

func TestComplexProviderMissingRule(t *testing.T) {
	provider := NewComplexProvider(NewMockRuleProvider(false))

	inputs := make(map[string]string)
	inputs[CATEGORY_RULE] = ""

	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_RULE, inputs)
	assert.Equal(t, ErrMissingInputs, err)
	assert.Equal(t, "{}", string(data))
}

func TestComplexProviderEvaluateFailure(t *testing.T) {
	provider := NewComplexProvider(NewMockRuleProvider(true))

	inputs := make(map[string]string)
	inputs[CATEGORY_IPADDR] = "1.2.3.4"
	inputs[common.INPUT_TYPE] = "JSON"
	inputs[common.INPUT_RULE] = "[ip/ipsum/blacklist.isBlacklisted]||[ip/uhb/blacklist.isBlacklisted]"

	data, err := provider.CategoryInfo(context.TODO(), CATEGORY_RULE, inputs)
	assert.Equal(t, ErrItemNotFound, err)
	assert.Equal(t, "", string(data))
}
