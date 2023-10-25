// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/tidwall/gjson"
)

type APIProvider struct {
	remote *API
	tools  *Toolbox
}

const (
	// NOTE:  these need to match what's in {vendor}.json
	GREYNOISE = "greynoise"
	APIVOID   = "apivoid"
)

func NewAPIProvider(tools *Toolbox, byod *API) *APIProvider {
	return &APIProvider{
		tools:  tools,
		remote: byod,
	}
}

func (x *APIProvider) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	client := x.tools.Client
	req := client.R()
	url := x.remote.URL

	if x.remote.Auth != nil {
		req.SetBasicAuth(x.tools.Secrets.Find(x.remote.Auth.Name), x.tools.Secrets.Find(x.remote.Auth.Pass))
	}

	for _, header := range x.remote.Headers {
		req = req.SetHeader(header.Key, header.Value)
	}

	for k, v := range inputs {
		url = strings.ReplaceAll(url, "{"+k+"}", v)
	}

	if x.remote.ApiKey != nil {
		apikey := x.tools.Secrets.Find(*x.remote.ApiKey)
		url = strings.ReplaceAll(url, "{key}", apikey)
	}

	resp, err := req.Get(url)

	switch x.remote.Name {
	case GREYNOISE:
		if resp.StatusCode() == http.StatusNotFound {
			return resp.Body(), nil
		}
	case APIVOID:
		errorString := gjson.GetBytes(resp.Body(), "error").String()
		if errorString != "" {
			return nil, common.ErrNoDataPresent
		}
	default:
		if resp.StatusCode() != http.StatusOK || err != nil {
			log.Error().Err(err).Int("status", resp.StatusCode()).Str("component", "api item provider").Str("source", x.remote.Name).Msg("request failed")
			return nil, common.ErrNoDataPresent
		}
	}

	return resp.Body(), err
}

func (x *APIProvider) IsCached() bool {
	return true
}
