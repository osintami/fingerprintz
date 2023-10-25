// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestHandlerItem(t *testing.T) {
	server := nodsServer(false)

	pParams := make(map[string]string)

	pParams["category"] = "ip"
	pParams["vendor"] = "ipsum"
	pParams["item"] = "blacklist.isBlacklisted"

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	r := buildItemRequest(pParams, qParams)
	w := httptest.NewRecorder()

	server.GetItemHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	content := "{\n    \"Item\": \"ip/ipsum/blacklist.isBlacklisted\",\n    \"Result\": {\n        \"Type\": 1,\n        \"Bool\": true\n    },\n    \"Keys\": {\n        \"ip\": \"1.2.3.4\"\n    },\n    \"Error\": \"\"\n}\n"
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerItemStripKeys(t *testing.T) {
	server := nodsServer(false)

	pParams := make(map[string]string)

	pParams["category"] = "ip"
	pParams["vendor"] = "ipsum"
	pParams["item"] = "blacklist.isBlacklisted"

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"
	// add some reserved keys and make sure they aren't returned
	qParams["role"] = "test-role"
	qParams["key"] = "test-api-key"
	qParams["type"] = "test-type"
	qParams["rule"] = "test-rule"
	// add an extra non-reserved key
	qParams["nope"] = "yep"

	r := buildItemRequest(pParams, qParams)
	w := httptest.NewRecorder()

	server.GetItemHandler(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	content := "{\n    \"Item\": \"ip/ipsum/blacklist.isBlacklisted\",\n    \"Result\": {\n        \"Type\": 1,\n        \"Bool\": true\n    },\n    \"Keys\": {\n        \"ip\": \"1.2.3.4\",\n        \"nope\": \"yep\"\n    },\n    \"Error\": \"\"\n}\n"
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerItemsNoData(t *testing.T) {
	// NOTE:  force a data failure
	server := nodsServer(true)

	pParams := make(map[string]string)

	pParams["category"] = "ip"
	pParams["vendor"] = "ipsum"
	pParams["item"] = "blacklist.isBlacklisted"

	qParams := make(map[string]string)
	qParams["ip"] = "0.0.0.0"

	r := buildItemRequest(pParams, qParams)
	w := httptest.NewRecorder()

	server.GetItemHandler(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
	content := common.BuildErrorResponse(ErrItemNotFound)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerItemsCategoryMismatch(t *testing.T) {
	server := nodsServer(false)

	pParams := make(map[string]string)

	pParams["category"] = "phone"
	pParams["vendor"] = "ipsum"
	pParams["item"] = "blacklist.isBlacklisted"

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	r := buildItemRequest(pParams, qParams)
	w := httptest.NewRecorder()

	server.GetItemHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrMissingInputs)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerItemsMissingKeys(t *testing.T) {
	server := nodsServer(false)

	r := buildItemRequest(nil, nil)
	w := httptest.NewRecorder()

	server.GetItemHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrMissingInputs)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerItemsInvalidCategory(t *testing.T) {
	server := nodsServer(false)

	pParams := make(map[string]string)

	pParams["category"] = "nope"
	pParams["vendor"] = "osintami"
	pParams["item"] = "isBlacklisted"

	qParams := make(map[string]string)
	qParams["nope"] = "1.2.3.4"

	r := buildItemRequest(pParams, qParams)
	w := httptest.NewRecorder()

	server.GetItemHandler(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
	content := common.BuildErrorResponse(ErrCategoryNotFound)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerItemsMissingCategory(t *testing.T) {
	server := nodsServer(false)

	pParams := make(map[string]string)

	pParams["category"] = ""
	pParams["vendor"] = "osintami"
	pParams["item"] = "isBlacklisted"

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	r := buildItemRequest(pParams, qParams)
	w := httptest.NewRecorder()

	server.GetItemHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrInvalidCategoryParam)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerItemsMissingSource(t *testing.T) {
	server := nodsServer(false)

	pParams := make(map[string]string)

	pParams["category"] = "rule"
	pParams["vendor"] = ""
	pParams["item"] = "isBlacklisted"

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	r := buildItemRequest(pParams, qParams)
	w := httptest.NewRecorder()

	server.GetItemHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrInvalidSourceParam)
	assert.Equal(t, content, w.Body.String())
}

func TestHandlerItemsMissingItem(t *testing.T) {
	server := nodsServer(false)

	pParams := make(map[string]string)

	pParams["category"] = "rule"
	pParams["vendor"] = "osintami"
	pParams["item"] = ""

	qParams := make(map[string]string)
	qParams["ip"] = "1.2.3.4"

	r := buildItemRequest(pParams, qParams)
	w := httptest.NewRecorder()

	server.GetItemHandler(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	content := common.BuildErrorResponse(ErrInvalidItemParam)
	assert.Equal(t, content, w.Body.String())
}
