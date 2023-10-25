// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"github.com/go-resty/resty/v2"

	"github.com/osintami/fingerprintz/common"
)

type Toolbox struct {
	Client   *resty.Client
	Cache    common.IFastCache
	Watcher  common.IFileWatcher
	Schema   IDataSchema
	Secrets  common.ISecrets
	DataPath string
}
