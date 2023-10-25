// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type IPixels interface {
	PixelFire(context.Context, *Pixel) error
}

type Pixel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	IpAddr    string
	CookieID  string `gorm:"index:idx_pixel_cookie_id"`
	UserAgent string
	Referrer  string
	Count     int
}

type Pixels struct {
	gorm *gorm.DB
}

func NewPixels(gorm *gorm.DB) *Pixels {
	return &Pixels{gorm: gorm}
}

func (x *Pixels) PixelFire(ctx context.Context, pixel *Pixel) error {
	if x.gorm.WithContext(ctx).Model(pixel).Where("cookie_id", pixel.CookieID).UpdateColumn("count", gorm.Expr("count + ?", 1)).RowsAffected == 0 {
		return x.gorm.WithContext(ctx).Create(pixel).Error
	}
	return nil
}
