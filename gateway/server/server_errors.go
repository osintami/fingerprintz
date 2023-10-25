// Copyright © 2023 OSINTAMI. This is not yours.
package server

import "errors"

var ErrNotAuthorized = errors.New("not authorized")
var ErrNotFound = errors.New("route not found")
var ErrBadRequest = errors.New("bad request")
var ErrNeedToPay = errors.New("payment required")
var ErrApiKey = errors.New("invalid API key")
var ErrCookieId = errors.New("invalid cookie id")
