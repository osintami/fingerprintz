// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"strconv"
	"strings"

	"github.com/osintami/fingerprintz/common"
	"github.com/tidwall/gjson"
)

type DataResponse struct {
	DefaultBool  bool
	DefaultFloat float64
	DefaultInt   int64
	DefaultDate  string
	DefaultStr   string
	DefaultJSON  string
}

func NewDataResponse() *DataResponse {
	return &DataResponse{DefaultBool: false, DefaultFloat: 0.0, DefaultInt: -1, DefaultDate: "0000-00-00 00:00:00", DefaultStr: "", DefaultJSON: "{}"}
}

func (x *DataResponse) EmptyResponse(dataType common.DataType, dataName string, inputs common.DataInputs, err error) *common.DataOutput {
	switch dataType {
	case common.Boolean:
		return &common.DataOutput{
			Item:   dataName,
			Keys:   inputs,
			Result: x.MarshalResult(gjson.Result{Type: gjson.False, Raw: strconv.FormatBool(x.DefaultBool)}, dataType),
			Error:  err.Error(),
		}
	case common.Integer:
		return &common.DataOutput{
			Item:   dataName,
			Keys:   inputs,
			Result: x.MarshalResult(gjson.Result{Type: gjson.Number, Raw: strconv.FormatInt(x.DefaultInt, 10), Num: float64(x.DefaultInt)}, dataType),
			Error:  err.Error(),
		}
	case common.Float:
		return &common.DataOutput{
			Item:   dataName,
			Keys:   inputs,
			Result: x.MarshalResult(gjson.Result{Type: gjson.Number, Raw: strconv.FormatFloat(x.DefaultFloat, 'f', 2, 32), Num: x.DefaultFloat}, dataType),
			Error:  err.Error(),
		}
	case common.String:
		return &common.DataOutput{
			Item:   dataName,
			Keys:   inputs,
			Result: x.MarshalResult(gjson.Result{Type: gjson.String, Raw: x.DefaultStr, Str: x.DefaultStr}, dataType),
			Error:  err.Error(),
		}
	case common.Date:
		return &common.DataOutput{
			Item:   dataName,
			Keys:   inputs,
			Result: x.MarshalResult(gjson.Result{Type: gjson.Null, Raw: "", Str: ""}, dataType),
			Error:  err.Error(),
		}
	case common.JSON:
		return &common.DataOutput{
			Item:   dataName,
			Keys:   inputs,
			Result: x.MarshalResult(gjson.Result{Type: gjson.String, Raw: x.DefaultJSON, Str: x.DefaultJSON}, dataType),
			Error:  err.Error(),
		}
	default:
		return &common.DataOutput{
			Item:   dataName,
			Keys:   inputs,
			Result: x.MarshalResult(gjson.Result{Type: gjson.Null, Raw: ""}, dataType),
			Error:  err.Error(),
		}
	}
}

func (x *DataResponse) MarshalResult(out gjson.Result, itemType common.DataType) common.DataResult {
	result := common.DataResult{}
	result.Raw = strings.TrimLeft(out.Raw, "\"")
	result.Raw = strings.TrimRight(result.Raw, "\"")
	result.Type = itemType
	switch result.Type {
	case common.String:
		result.Str = &out.Str
	case common.Date:
		result.Str = &out.Str
	case common.JSON:
		result.Str = &out.Raw
	case common.Integer:
		i := out.Float()
		result.Num = &i
	case common.Float:
		i := out.Float()
		result.Num = &i
	case common.Boolean:
		b := out.Bool()
		result.Bool = &b
	}
	return result
}
