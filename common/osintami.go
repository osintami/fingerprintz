// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/osintami/fingerprintz/log"
)

type WhoamiInfo struct {
	IpAddr    string
	Tor       bool
	CloudNode bool
	Proxy     bool
	VPN       bool
	Bot       bool
	Blacklist bool
	Latitude  float64
	Longitude float64
	City      string
	Country   string
	Continent string
	LastSeen  string
}

type DataType int

const (
	Null DataType = iota
	Boolean
	Float
	Integer
	String
	Date
	JSON
)

func (t DataType) String() string {
	switch t {
	default:
		return ""
	case Null:
		return "Null"
	case Boolean:
		return "Boolean"
	case Float:
		return "Float"
	case Integer:
		return "Integer"
	case String:
		return "String"
	case Date:
		return "Date"
	case JSON:
		return "JSON"
	}
}

func (t DataType) ToDataType(itemType string) DataType {
	switch itemType {
	case "Null":
		return Null
	case "Boolean":
		return Boolean
	case "String":
		return String
	case "Integer":
		return Integer
	case "Float":
		return Float
	case "Date":
		return Date
	case "JSON":
		return JSON
	}
	return Null
}

type DataResult struct {
	Type DataType
	Raw  string   `json:"-"`
	Str  *string  `json:",omitempty"`
	Num  *float64 `json:",omitempty"`
	Bool *bool    `json:",omitempty"`
}

func (x DataResult) IsEmpty() bool {
	return x.Raw == "" || x.Raw == "{}"
}

func (x DataResult) Boolean() bool {
	if x.Bool != nil {
		return *x.Bool
	}
	if x.Raw == "" {
		return false
	}
	return x.Raw == "true"
}

const (
	INPUT_TYPE = "type"
	INPUT_ROLE = "role"
	INPUT_KEY  = "key"
	INPUT_RULE = "rule"
)

type DataInputs map[string]string

func (x DataInputs) String() string {
	return x.Params()
}

func (x DataInputs) Params() string {
	params := "?"

	keys := make([]string, 0, len(x))
	for k := range x {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if key == INPUT_RULE {
			continue
		}
		params += key + "=" + x[key] + "&"
	}
	params = strings.TrimRight(params, "&")
	return params
}

type DataOutput struct {
	Item   string
	Result DataResult
	Keys   DataInputs
	Error  string
}

type INods interface {
	Whoami(ip string) (*WhoamiInfo, error)
	Item(name string, keys map[string]string) (*DataOutput, error)
}

type OSINTAMIClient struct {
	client  *resty.Client
	nodsURL string
}

func NewOSINTAMIClient(client *resty.Client, nodsURL string) *OSINTAMIClient {
	return &OSINTAMIClient{client: client, nodsURL: nodsURL}
}

func (x *OSINTAMIClient) Whoami(ip string) (*WhoamiInfo, error) {
	resp, err := x.client.R().Get(x.nodsURL + "/v1/data/whoami?ip=" + ip)
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Error().Err(err).Str("component", "osintami").Int("status", resp.StatusCode()).Msg("osintami not reachable")
		return &WhoamiInfo{}, ErrItemNotFound
	}
	defer resp.RawBody().Close()
	info := &WhoamiInfo{}
	err = x.parseBody(resp.Body(), info)
	return info, err
}

// func (x *OSINTAMIClient) Fingerprint(keys map[string]string) (string, error) {
// 	inputs, _ := json.Marshal(keys)
// 	resp, err := x.client.R().SetBody(inputs).Post(x.nodsURL + "/v1/data/fingerprint")
// 	if err != nil || resp.StatusCode() != http.StatusOK {
// 		log.Error().Err(err).Str("component", "osintami").Int("status", resp.StatusCode()).Msg("osintami not reachable")
// 		return "", ErrItemNotFound
// 	}
// 	defer resp.RawBody().Close()
// 	return string(resp.Body()), nil
// }

func (x *OSINTAMIClient) Item(name string, keys map[string]string) (*DataOutput, error) {
	resp, err := x.client.R().SetQueryParams(keys).Get(x.nodsURL + "/data/" + name)
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Error().Err(err).Str("component", "osintami").Int("status", resp.StatusCode()).Msg("osintami not reachable")
		return &DataOutput{}, ErrItemNotFound
	}
	response := &DataOutput{}
	err = x.parseBody(resp.Body(), response)
	return response, err
}

func (x *OSINTAMIClient) parseBody(data []byte, response interface{}) error {
	err := json.Unmarshal(data, response)
	if err != nil {
		log.Error().Err(err).Str("component", "osintami").Msg("decode osintami response")
		return err
	}
	return nil
}
