// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestDataResponse(t *testing.T) {
	dr := NewDataResponse()
	inputs := common.DataInputs{}
	inputs["ip"] = "1.2.3.4"
	dataName := "ip/nope/blacklist.isBlacklisted"

	response := dr.EmptyResponse(common.Null, dataName, inputs, ErrItemNotFound)
	assert.Equal(t, "", response.Result.Raw)
	assert.Equal(t, ErrItemNotFound.Error(), response.Error)
	assert.Equal(t, dataName, dataName, response.Item)
	assert.Equal(t, common.Null, response.Result.Type)

	response = dr.EmptyResponse(common.Boolean, dataName, inputs, ErrItemNotFound)
	assert.Equal(t, "false", response.Result.Raw)

	response = dr.EmptyResponse(common.Integer, dataName, inputs, ErrItemNotFound)
	assert.Equal(t, "-1", response.Result.Raw)

	response = dr.EmptyResponse(common.Float, dataName, inputs, ErrItemNotFound)
	assert.Equal(t, "0.00", response.Result.Raw)

	response = dr.EmptyResponse(common.String, dataName, inputs, ErrItemNotFound)
	assert.Equal(t, "", response.Result.Raw)

	response = dr.EmptyResponse(common.JSON, dataName, inputs, ErrItemNotFound)
	assert.Equal(t, "{}", response.Result.Raw)

	response = dr.EmptyResponse(common.Date, dataName, inputs, ErrItemNotFound)
	assert.Equal(t, "", response.Result.Raw)
}
