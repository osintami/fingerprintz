// Copyright © 2023 OSINTAMI. This is not yours.
package pwned

import (
	"net/http"

	"github.com/TRIKKSS/haveibeenpwnedpkg"
	gopwned "github.com/mavjs/goPwned"
)

type IPwned interface {
	GetAccountPastes(string) ([]*gopwned.Paste, error)
	GetAccountBreaches(string, string, bool, bool) ([]*gopwned.Breach, error)
	HaveIBeenPwned(string) (int, error)
}

type Pwned struct {
	client *gopwned.Client
}

func NewPwned(client *http.Client, apikey string) IPwned {
	return &Pwned{client: gopwned.NewClient(client, apikey)}
}

func (x *Pwned) GetAccountPastes(email string) ([]*gopwned.Paste, error) {
	return x.client.GetAccountPastes(email)
}

func (x *Pwned) GetAccountBreaches(account string, domain string, truncate bool, unverified bool) ([]*gopwned.Breach, error) {
	return x.client.GetAccountBreaches(account, domain, truncate, unverified)
}

func (x *Pwned) HaveIBeenPwned(password string) (int, error) {
	return haveibeenpwnedpkg.HaveibeenpwnPassword(password)
}
