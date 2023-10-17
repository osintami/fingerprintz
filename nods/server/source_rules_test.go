// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestRules(t *testing.T) {
	source := NewRuleSource(mockToolbox(), NewMockDataRouter(false), NewMockDataSchema())

	inputs := make(map[string]string)
	inputs[CATEGORY_IPADDR] = "1.2.3.4"
	inputs[common.INPUT_TYPE] = "JSON"
	inputs[common.INPUT_RULE] = "rule/osintami/isBlacklisted"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_RULE, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "isBlacklisted")
	assert.Equal(t, true, result.Bool())
}

func TestRulesInvalidCategory(t *testing.T) {
	source := NewRuleSource(mockToolbox(), NewMockDataRouter(false), NewMockDataSchema())
	inputs := make(map[string]string)
	// invalid category
	inputs[CATEGORY_PHONE] = "18001234567"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_PHONE, inputs)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Equal(t, "{}", string(data))
}

func TestRulesWithDate(t *testing.T) {
	source := NewRuleSource(mockToolbox(), NewMockDataRouter(false), NewMockDataSchema())

	inputs := make(map[string]string)
	inputs[CATEGORY_EMAIL] = "1@2.com"
	inputs[common.INPUT_TYPE] = "JSON"
	inputs[common.INPUT_RULE] = "rule/osintami/email.hasRecentBreach"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_RULE, inputs)
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "isCompromised")
	assert.Equal(t, true, result.Bool())
}

func TestRulesBadSyntax(t *testing.T) {
	source := NewRuleSource(mockToolbox(), NewMockDataRouter(false), NewMockDataSchema())

	inputs := make(map[string]string)
	inputs[CATEGORY_EMAIL] = "1@2.com"
	inputs[common.INPUT_TYPE] = "JSON"
	inputs[common.INPUT_RULE] = "rule/osintami/expression"
	data, err := source.CategoryInfo(context.TODO(), CATEGORY_RULE, inputs)
	assert.NotNil(t, err)
	assert.Equal(t, "{}", string(data))
}

// TODO:  nested rules, string and float data types in rules
