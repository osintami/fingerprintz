// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mcnijman/go-emailaddress"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"gorm.io/gorm"
)

type IAccounts interface {
	CreateAccount(ctx context.Context, user *Account) error
	FindByEmail(ctx context.Context, email string) (*Account, error)
	FindByIpAddr(ctx context.Context, ipaddr string) (*Account, error)
	FindByApiKey(ctx context.Context, apikey string) (*Account, error)
	FindByStripeId(ctx context.Context, apikey string) (*Account, error)
	EnableOrDisableAccount(ctx context.Context, account *Account) error
	UpdateLastPayment(ctx context.Context, account *Account) error
	BurnToken(ctx context.Context, account *Account) error
	WelcomeEmail(user *Account) error
}

type Account struct {
	gorm.Model  `json:"-"`
	Name        string
	Email       string `gorm:"index:idx_account_email_ipaddr, unique"`
	IpAddr      string `gorm:"index:idx_account_email_ipaddr, unique"`
	ApiKey      string
	Tokens      int
	Role        string
	LastPayment time.Time
	Enabled     bool
	StripeId    string
	Calls       []Call `json:"-"`
}

func (x Account) IsAdmin() bool {
	return x.Role == ROLE_ADMIN
}

const (
	ROLE_ADMIN = "admin"
	ROLE_USER  = "user"
	ROLE_TRIAL = "trial"
)

var ErrUserExists = errors.New("user exists")

type Accounts struct {
	gorm         *gorm.DB
	mail         *Email
	templateFile string
}

func NewAccounts(gorm *gorm.DB, sender IMail, templateFile string) *Accounts {
	return &Accounts{
		gorm:         gorm,
		mail:         NewEmail(sender),
		templateFile: templateFile}
}

func (x *Accounts) CreateAccount(ctx context.Context, account *Account) error {
	if account.ApiKey == "" {
		key := make([]byte, 32)
		rand.Read(key)
		account.ApiKey = base64.StdEncoding.EncodeToString([]byte(key))
	}
	if err := x.gorm.WithContext(ctx).Create(account).Error; err != nil {
		log.Error().Err(err).Str("component", "accounts").Str("email", account.Email).Msg("user exists")
		return ErrUserExists
	}
	return nil
}

func (x *Accounts) FindByStripeId(ctx context.Context, stripeId string) (*Account, error) {
	var account Account
	err := x.gorm.WithContext(ctx).First(&account, "stripe_id", stripeId).Error
	if err != nil {
		log.Error().Err(err).Str("component", "accounts").Str("find", stripeId).Msg("stripe doesn't exist")
		return nil, err
	}
	return &account, nil
}

func (x *Accounts) FindByEmail(ctx context.Context, email string) (*Account, error) {
	var account Account
	err := x.gorm.WithContext(ctx).First(&account, "email", email).Error
	if err != nil {
		log.Error().Err(err).Str("component", "accounts").Str("find", email).Msg("email doesn't exist")
		return nil, err
	}
	return &account, nil
}

func (x *Accounts) FindByIpAddr(ctx context.Context, ipaddr string) (*Account, error) {
	var account Account
	err := x.gorm.WithContext(ctx).First(&account, "ip_addr", ipaddr).Error
	if err != nil {
		log.Error().Err(err).Str("component", "accounts").Str("find", ipaddr).Msg("ip addr doesn't exist")
		return nil, err
	}
	return &account, nil
}

func (x *Accounts) FindByApiKey(ctx context.Context, apikey string) (*Account, error) {
	var account Account
	err := x.gorm.WithContext(ctx).First(&account, "api_key", apikey).Error
	if err != nil {
		log.Error().Err(err).Str("component", "accounts").Str("find", apikey).Msg("apikey doesn't exist")
		return nil, err
	}
	return &account, nil
}

func (x *Accounts) EnableOrDisableAccount(ctx context.Context, account *Account) error {
	return x.gorm.WithContext(ctx).Model(account).Where("api_key", account.ApiKey).Update("enabled", account.Enabled).Error
}

func (x *Accounts) UpdateLastPayment(ctx context.Context, account *Account) error {
	err := x.gorm.WithContext(ctx).Model(account).Where("api_key", account.ApiKey).Update("last_payment", account.LastPayment).Error
	if err != nil {
		log.Error().Err(err).Str("component", "accounts").Str("apikey", account.ApiKey).Msg("api key doesn't exist")
	}
	return err
}

// func (x *Accounts) AddTokens(ctx context.Context, account *Account) error {
// 	err := x.gorm.WithContext(ctx).Model(&account).Where("api_key = ?", account.ApiKey).Update("tokens", account.Tokens).Error
// 	if err != nil {
// 		log.Error().Err(err).Str("component", "accounts").Str("apikey", account.ApiKey).Msg("failed to increment tokens")
// 		return err
// 	}
// 	return nil
// }

func (x *Accounts) BurnToken(ctx context.Context, account *Account) error {
	account.Tokens += 1
	return x.gorm.WithContext(ctx).Model(account).Where("api_key", account.ApiKey).Update("tokens", account.Tokens).Error
}

func (x *Accounts) WelcomeEmail(account *Account) error {
	recipients := []string{}
	recipients = append(recipients, account.Email)

	receipt := time.Now().Format(common.GO_DEFAULT_DATE)
	subject := fmt.Sprintf("OSINTAMI Receipt - %s", receipt)

	body, err := x.mail.loadTemplate(x.templateFile)

	if err == nil {
		content := string(body)
		content = strings.ReplaceAll(content, "$token", account.ApiKey)
		content = strings.ReplaceAll(content, "$email", account.Email)
		content = strings.ReplaceAll(content, "$name", account.Name)

		for _, recipient := range recipients {
			buyer, _ := emailaddress.Parse(recipient)
			internal, _ := emailaddress.Parse("sales@osintami.com")
			//fmt.Println("[EMAIL]", email.String()+"<br>", "\t[SUBJECT] "+subject+"<br>", "\t[BODY]<br>"+body)
			err = x.mail.Send(subject, content, internal, true)
			if err != nil {
				log.Error().Err(err).Str("component", "welcome").Str("email", internal.String()).Msg("email send failed")
			}
			err = x.mail.Send(subject, content, buyer, true)
			if err != nil {
				log.Error().Err(err).Str("component", "welcome").Str("email", buyer.String()).Msg("email send failed")
			}
		}
	}
	return err
}
