// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
	"github.com/osintami/fingerprintz/pwned"
)

type PwnedInfo struct {
	Email                   string
	BreachCount             int
	PastebinCount           int
	MostRecentBreachInDays  int
	MostRecentBreachInYears int
	MostRecentBreachDate    string
	BreachInfo              []*BreachInfo
}

type BreachInfo struct {
	Name  string
	Title string
	Date  string
}

type PwnedPasswordInfo struct {
	PastebinCount int
}

type PwnedSource struct {
	client pwned.IPwned
}

func NewPwnedSource(tools *Toolbox, client pwned.IPwned) IDataSource {
	return NewDataInstance(
		tools,
		SOURCE_PWNED_NAME,
		&PwnedSource{
			client: client,
		})
}

func (x *PwnedSource) IsCached() bool {
	return true
}

func (x *PwnedSource) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	switch categoryName {
	case CATEGORY_EMAIL:
		return x.GetEmailInfo(ctx, inputs[categoryName])
	case CATEGORY_PASSWORD:
		return x.GetPasswordInfo(ctx, inputs[categoryName])
	}
	return nil, ErrNotImplemented
}

func (x *PwnedSource) GetEmailInfo(ctx context.Context, email string) (json.RawMessage, error) {
	info, err := x.Breaches(email)
	if err != nil {
		data, _ := json.Marshal(info)
		return data, common.ErrNoDataPresent
	}

	count, err := x.Pastebin(email)
	if err != nil {
		data, _ := json.Marshal(info)
		return data, common.ErrNoDataPresent
	}

	info.PastebinCount = count
	return json.Marshal(info)
}

func (x *PwnedSource) GetPasswordInfo(ctx context.Context, password string) (json.RawMessage, error) {
	var info PwnedPasswordInfo
	count, err := x.client.HaveIBeenPwned(password)
	if err != nil {
		log.Error().Err(err).Str("component", "pwned").Str("password", password).Msg("pwned password")
		data, _ := json.Marshal(info)
		return data, common.ErrNoDataPresent
	}
	info.PastebinCount = count
	return json.Marshal(info)
}

func (x *PwnedSource) Pastebin(email string) (int, error) {
	pastes, err := x.client.GetAccountPastes(email)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Not found") {
			return 0, nil
		} else {
			log.Error().Err(err).Str("component", "pwned").Str("email", email).Msg("pwned account paste")
			return 0, err
		}
	}
	count := 0
	for _, paste := range pastes {
		count += paste.EmailCount
	}
	return count, err
}

func (x *PwnedSource) Breaches(email string) (*PwnedInfo, error) {
	var info PwnedInfo
	info.Email = email
	info.BreachCount = 0
	info.BreachInfo = []*BreachInfo{}
	breaches, err := x.client.GetAccountBreaches(email, "", true, true)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			// NOTE:  the pwned API does not return a breachCount of zero, they return an error
			info.BreachCount = 0
			return &info, nil
		}
		log.Error().Err(err).Str("component", "pwned").Str("email", email).Msg("pwned account breaches")
		return &info, err
	}

	var l = len(breaches)
	if l > 0 {
		info.BreachCount = l
		days, err := common.GetDaysFromDate(breaches[0].BreachDate)
		if err == nil {
			info.MostRecentBreachInDays = days
		}
		for _, breach := range breaches {
			days, err = common.GetDaysFromDate(breach.BreachDate)
			if err == nil { // && days < info.MostRecentBreachInDays {
				info.MostRecentBreachInDays = days
				info.MostRecentBreachDate = breach.BreachDate
			}
			bi := &BreachInfo{Name: breach.Name, Title: breach.Title, Date: breach.BreachDate}
			info.BreachInfo = append(info.BreachInfo, bi)
		}
		info.MostRecentBreachInYears = info.MostRecentBreachInDays / 365
	}
	return &info, nil
}
