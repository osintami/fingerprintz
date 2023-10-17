// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

func init() {
	os.Remove("./test/test.log")
	log.InitLogger("./test/", "test.log", "ERROR", false)
}

func mockToolbox() *Toolbox {
	return &Toolbox{
		Secrets:  NewMockSecretsManager(),
		Cache:    NewMockCache(),
		Watcher:  NewMockWatcher(),
		Schema:   NewMockDataSchema(),
		DataPath: ""}
}

func buildItemRequest(pathParams map[string]string, queryParams map[string]string) *http.Request {
	path := fmt.Sprintf("%s/%s/%s/%s", "/v1/data", pathParams["category"], pathParams["vendor"], pathParams["item"])
	return common.BuildRequest(http.MethodGet, path, pathParams, queryParams)
}

func nodsServer(dataFail bool) *NormalizedDataServer {
	schema := NewMockDataSchema()
	router := NewMockDataRouter(dataFail)
	// TODO:  we have good test coverage on the server with this real object, need to mock it and recheck
	rules := NewRuleProvider(router, schema)
	return &NormalizedDataServer{
		router:  router,
		schema:  schema,
		secrets: NewMockSecretsManager(),
		params:  common.NewParameterHelper(),
		rules:   NewComplexProvider(rules)}
}
