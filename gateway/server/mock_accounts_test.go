// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MockAccounts struct {
	table map[string]bool
}

func NewMockAccounts(table map[string]bool) IAccounts {
	return &MockAccounts{table: table}
}

func (x *MockAccounts) Name() string {
	return "account"
}

func (x *MockAccounts) CreateAccount(ctx context.Context, user *Account) error {
	if x.table["CreateAccount"] == true {
		return gorm.ErrRecordNotFound
	}
	user.ApiKey = uuid.NewString()
	return nil
}

func (x *MockAccounts) FindByEmail(ctx context.Context, email string) (*Account, error) {
	if x.table["FindByEmail"] == true {
		return nil, gorm.ErrRecordNotFound
	}
	var account Account
	if email == "test@example.com" {
		account.ApiKey = "test-api-key"
		account.Email = email
		account.Role = "user"
		account.Enabled = true
		account.Tokens = 0
		return &account, nil
	} else {
		return nil, gorm.ErrRecordNotFound
	}
}

func (x *MockAccounts) FindByStripeId(ctx context.Context, email string) (*Account, error) {
	if x.table["FindByStripeId"] == true {
		return nil, gorm.ErrRecordNotFound
	}
	var account Account
	account.ApiKey = "stripe-api-key"
	account.Email = email
	account.Role = "user"
	account.Enabled = true
	account.Tokens = 0
	return &account, nil
}

func (x *MockAccounts) FindByApiKey(ctx context.Context, apikey string) (*Account, error) {
	if x.table["FindByApiKey"] == true {
		return nil, gorm.ErrRecordNotFound
	}
	var account Account
	if apikey != "admin_api_key" && apikey != "user_api_key" {
		return nil, errors.New("user not found")
	}
	account.ApiKey = apikey
	// HACK:  for testing use "user" or "admin" for the apikey
	if apikey == "admin_api_key" {
		account.Role = "admin"
	} else {
		account.Role = "user"
	}
	account.Enabled = true
	account.Tokens = 0
	return &account, nil
}

func (x *MockAccounts) EnableOrDisableAccount(ctx context.Context, account *Account) error {
	if x.table["EnableOrDisableAccount"] == true {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (x *MockAccounts) UpdateLastPayment(ctx context.Context, account *Account) error {
	if x.table["UpdateLastPayment"] == true {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (x *MockAccounts) BurnToken(ctx context.Context, account *Account) error {
	if x.table["BurnToken"] == true {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (x *MockAccounts) ReadEmailTemplate(fileName string) (string, error) {
	if x.table["ReadEmailTemplate"] == true {
		return "", gorm.ErrRecordNotFound
	}
	return "", nil
}

func (x *MockAccounts) WelcomeEmail(user *Account) error {
	if x.table["WelcomeEmail"] == true {
		return gorm.ErrRecordNotFound
	}
	return nil
}
