package etl

import "errors"

var ErrBadSourceData = errors.New("bad source data")
var ErrVendorNotFound = errors.New("vendor not found")
var ErrVendorDisabled = errors.New("vendor is disabled")
var ErrEmptySchema = errors.New("empty schema")
