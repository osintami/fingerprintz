// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/etlr/utils"
)

type Toolbox struct {
	Network    utils.INetworking
	FileSystem utils.IFileSystem
	CSV        utils.ICSV
	Secrets    common.ISecrets
	Items      map[string]Item
}

type Item struct {
	Item        string
	Enabled     bool
	GJSON       string
	Description string
	Type        string
}

type Source struct {
	Name       string
	Enabled    bool
	URL        string `json:",omitempty"`
	File       string `json:",omitempty"`
	ApiKey     string `json:",omitempty"`
	InputType  string
	OutputType string
	Separator  string
}

type IWriter interface {
	Type() string
	Create(fileName string) error
	Load(job IETLJob) error
	Insert(key, value interface{}) error
}
