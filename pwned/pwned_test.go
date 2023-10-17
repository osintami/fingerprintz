// Copyright © 2023 OSINTAMI. This is not yours.
package pwned

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestPwnedPasswords(t *testing.T) {
	pwned := NewPwned(nil, "")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	content := "00B0050BCC9F26E6D46C503FC2ED9627B5E:1\r\n00C0DFF90CFFAD774358304292784A146B8:3\r\n0167D0A60E720CD4B52ED263AABD9878249:1\r\n01B42E8D024D401E5F83FF48B883A3163EA:2\r\n01EF763E5AF82CCD345CD980ABEEC3576FA:1\r\n029A1C5E9F2F457897C85DD31FD37F42B7D:2\r\n02B6DAFD9C2B6EF89139E3F0DE11228474C:1\r\n02F0B40D6E097F7287E9241137D527A69BA:1\r\n0321DFCD8CAA0C7C7AE88972D7F75131239:1\r\n03A37E6B7CF45CB9A09C89ABD1D0054DBE6:4\r\n04476DFE56FC54C57C24477A446A16C9D5A:1\r\n049EA20513D919FAED92860A9227A03D019:2\r\n04DF32D8FFE77FE134B7F75D160164071B2:1\r\n04EE5"
	url := "https://api.pwnedpasswords.com/range/76272"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(200, content))

	count, err := pwned.HaveIBeenPwned("nope")
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}

func TestPwnedGetAccountPastes(t *testing.T) {
	client := resty.New().GetClient()
	pwned := NewPwned(client, "test-api-key")

	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	// breaches API
	content := "[{\"Name\":\"Adobe\",\"Title\":\"Adobe\",\"Domain\":\"adobe.com\",\"BreachDate\":\"2013-10-04\"}]"
	url := "https://haveibeenpwned.com/api/v3/breachedaccount/1%402.com?includeUnverified=true&truncateResponse=true"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(200, content))

	breaches, err := pwned.GetAccountBreaches("1@2.com", "", true, true)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(breaches))
}

func TestPwnedGetAccountBreaches(t *testing.T) {
	client := resty.New().GetClient()
	pwned := NewPwned(client, "test-api-key")
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	// pastbin API
	content := " [{\"Id\":\"za6e3zSs\",\"Source\":\"Pastebin\",\"Title\":null,\"Date\":\"2014-12-11T07:12:00Z\",\"EmailCount\":5046}]"
	url := "https://haveibeenpwned.com/api/v3/pasteaccount/1@2.com"
	httpmock.RegisterResponder(
		"GET", url, httpmock.NewStringResponder(200, content))

	pastes, err := pwned.GetAccountPastes("1@2.com")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(pastes))
}
