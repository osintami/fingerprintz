// Copyright © 2023 OSINTAMI. This is not yours.
package server

import "errors"

var ErrMissingInputs = errors.New("missing input parameters")
var ErrFingerprintSmudged = errors.New("fingerprint smudged")
