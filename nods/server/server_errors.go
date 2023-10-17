// Copyright © 2023 OSINTAMI. This is not yours.
package server

import "errors"

var ErrInvalidItemParam = errors.New("data item parameter missing or invalid")
var ErrInvalidSourceParam = errors.New("source parameter missing or invalid")
var ErrInvalidCategoryParam = errors.New("category parameter missing or invalid")
var ErrInvalidCategoryInputParam = errors.New("path parameter category does not match query parameter name")
var ErrInvalidUserRole = errors.New("role parameter missing or invalid")
