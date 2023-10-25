// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/tidwall/sjson"
)

type IRuleProvider interface {
	IsCached() bool
	CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error)
	Evaluate(ctx context.Context, nods Item, inputs common.DataInputs) ([]byte, error)
}

type RuleProvider struct {
	router IDataRouter
	rules  map[string]Item
	dregex *regexp.Regexp
	sregex *regexp.Regexp
	dates  common.DateStuff
}

func NewRuleSource(tools *Toolbox, router IDataRouter, schema IDataSchema) IDataSource {
	return NewDataInstance(tools, SOURCE_OSINTAMI_NAME, NewRuleProvider(router, schema))
}

func NewRuleProvider(router IDataRouter, schema IDataSchema) IRuleProvider {
	return &RuleProvider{
		router: router,
		rules:  schema.ListRulesItems(),
		sregex: regexp.MustCompile(`\[(.*?)\]`),
		dregex: regexp.MustCompile(`@{(.*?)}`),
		dates:  *common.NewDateStuff()}
}

func (x *RuleProvider) IsCached() bool {
	return false
}

func (x *RuleProvider) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	if categoryName != CATEGORY_RULE {
		return []byte("{}"), ErrNotImplemented
	}

	if nods, ok := x.rules[inputs[common.INPUT_RULE]]; ok {
		output, err := x.Evaluate(ctx, nods, inputs)
		if err != nil {
			return []byte("{}"), err
		}
		return output, err
	}
	return []byte("{}"), common.ErrNoDataPresent
}

func (x *RuleProvider) Evaluate(ctx context.Context, nods Item, inputs common.DataInputs) ([]byte, error) {
	// insert dates into rule
	matches := x.dregex.FindAllStringSubmatch(nods.Query, -1)
	for i, match := range matches {
		date := x.dates.AgoStringToDate(match[1])
		nods.Query = strings.ReplaceAll(nods.Query, matches[i][0], "'"+date+"'")
	}

	// insert item values into the parameter map
	parameters := make(map[string]interface{}, 8)
	matches = x.sregex.FindAllStringSubmatch(nods.Query, -1)
	for _, match := range matches {
		// NOTE:  internal source item routing requires the item to be passed with
		//   the other parameters for some sources (ie. thing and nods sources)
		item := match[1]
		inputs[common.INPUT_RULE] = item
		result, err := x.router.DataValue(ctx, NewItemSplitter(item), inputs)
		log.Debug().Str("component", "rules engine").Str("rule", nods.Query).Str(common.INPUT_RULE, item).Str("result", result.Result.Raw).Msg("rule partial eval")
		// TODO:  this only works because we create empty result sets with smart defaults, at some point this will BITE us
		if err == nil || err == ErrSourceNotFound || err == common.ErrNoDataPresent {
			param := match[1]
			switch result.Result.Type {
			case common.Boolean:
				parameters[param] = *result.Result.Bool
			case common.Date:
				// NOTE:  greater/less operators don't like strings when parameters
				//   are used so do the substitution manually, otherwise let
				//   govaluate handle it
				nods.Query = strings.ReplaceAll(nods.Query, match[0], "'"+*result.Result.Str+"'")

			case common.String:
				parameters[param] = *result.Result.Str
			case common.Float:
				parameters[param] = *result.Result.Num
			case common.Integer:
				parameters[param] = *result.Result.Num
			case common.JSON:
				parameters[param] = *result.Result.Str
			case common.Null:
				// HACK:  really should figure out what the type was supposed to be and take the proper default, might be a string or a number
				parameters[param] = false
			}
		}
	}

	// NOTE:  this is slow (no precompile) but requires all parameters to work
	//   exactly as govaluate wants and we are bending the rules a bit with dates
	expression, err := govaluate.NewEvaluableExpression(nods.Query)
	if err != nil {
		log.Error().Err(err).Str("component", "rules engine").Str("rule", nods.Query).Msg("rule expression")
		return nil, err
	}

	// evaluate the expression, returns an interface
	govaluateResult, err := expression.Evaluate(parameters)
	if err != nil {
		log.Error().Err(err).Str("component", "rules engine").Str("rule", nods.Query).Msg("rule evaluation")
		return nil, err
	}

	// leverage sjson to handle putting the interface into proper formatted JSON
	if nods.Gjson != PASSTHROUGH {
		jsonResult, err := sjson.Set("", nods.Gjson, govaluateResult)
		return []byte(jsonResult), err
	} else {
		return json.Marshal(govaluateResult)
	}
}
