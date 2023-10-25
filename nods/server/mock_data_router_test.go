// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"

	"github.com/osintami/fingerprintz/common"
)

type MockDataRouter struct {
	response *DataResponse
	fail     bool
}

func NewMockDataRouter(fail bool) IDataRouter {
	return &MockDataRouter{
		response: NewDataResponse(),
		fail:     fail,
	}
}

func (x *MockDataRouter) Init() {
}

func (x *MockDataRouter) DataValue(ctx context.Context, dataURI *DataURI, inputs common.DataInputs) (*common.DataOutput, error) {
	if x.fail {
		return x.response.EmptyResponse(common.Null, dataURI.Key(), inputs, ErrItemNotFound), ErrItemNotFound
	}

	// category not found
	if dataURI.CategoryName == "nope" {
		return x.response.EmptyResponse(common.Null, dataURI.ItemName, inputs, ErrCategoryNotFound), ErrCategoryNotFound
	}

	// source not found
	if dataURI.SourceName != "ipsum" && dataURI.SourceName != "uhb" && dataURI.SourceName != "osintami" && dataURI.SourceName != "maxmind" && dataURI.SourceName != "pwned" {
		return x.response.EmptyResponse(common.Null, dataURI.ItemName, inputs, ErrSourceNotFound), ErrSourceNotFound
	}

	// only mocking data for one IP address
	if inputs["ip"] != "1.2.3.4" && inputs["ip"] != "" {
		return x.response.EmptyResponse(common.Null, dataURI.ItemName, inputs, common.ErrNoDataPresent), common.ErrNoDataPresent
	}

	if dataURI.URI == "email/pwned/breachCount" {
		value := 3.0
		output := &common.DataOutput{
			Item: "email/pwned/breachCount",
			Result: common.DataResult{
				Type: common.Integer,
				Raw:  "3",
				Num:  &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	if dataURI.URI == "email/pwned/breachAgeDate" {
		value := "2023-10-01"
		output := &common.DataOutput{
			Item: "email/pwned/breachAgeDate",
			Result: common.DataResult{
				Type: common.Date,
				Raw:  value,
				Str:  &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	// items supported for mock purposes
	if dataURI.URI == "ip/ipsum/blacklist.isBlacklisted" {
		value := true
		output := &common.DataOutput{
			Item: "ip/ipsum/blacklist.isBlacklisted",
			Result: common.DataResult{
				Type: common.Boolean,
				Raw:  "true",
				Bool: &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	if dataURI.URI == "ip/maxmind/location" {
		value := "{\"city\":\"Council Bluffs\",\"continent\":\"NA\",\"country\":\"US\",\"latitude\":41.2591,\"longitude\":-95.8517}"
		output := &common.DataOutput{
			Item: "ip/maxmind/location",
			Result: common.DataResult{
				Type: common.JSON,
				Raw:  value,
				Str:  &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	if dataURI.URI == "rule/osintami/isTor" {
		value := true
		output := &common.DataOutput{
			Item: "rule/osintami/isTor",
			Result: common.DataResult{
				Type: common.Boolean,
				Raw:  "true",
				Bool: &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	if dataURI.URI == "rule/osintami/isCloudNode" {
		value := true
		output := &common.DataOutput{
			Item: "rule/osintami/isCloudNode",
			Result: common.DataResult{
				Type: common.Boolean,
				Raw:  "true",
				Bool: &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	if dataURI.URI == "rule/osintami/isProxy" {
		value := true
		output := &common.DataOutput{
			Item: "rule/osintami/isProxy",
			Result: common.DataResult{
				Type: common.Boolean,
				Raw:  "true",
				Bool: &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	if dataURI.URI == "rule/osintami/isVPN" {
		value := true
		output := &common.DataOutput{
			Item: "rule/osintami/isVPN",
			Result: common.DataResult{
				Type: common.Boolean,
				Raw:  "true",
				Bool: &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	if dataURI.URI == "rule/osintami/isBot" {
		value := true
		output := &common.DataOutput{
			Item: "rule/osintami/isBot",
			Result: common.DataResult{
				Type: common.Boolean,
				Raw:  "true",
				Bool: &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	if dataURI.URI == "rule/osintami/isBlacklisted" {
		value := true
		output := &common.DataOutput{
			Item: "rule/osintami/isBlacklisted",
			Result: common.DataResult{
				Type: common.Boolean,
				Raw:  "true",
				Bool: &value,
			},
			Keys:  inputs,
			Error: "",
		}
		return output, nil
	}

	// TODO:  add rules data items for whoami handler tests

	// item not found
	return x.response.EmptyResponse(common.Null, dataURI.Key(), inputs, ErrItemNotFound), ErrItemNotFound
}

func (x *MockDataRouter) CategoryValues(ctx context.Context, categoryName string, inputs common.DataInputs) ([]*common.DataOutput, error) {
	itemResults := []*common.DataOutput{}
	if x.fail {
		return itemResults, ErrItemNotFound
	}

	value := true
	output := &common.DataOutput{
		Item: "ip/ipsum/blacklist.isBlacklisted",
		Result: common.DataResult{
			Type: common.Boolean,
			Raw:  "true",
			Bool: &value,
		},
		Keys:  inputs,
		Error: "",
	}

	itemResults = append(itemResults, output)

	value = true
	output = &common.DataOutput{
		Item: "ip/uhb/blacklist.isBlacklisted",
		Result: common.DataResult{
			Type: common.Boolean,
			Raw:  "true",
			Bool: &value,
		},
		Keys:  inputs,
		Error: "",
	}

	itemResults = append(itemResults, output)

	if len(itemResults) == 0 {
		return itemResults, ErrItemNotFound
	}

	return itemResults, nil
}
