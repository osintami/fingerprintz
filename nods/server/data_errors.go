package server

import "errors"

var ErrNotImplemented = errors.New("not implemented")
var ErrBadData = errors.New("bad data")
var ErrMissingInputs = errors.New("missing input")
var ErrItemNotFound = errors.New("item not found")
var ErrCategoryNotFound = errors.New("category not found")
var ErrStatusNotOk = errors.New("status != http.StatusOK")
var ErrSourceNotFound = errors.New("vendor not found")
var ErrBadItemNomenclature = errors.New("malformed item {category}.{vendor}.{name}")
var ErrKeyMissing = errors.New("item key missing ip, email or phone")
