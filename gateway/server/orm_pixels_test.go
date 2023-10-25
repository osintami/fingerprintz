// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPixelFire(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	pixels := NewPixels(db)
	pixel := &Pixel{
		CookieID:  "test-uuid",
		IpAddr:    "1.2.3.4",
		UserAgent: "test-user-agent-string",
		Referrer:  "test-referrer",
		Count:     1,
	}

	// fail the update, the record doesn't exist yet
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "pixels" SET "count"=count + $1 WHERE "cookie_id" = $2`)).
		WithArgs(
			1,
			pixel.CookieID).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	// insert the new record
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "pixels" ("created_at","ip_addr","cookie_id","user_agent","referrer","count") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			pixel.IpAddr,
			pixel.CookieID,
			pixel.UserAgent,
			pixel.Referrer,
			1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pixel.ID))
	mock.ExpectCommit()

	err := pixels.PixelFire(context.TODO(), pixel)
	assert.Nil(t, err)

	// update the existing record
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "pixels" SET "count"=count + $1 WHERE "cookie_id" = $2`)).
		WithArgs(
			1,
			pixel.CookieID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	pixel.Count += 1
	err = pixels.PixelFire(context.TODO(), pixel)
	assert.Nil(t, err)

	assert.Nil(t, mock.ExpectationsWereMet())
}
