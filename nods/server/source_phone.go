package server

import (
	"context"
	"encoding/json"

	"github.com/davegardnerisme/phonegeocode"
	"github.com/nyaruka/phonenumbers"
	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

type PhoneInfo struct {
	RawPhoneNumber      string
	IsValid             bool
	IsTollFree          bool
	CountryISOCode      string
	CountryCode         int32
	PhoneNumber         uint64
	NationalNumber      string
	InternationalNumber string
}

type PhoneSource struct {
}

func NewPhoneSource(tools *Toolbox) IDataSource {
	return NewDataInstance(tools, SOURCE_OSINTAMI_NAME, &PhoneSource{})
}

func (x *PhoneSource) IsCached() bool {
	return false
}

func (x *PhoneSource) CategoryInfo(ctx context.Context, categoryName string, inputs common.DataInputs) (json.RawMessage, error) {
	switch categoryName {
	case CATEGORY_PHONE:
		return x.getPhoneInfo(inputs[categoryName])
	}

	return nil, ErrNotImplemented
}

func (x *PhoneSource) getPhoneInfo(phone string) (json.RawMessage, error) {
	info := PhoneInfo{}
	info.RawPhoneNumber = phone
	if len(info.RawPhoneNumber) > 10 {
		cc, err := phonegeocode.New().Country(info.RawPhoneNumber)
		if err != nil {
			log.Error().Err(err).Str("component", SOURCE_OSINTAMI_NAME).Msg("country code")
			info.IsValid = false
			data, _ := json.Marshal(info)
			return data, common.ErrNoDataPresent
		}
		info.CountryISOCode = cc
	} else {
		// NOTE:  not sure this is a good guess
		info.CountryISOCode = "US"
	}

	pns, err := phonenumbers.Parse(info.RawPhoneNumber, info.CountryISOCode)
	if err != nil {
		log.Error().Err(err).Str("component", SOURCE_OSINTAMI_NAME).Msg("country ISO code")
		info.IsValid = false
		data, _ := json.Marshal(info)
		return data, common.ErrNoDataPresent
	} else {
		info.CountryCode = *pns.CountryCode
		info.PhoneNumber = *pns.NationalNumber
		info.NationalNumber = phonenumbers.Format(pns, phonenumbers.NATIONAL)
		info.InternationalNumber = phonenumbers.Format(pns, phonenumbers.INTERNATIONAL)
		info.IsValid = true

		tollFreeCheck := info.NationalNumber[1:4]
		if tollFreeCheck == "800" || tollFreeCheck == "888" || tollFreeCheck == "867" || tollFreeCheck == "866" {
			info.IsTollFree = true
		}
	}

	return json.Marshal(info)
}
