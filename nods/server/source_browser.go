// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/ua-parser/uap-go/uaparser"
)

type BrowserSource struct {
	ua *uaparser.Parser
}

func NewBrowserSource(tools *Toolbox) IDataSource {
	parser, err := uaparser.New(tools.DataPath + "useragent.yaml")
	if err != nil {
		log.Error().Err(err).Msg("user agent regexes missing or corrupt")
		return nil
	}
	return NewDataInstance(tools, SOURCE_OSINTAMI_NAME, &BrowserSource{ua: parser})
}

func (x *BrowserSource) IsCached() bool {
	return false
}

func (x *BrowserSource) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	switch categoryName {
	case CATEGORY_BROWSER:
		return x.getUserAgentInfo(inputs[categoryName])
	}

	return nil, ErrNotImplemented
}

func (x *BrowserSource) getUserAgentInfo(ua string) (json.RawMessage, error) {
	return json.Marshal(x.ua.Parse(strings.ReplaceAll(ua, "%20", " ")))
}
