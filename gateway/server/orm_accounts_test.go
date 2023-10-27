// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createAccount(t *testing.T, mock sqlmock.Sqlmock, accounts *Accounts, account *Account) {
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "accounts" ("created_at","updated_at","deleted_at","name","email","ip_addr","api_key","tokens","role","last_payment","enabled","stripe_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			account.Name,
			account.Email,
			account.IpAddr,
			sqlmock.AnyArg(),
			account.Tokens,
			account.Role,
			account.LastPayment,
			account.Enabled,
			account.StripeId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(account.ID))
	mock.ExpectCommit()

	err := accounts.CreateAccount(context.TODO(), account)
	assert.Nil(t, err)
	assert.Equal(t, "test-name", account.Name)
}

func setupAccount(db *gorm.DB) (*Accounts, *Account) {
	accounts := NewAccounts(db, NewMockSender(false), "../welcome.template")
	account := &Account{
		Name:        "test-name",
		Email:       "test@example.com",
		IpAddr:      "1.2.3.4",
		ApiKey:      "",
		Tokens:      0,
		Role:        "admin",
		LastPayment: time.Now(),
		Enabled:     true,
		StripeId:    "test-stripe-id",
	}

	return accounts, account
}

func TestAccountCreate(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)
	createAccount(t, mock, accounts, account)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.True(t, account.IsAdmin())
}

func TestAccountCreateDuplicate(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)
	createAccount(t, mock, accounts, account)
	assert.Nil(t, mock.ExpectationsWereMet())

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "accounts" ("created_at","updated_at","deleted_at","name","email","ip_addr","api_key","tokens","role","last_payment","enabled","stripe_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			account.Name,
			account.Email,
			account.IpAddr,
			sqlmock.AnyArg(),
			account.Tokens,
			account.Role,
			account.LastPayment,
			account.Enabled,
			account.StripeId).
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	err := accounts.CreateAccount(context.TODO(), account)
	assert.NotNil(t, err)
}

func TestAccountFindByIpAddr(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// find the account
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "ip_addr" = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`)).
		WithArgs(account.IpAddr).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "email", "ip_addr", "api_key", "tokens", "role", "last_payment", "enabled", "stripe_id"}).
			AddRow(0, nil, nil, nil, account.Name, account.Email, account.IpAddr, account.ApiKey, account.Tokens, account.Role, account.LastPayment, account.Enabled, account.StripeId))

	account, err := accounts.FindByIpAddr(context.TODO(), account.IpAddr)
	assert.Nil(t, err)
	assert.NotNil(t, account)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAccountFindByIpAddrDoesNotExist(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// find the account
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "ip_addr" = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`)).
		WithArgs(account.IpAddr).
		WillReturnError(gorm.ErrRecordNotFound)

	account, err := accounts.FindByIpAddr(context.TODO(), account.IpAddr)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, account)
}

func TestAccountFindByEmail(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// find the account
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "email" = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`)).
		WithArgs(account.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "email", "api_key", "tokens", "role", "last_payment", "enabled", "stripe_id"}).
			AddRow(0, nil, nil, nil, account.Name, account.Email, account.ApiKey, account.Tokens, account.Role, account.LastPayment, account.Enabled, account.StripeId))

	account, err := accounts.FindByEmail(context.TODO(), account.Email)
	assert.Nil(t, err)
	assert.NotNil(t, account)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAccountFindByEmailDoesNotExist(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// find the account
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "email" = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`)).
		WithArgs(account.Email).
		WillReturnError(gorm.ErrRecordNotFound)

	account, err := accounts.FindByEmail(context.TODO(), account.Email)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, account)
}

