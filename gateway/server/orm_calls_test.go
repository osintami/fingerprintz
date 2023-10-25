// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
)

func TestCalls(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	calls := NewCalls(db)
	call := &Call{
		CreatedAt: time.Now(),
		AccountId: 0,
		IpAddr:    "1.2.3.4",
		RequestId: uuid.NewString(),
		API:       "nods"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "calls" ("created_at","account_id","ip_addr","request_id","api","page") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			call.AccountId,
			call.IpAddr,
			call.RequestId,
			call.API,
			call.Page).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(call.ID))
	mock.ExpectCommit()

	err := calls.Call(context.TODO(), call)
	assert.Nil(t, err)

	assert.Nil(t, mock.ExpectationsWereMet())
}
