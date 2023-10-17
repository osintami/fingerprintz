// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"errors"

	gopwned "github.com/mavjs/goPwned"
	"github.com/osintami/fingerprintz/pwned"
)

type MockPwned struct {
	failPastes   bool
	failBreaches bool
	failNotFound bool
}

func NewMockPwned(failPastes, failBreaches, failNotFound bool) pwned.IPwned {
	return &MockPwned{failPastes: failPastes, failBreaches: failBreaches, failNotFound: failNotFound}
}

func (x *MockPwned) GetAccountPastes(string) ([]*gopwned.Paste, error) {
	if x.failNotFound {
		return nil, errors.New("Not found")
	}
	if x.failPastes {
		return nil, ErrBadData
	}
	pastes := make([]*gopwned.Paste, 0)
	pastes = append(pastes, &gopwned.Paste{EmailCount: 1})
	return pastes, nil
}

func (x *MockPwned) GetAccountBreaches(string, string, bool, bool) ([]*gopwned.Breach, error) {
	if x.failNotFound {
		return nil, errors.New("Not found")
	}
	if x.failBreaches {
		return nil, ErrBadData
	}
	breaches := make([]*gopwned.Breach, 0)
	breaches = append(breaches, &gopwned.Breach{BreachDate: "2023-10-01"})
	return breaches, nil
}

func (x *MockPwned) HaveIBeenPwned(string) (int, error) {
	if x.failNotFound {
		return 0, errors.New("Not found")
	}
	return 0, nil
}
