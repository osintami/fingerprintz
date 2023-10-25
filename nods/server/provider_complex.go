// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/osintami/fingerprintz/common"
	"github.com/tidwall/gjson"
)

type EvalResult struct {
	Type int
	Bool bool
}

type EvalOutput struct {
	Item   string
	Result EvalResult
	Keys   common.DataInputs
	Error  string
}

type ComplexProvider struct {
	rules IRuleProvider
}

func NewComplexProvider(rules IRuleProvider) IDataProvider {
	return &ComplexProvider{rules: rules}
}

func (x *ComplexProvider) IsCached() bool {
	return false
}

func (x *ComplexProvider) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	if categoryName != CATEGORY_RULE {
		return []byte("{}"), ErrNotImplemented
	}

	// short circuit on missing data key
	if inputs[categoryName] == "" {
		return []byte("{}"), ErrMissingInputs
	}

	query := inputs[CATEGORY_RULE]

	query = strings.ReplaceAll(query, " OR ", " || ")
	query = strings.ReplaceAll(query, "]OR[", "] || [")
	query = strings.ReplaceAll(query, "]||[", "] || [")
	query = strings.ReplaceAll(query, "]|[", "] || [")
	query = strings.ReplaceAll(query, "] | [", "] || [")

	query = strings.ReplaceAll(query, " AND ", " && ")
	query = strings.ReplaceAll(query, "]AND[", "] && [")
	query = strings.ReplaceAll(query, "]&&[", "] && [")
	query = strings.ReplaceAll(query, "]&[", "] && [")
	query = strings.ReplaceAll(query, "] & [", "] && [")

	nods := Item{
		Query:    query,
		TypeName: inputs["type"],
		Gjson:    "Value"}

	data, err := x.rules.Evaluate(ctx, nods, inputs)
	if err != nil {
		return nil, ErrItemNotFound
	}

	item := inputs[common.INPUT_RULE]

	delete(inputs, common.INPUT_ROLE)
	delete(inputs, common.INPUT_KEY)
	delete(inputs, common.INPUT_TYPE)
	delete(inputs, common.INPUT_RULE)

	value := gjson.GetBytes(data, "Value")
	// TODO:  figure out how to handle any type output from a rule
	result := EvalResult{Type: 0, Bool: value.Bool()}
	output := EvalOutput{
		Item:   item,
		Result: result,
		Keys:   inputs,
		Error:  ""}

	return json.Marshal(output)
}
