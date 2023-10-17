// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"context"
	"encoding/json"
	"strings"
	"unicode"

	"github.com/mcnijman/go-emailaddress"
	"github.com/osintami/fingerprintz/common"
)

type EmailInfo struct {
	Email              string
	IsWeirdUserName    bool
	IsWeirdDomainName  bool
	IsEmailAlias       bool
	IsValidEmail       bool
	IsValidIcannSuffix bool
	IsNefariusDomain   bool
}

var riskyTopLevelDomains = []string{
	".xyz",
	".de",
	".icu",
	".ru",
	".cn",
	".uk",
	".tk",
	".ga",
	".cf",
	".org",
	".ml",
	".pw",
	".top",
	".icu",
	".info",
	".co",
	".work",
	".net",
	".club",
	".gq",
	".zw",
	".bd",
	".ke",
	".am",
	".sbs",
	".date",
	".quest",
	".cd",
	".bid"}

type EmailSource struct {
}

func NewEmailSource(tools *Toolbox) IDataSource {
	return NewDataInstance(tools, SOURCE_OSINTAMI_NAME, &EmailSource{})
}

func (x *EmailSource) IsCached() bool {
	return false
}

func (x *EmailSource) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	switch categoryName {
	case CATEGORY_EMAIL:

		return x.getEmailInfo(inputs[categoryName])
	}

	return nil, ErrNotImplemented
}

func (x *EmailSource) getEmailInfo(email string) (json.RawMessage, error) {
	var emailInfo EmailInfo
	emailInfo.Email = email
	em, err := emailaddress.Parse(email)
	if err != nil {
		emailInfo.IsValidEmail = false
		return json.Marshal(emailInfo)
	}

	emailInfo.IsValidEmail = true
	emailInfo.IsWeirdUserName = x.isNonHumanUserName(em.LocalPart)
	emailInfo.IsWeirdDomainName = x.isNonHumanDomainName(em.Domain)
	emailInfo.IsEmailAlias = strings.Contains(em.LocalPart, "+")
	emailInfo.IsValidIcannSuffix = (em.ValidateIcanSuffix() == nil)
	emailInfo.IsNefariusDomain = x.isRiskyTopLevelDomain(em.Domain)

	return json.Marshal(emailInfo)
}

func (x *EmailSource) isNonHumanUserName(name string) bool {
	if x.numerics(name) > 3 || x.nonAlphanumerics(name) > 2 {
		return true
	} else {
		return false
	}
}

func (x *EmailSource) nonAlphanumerics(name string) int {
	var count int = 0
	for _, letter := range name {
		if unicode.IsLetter(letter) || unicode.IsDigit(letter) {
			count++
		}
	}
	return len(name) - count
}

func (x *EmailSource) isNonHumanDomainName(name string) bool {
	if x.numerics(name) > 3 {
		return true
	} else {
		return false
	}
}

func (x *EmailSource) isRiskyTopLevelDomain(domain string) bool {
	for _, tld := range riskyTopLevelDomains {
		if strings.HasSuffix(domain, tld) {
			return true
		}
	}
	return false
}

func (x *EmailSource) numerics(name string) int {
	count := strings.Count(name, "0") +
		strings.Count(name, "1") +
		strings.Count(name, "2") +
		strings.Count(name, "3") +
		strings.Count(name, "4") +
		strings.Count(name, "5") +
		strings.Count(name, "6") +
		strings.Count(name, "7") +
		strings.Count(name, "8") +
		strings.Count(name, "9")
	return count
}
