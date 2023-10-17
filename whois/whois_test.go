// Copyright © 2023 OSINTAMI. This is not yours.
package whois

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// NOTE:  these tests use a live network connection
func TestWhois(t *testing.T) {
	whois := NewWhois()
	out, err := whois.Get("nope.com", "1s")
	assert.Nil(t, err)
	assert.NotNil(t, out)
}

func TestWhoisTLD(t *testing.T) {
	whois := NewWhois()
	out, err := whois.Get("xyz", "1s")
	assert.Equal(t, ErrInvalidDomainName, err)
	assert.Equal(t, "", out)
}

func TestWhoisInvalidZone(t *testing.T) {
	whois := NewWhois()
	out, err := whois.Get("xyz.123", "1s")
	assert.Equal(t, ErrNoServerForZone, err)
	assert.Equal(t, "", out)
}

func TestWhoisTimeout(t *testing.T) {
	whois := NewWhois()
	out, err := whois.Get("1.com", "0ms")
	assert.Equal(t, ErrNetwork, err)
	assert.Equal(t, "", out)
}

func TestWhoisInvalidTimeout(t *testing.T) {
	whois := NewWhois()
	out, err := whois.Get("1.com", "")
	assert.Equal(t, ErrNetwork, err)
	assert.Equal(t, "", out)
}

func TestWhoisParseInfo(t *testing.T) {
	whois := NewWhois()

	var output = `Domain Name: OSINTAMI.COM
Registry Domain ID: 2740471795_DOMAIN_COM-VRSN
Registrar WHOIS Server: whois.squarespace.domains
Registrar URL: http://domains2.squarespace.com
Updated Date: 2022-11-24T00:48:06Z
Creation Date: 2022-11-24T00:48:04Z
Registry Expiry Date: 2023-11-24T00:48:04Z
Registrar: Squarespace Domains II LLC
Registrar IANA ID: 895
Registrar Abuse Contact Email: abuse-complaints@squarespace.com
Registrar Abuse Contact Phone: +1.6466935324
Domain Status: clientTransferProhibited https://icann.org/epp#clientTransferProhibited
Name Server: NS-CLOUD-B1.GOOGLEDOMAINS.COM
Name Server: NS-CLOUD-B2.GOOGLEDOMAINS.COM
Name Server: NS-CLOUD-B3.GOOGLEDOMAINS.COM
Name Server: NS-CLOUD-B4.GOOGLEDOMAINS.COM
DNSSEC: signedDelegation
DNSSEC DS Data: 61462 8 2 3C1D289A89DE6D411B4E0ADD0C579A49D92509384E24AF4A19C058842ADCA5F8
URL of the ICANN Whois Inaccuracy Complaint Form: https://www.icann.org/wicf/
>>> Last update of whois database: 2023-10-08T05:43:50Z <<<`

	info, err := whois.ParseInfo("osintami.com", output)
	assert.Nil(t, err)
	assert.Equal(t, "osintami.com", info.Domain)
}

func TestWhoisParseInfoBadDate(t *testing.T) {
	whois := NewWhois()

	var output = `Domain Name: OSINTAMI.COM
Registry Domain ID: 2740471795_DOMAIN_COM-VRSN
Registrar WHOIS Server: whois.squarespace.domains
Registrar URL: http://domains2.squarespace.com
Updated Date: 2022-11-24T00:48:06Z
Creation Date: nope
Registry Expiry Date: 2023-11-24T00:48:04Z
Registrar: Squarespace Domains II LLC
Registrar IANA ID: 895
Registrar Abuse Contact Email: abuse-complaints@squarespace.com
Registrar Abuse Contact Phone: +1.6466935324
Domain Status: clientTransferProhibited https://icann.org/epp#clientTransferProhibited
Name Server: NS-CLOUD-B1.GOOGLEDOMAINS.COM
Name Server: NS-CLOUD-B2.GOOGLEDOMAINS.COM
Name Server: NS-CLOUD-B3.GOOGLEDOMAINS.COM
Name Server: NS-CLOUD-B4.GOOGLEDOMAINS.COM
DNSSEC: signedDelegation
DNSSEC DS Data: 61462 8 2 3C1D289A89DE6D411B4E0ADD0C579A49D92509384E24AF4A19C058842ADCA5F8
URL of the ICANN Whois Inaccuracy Complaint Form: https://www.icann.org/wicf/
>>> Last update of whois database: 2023-10-08T05:43:50Z <<<`

	info, err := whois.ParseInfo("osintami.com", output)
	assert.NotNil(t, err)
	assert.Equal(t, "osintami.com", info.Domain)
}
