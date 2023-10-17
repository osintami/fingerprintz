// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"

	"gorm.io/gorm"
)

type MockPixels struct {
	table map[string]bool
}

func NewMockPixels(table map[string]bool) IPixels {
	return &MockPixels{table: table}
}

func (x *MockPixels) PixelFire(ctx context.Context, pixel *Pixel) error {
	if x.table["PixelFire"] == true {
		return gorm.ErrDuplicatedKey
	}
	return nil
}
