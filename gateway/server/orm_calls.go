// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ICalls interface {
	Call(ctx context.Context, call *Call) error
}

type Call struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	AccountId uint `gorm:"index:idx_call_key_id"`
	IpAddr    string
	RequestId string
	API       string
	Page      string
}

type Calls struct {
	gorm *gorm.DB
}

func NewCalls(gorm *gorm.DB) *Calls {
	return &Calls{gorm: gorm}
}

func (x *Calls) Call(ctx context.Context, call *Call) error {
	return x.gorm.WithContext(ctx).Create(&call).Error
}
