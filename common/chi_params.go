// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	CATEGORY_BROWSER = "browser"
)

type IParameterHelper interface {
	Keys(r *http.Request) (map[string]string, error)
}

type ParameterHelper struct {
	// uaregex *regexp.Regexp
}

func NewParameterHelper() IParameterHelper {
	// uaregex := regexp.MustCompile(`browser=([^&#]*)`)
	return &ParameterHelper{}
}

func (x *ParameterHelper) Keys(r *http.Request) (map[string]string, error) {
	keys := make(map[string]string)
	if r.Method == "GET" {
		params := r.URL.Query()
		// TODO:  find a bad browser agent string so we can prove this is still required
		// 	useragent := x.uaregex.Find([]byte(r.URL.RawQuery))
		// 	// NOTE:  we need handle disallowed HTTP param characters found in some user-agent strings
		// 	if useragent != nil {
		// 		keys[CATEGORY_BROWSER] = strings.TrimLeft(string(useragent), "browser=")
		// 	}
		// }
		for k := range params {
			keys[k] = strings.TrimSpace(params.Get(k))
		}
	} else {
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&keys)
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}
