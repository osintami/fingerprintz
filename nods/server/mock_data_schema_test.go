// Copyright © 2023 OSINTAMI. This is not yours.
package server

import "github.com/osintami/fingerprintz/common"

type MockDataSchema struct {
}

func NewMockDataSchema() *MockDataSchema {
	return &MockDataSchema{}
}

func (x *MockDataSchema) Item(dataURI *DataURI) (Item, error) {
	if dataURI.URI == "ip/ipsum/blacklist.isBlacklisted" {
		item := Item{
			Path:         "ip/ipsum/blacklist.isBlacklisted",
			CategoryName: CATEGORY_IPADDR,
			SourceName:   "ipsum",
			Enabled:      true,
			Gjson:        "ipsum.blacklist.isBlacklisted",
			Description:  "This IP address is blacklisted by one or more sources.",
			TypeName:     "Boolean",
			Type:         common.Boolean,
			Query:        "",
		}
		return item, nil
	}

	if dataURI.URI == "ip/uhb/blacklist.IsBlacklisted" {
		item := Item{
			Path:         "ip/uhb/blacklist.isBlacklisted",
			CategoryName: CATEGORY_IPADDR,
			SourceName:   "uhb",
			Enabled:      true,
			Gjson:        "uhb.blacklist.isBlacklisted",
			Description:  "This IP address is blacklisted by one or more sources.",
			TypeName:     "Boolean",
			Type:         common.Boolean,
			Query:        "",
		}
		return item, nil
	}

	//"[email/pwned/breachCount] > 0 && [email/pwned/breachAgeDate] > @{1.years.ago}",

	if dataURI.URI == "email/pwned/breachCount" {
		item := Item{
			Path:         "email/pwned/breachCount",
			CategoryName: CATEGORY_EMAIL,
			SourceName:   "pwned",
			Enabled:      true,
			Gjson:        "BreachCount",
			Description:  "Email breach count.",
			TypeName:     "Integer",
			Type:         common.Integer,
			Query:        "",
		}
		return item, nil
	}

	if dataURI.URI == "email/pwned/breachAgeDate" {
		item := Item{
			Path:         "email/pwned/breachAgeDate",
			CategoryName: CATEGORY_EMAIL,
			SourceName:   "pwned",
			Enabled:      true,
			Gjson:        "MostRecentBreachDate",
			Description:  "Most recent breach date.",
			TypeName:     "Date",
			Type:         common.Date,
			Query:        "",
		}
		return item, nil
	}

	return Item{}, ErrItemNotFound
}

func (x *MockDataSchema) IsValidItem(dataURI *DataURI) bool {
	_, err := x.Item(dataURI)
	return err == nil
}

func (x *MockDataSchema) IsValidCategory(categoryName string) bool {
	return categoryName != "nope"
}

func (x *MockDataSchema) IsEnabled(sourceName string) bool {
	if sourceName == "ipsum" || sourceName == "uhb" {
		return true
	}
	return false
}

func (x *MockDataSchema) Source(sourceName string) (SourceInfo, error) {
	return SourceInfo{}, nil
}

func (x *MockDataSchema) ListItems() []Item {
	items := []Item{}

	item := Item{
		Path:         "ip/ipsum/blacklist.isBlacklisted",
		CategoryName: CATEGORY_IPADDR,
		SourceName:   "ipsum",
		Enabled:      true,
		Gjson:        "ipsum.blacklist.isBlacklisted",
		Description:  "ipsum blacklisted",
		TypeName:     "Boolean",
		Type:         common.Boolean,
		Query:        "",
	}

	items = append(items, item)

	item = Item{
		Path:         "ip/uhb/blacklist.isBlacklisted",
		CategoryName: CATEGORY_IPADDR,
		SourceName:   "uhb",
		Enabled:      true,
		Gjson:        "uhb.blacklist.isBlacklisted",
		Description:  "uhb blacklisted",
		TypeName:     "Boolean",
		Type:         common.Boolean,
		Query:        "",
	}

	items = append(items, item)

	item = Item{
		Path:         "rule/osintami/isBlacklisted",
		CategoryName: CATEGORY_RULE,
		SourceName:   "osintami",
		Enabled:      true,
		Gjson:        "isBlacklisted",
		Description:  "rule blacklisted",
		TypeName:     "Boolean",
		Type:         common.Boolean,
		Query:        "[ip/ipsum/blacklist.isBlacklisted] || [ip/uhb/blacklist.isBlacklisted]",
	}

	items = append(items, item)

	item = Item{
		Path:         "rule/osintami/email.hasRecentBreach",
		CategoryName: CATEGORY_RULE,
		SourceName:   "osintami",
		Enabled:      true,
		Gjson:        "isCompromised",
		Description:  "rule recent data breach",
		TypeName:     "Boolean",
		Type:         common.Boolean,
		Query:        "[email/pwned/breachCount] > 0 && [email/pwned/breachAgeDate] > @{1.years.ago}",
	}

	items = append(items, item)

	item = Item{
		Path:         "rule/osintami/expression",
		CategoryName: CATEGORY_RULE,
		SourceName:   "osintami",
		Enabled:      true,
		Gjson:        "isBroken",
		Description:  "broken rule",
		TypeName:     "Boolean",
		Type:         common.Boolean,
		Query:        "1<=nope",
	}

	items = append(items, item)

	return items
}

func (x *MockDataSchema) ListSources() []SourceInfo {
	out := []SourceInfo{}

	info := SourceInfo{
		Name:     "ipsum",
		Database: "mmdb",
		Enabled:  true}

	out = append(out, info)

	info = SourceInfo{
		Name:     "uhb",
		Database: "mmdb",
		Enabled:  true}

	out = append(out, info)

	info = SourceInfo{
		Name:     "osintami",
		Database: "code",
		Enabled:  true}

	out = append(out, info)

	info = SourceInfo{
		Name:     "palo",
		Database: "byod",
		Enabled:  true}

	out = append(out, info)

	info = SourceInfo{
		Name:     "fakefilter",
		Database: "fast",
		Enabled:  true}

	out = append(out, info)

	info = SourceInfo{
		Name:     "disabled",
		Database: "nope",
		Enabled:  false}

	out = append(out, info)

	return out
}

func (x *MockDataSchema) ListCategories() []string {
	return CATEGORIES
}

func (x *MockDataSchema) ListItemsByCategory(categoryName string) []Item {
	out := []Item{}
	for _, item := range x.ListItems() {
		if categoryName == item.CategoryName {
			out = append(out, item)
		}
	}
	return out
}

func (x *MockDataSchema) ListRulesItems() map[string]Item {
	rules := make(map[string]Item)
	for _, item := range x.ListItems() {
		if item.CategoryName == CATEGORY_RULE || item.Query != "" {
			rules[item.Path] = item
		}
	}
	return rules
}