func TestAccountFindByApiKey(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// find the account
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "api_key" = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`)).
		WithArgs(account.ApiKey).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "email", "api_key", "tokens", "role", "last_payment", "enabled", "stripe_id"}).
			AddRow(0, nil, nil, nil, account.Name, account.Email, account.ApiKey, account.Tokens, account.Role, account.LastPayment, account.Enabled, account.StripeId))

	account, err := accounts.FindByApiKey(context.TODO(), account.ApiKey)
	assert.Nil(t, err)
	assert.NotNil(t, account)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAccountFindByApiKeyDoesNotExist(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// find the account
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "api_key" = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`)).
		WithArgs(account.ApiKey).
		WillReturnError(gorm.ErrRecordNotFound)

	account, err := accounts.FindByApiKey(context.TODO(), account.ApiKey)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, account)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAccountFindByStripeId(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// find the account
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "stripe_id" = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`)).
		WithArgs(account.StripeId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "email", "api_key", "tokens", "role", "last_payment", "enabled", "stripe_id"}).
			AddRow(0, nil, nil, nil, account.Name, account.Email, account.ApiKey, account.Tokens, account.Role, account.LastPayment, account.Enabled, account.StripeId))

	account, err := accounts.FindByStripeId(context.TODO(), account.StripeId)
	assert.Nil(t, err)
	assert.NotNil(t, account)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAccountFindByStripeIdDoesNotExist(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// find the account
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "stripe_id" = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`)).
		WithArgs(account.StripeId).
		WillReturnError(gorm.ErrRecordNotFound)

	account, err := accounts.FindByStripeId(context.TODO(), account.StripeId)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, account)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAccountEnableOrDisable(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// disable the account
	account.Enabled = false
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "accounts" SET "enabled"=$1,"updated_at"=$2 WHERE "api_key" = $3 AND "accounts"."deleted_at" IS NULL`)).
		WithArgs(
			account.Enabled,
			sqlmock.AnyArg(),
			account.ApiKey).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := accounts.EnableOrDisableAccount(context.TODO(), account)
	assert.Nil(t, err)
	assert.False(t, account.Enabled)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAccountUpdateLastPayment(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// update payment
	account.LastPayment = time.Now()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "accounts" SET "last_payment"=$1,"updated_at"=$2 WHERE "api_key" = $3 AND "accounts"."deleted_at" IS NULL`)).
		WithArgs(
			account.LastPayment,
			sqlmock.AnyArg(),
			account.ApiKey).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := accounts.UpdateLastPayment(context.TODO(), account)
	assert.Nil(t, err)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAccountUpdateLastPaymentDoesNotExist(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// update payment
	account.LastPayment = time.Now()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "accounts" SET "last_payment"=$1,"updated_at"=$2 WHERE "api_key" = $3 AND "accounts"."deleted_at" IS NULL`)).
		WithArgs(
			account.LastPayment,
			sqlmock.AnyArg(),
			account.ApiKey).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	err := accounts.UpdateLastPayment(context.TODO(), account)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAccountBurnToken(t *testing.T) {
	db, mock := common.CreateMockDatabase()
	accounts, account := setupAccount(db)

	// create an account
	createAccount(t, mock, accounts, account)

	// burn an API call
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "accounts" SET "tokens"=$1,"updated_at"=$2 WHERE "api_key" = $3 AND "accounts"."deleted_at" IS NULL`)).
		WithArgs(
			account.Tokens+1,
			sqlmock.AnyArg(),
			account.ApiKey).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := accounts.BurnToken(context.TODO(), account)
	assert.Nil(t, err)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestWelcomeEmail(t *testing.T) {
	accounts, account := setupAccount(nil)
	err := accounts.WelcomeEmail(account)
	assert.Nil(t, err)
}

func TestWelcomeEmailUserFailedToSend(t *testing.T) {
	accounts := NewAccounts(nil, NewMockSender(true), "../welcome.template")
	account := &Account{
		Name:        "test-name",
		Email:       "test@example.com",
		ApiKey:      "",
		Tokens:      0,
		Role:        "admin",
		LastPayment: time.Now(),
		Enabled:     true,
		StripeId:    "test-stripe-id",
	}
	err := accounts.WelcomeEmail(account)
	assert.NotNil(t, err)
}
