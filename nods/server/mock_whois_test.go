// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"time"

	"github.com/osintami/fingerprintz/whois"
)

// -----------------------------------------------------------------------------
type MockWhois struct {
	failParse bool
	failGet   bool
}

func NewMockWhois(failParse, failGet bool) whois.IWhois {
	return &MockWhois{
		failParse: failParse,
		failGet:   failGet}
}

func (x *MockWhois) ParseInfo(domain, output string) (*whois.WhoisInfo, error) {
	if x.failParse {
		return nil, ErrBadData
	}
	return whois.NewWhois().ParseInfo(domain, output)
}

func (x *MockWhois) Get(domain string, timeout string) (string, error) {
	if domain == "com" {
		return "", whois.ErrInvalidDomainName
	}

	if x.failGet {
		return "", ErrBadData
	}

	return `WHOIS    Domain Name: OSINTAMI.COM
	Registry Domain ID: 2740471795_DOMAIN_COM-VRSN
	Registrar WHOIS Server: whois.google.com
	Registrar URL: http://domains2.squarespace.com
	Updated Date: 2022-11-24T00:48:06Z
	Creation Date: 2022-11-24T00:48:04Z
	Registry Expiry Date: 2023-11-24T00:48:04Z
	Registrar: Squarespace Domains II LLC
	Registrar IANA ID: 895
	Registrar Abuse Contact Email: registrar-abuse@google.com
	Registrar Abuse Contact Phone: +1.8772376466
	Domain Status: clientTransferProhibited https://icann.org/epp#clientTransferProhibited
	Name Server: NS-CLOUD-B1.GOOGLEDOMAINS.COM
	Name Server: NS-CLOUD-B2.GOOGLEDOMAINS.COM
	Name Server: NS-CLOUD-B3.GOOGLEDOMAINS.COM
	Name Server: NS-CLOUD-B4.GOOGLEDOMAINS.COM
	DNSSEC: signedDelegation
	DNSSEC DS Data: 61462 8 2 3C1D289A89DE6D411B4E0ADD0C579A49D92509384E24AF4A19C058842ADCA5F8
	URL of the ICANN Whois Inaccuracy Complaint Form: https://www.icann.org/wicf/
 >>> Last update of whois database: 2023-09-21T17:57:27Z <<<
 
 For more information on Whois status codes, please visit https://icann.org/epp
 
 NOTICE: The expiration date displayed in this record is the date the
 registrar's sponsorship of the domain name registration in the registry is
 currently set to expire. This date does not necessarily reflect the expiration
 date of the domain name registrant's agreement with the sponsoring
 registrar.  Users may consult the sponsoring registrar's Whois database to
 view the registrar's reported date of expiration for this registration.
 
 TERMS OF USE: You are not authorized to access or query our Whois
 database through the use of electronic processes that are high-volume and
 automated except as reasonably necessary to register domain names or
 modify existing registrations; the Data in VeriSign Global Registry
 Services' ("VeriSign") Whois database is provided by VeriSign for
 information purposes only, and to assist persons in obtaining information
 about or related to a domain name registration record. VeriSign does not
 guarantee its accuracy. By submitting a Whois query, you agree to abide
 by the following terms of use: You agree that you may use this Data only
 for lawful purposes and that under no circumstances will you use this Data
 to: (1) allow, enable, or otherwise support the transmission of mass
 unsolicited, commercial advertising or solicitations via e-mail, telephone,
 or facsimile; or (2) enable high volume, automated, electronic processes
 that apply to VeriSign (or its computer systems). The compilation,
 repackaging, dissemination or other use of this Data is expressly
 prohibited without the prior written consent of VeriSign. You agree not to
 use electronic processes that are automated and high-volume to access or
 query the Whois database except as reasonably necessary to register
 domain names or modify existing registrations. VeriSign reserves the right
 to restrict your access to the Whois database in its sole discretion to ensure
 operational stability.  VeriSign may restrict or terminate your access to the
 Whois database for failure to abide by these terms of use. VeriSign
 reserves the right to modify these terms at any time.
 
 The Registry database contains ONLY .COM, .NET, .EDU domains and
 Registrars.`, nil
}

func (x *MockWhois) GetWithTimeout(domain string, timeout time.Duration) (string, error) {
	return "", nil
}
