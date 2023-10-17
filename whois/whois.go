// -------------------------
//
// Copyright 2015, undiabler
//
// git: github.com/undiabler/golang-whois
//
// http://undiabler.com
//
// Released under the Apache License, Version 2.0
//
//--------------------------

package whois

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/osintami/fingerprintz/common"
	"github.com/osintami/fingerprintz/log"
)

var ErrInvalidDomainName = errors.New("invalid domain name")
var ErrNoServerForZone = errors.New("no server for zone")
var ErrNetwork = errors.New("whois server network error")
var ErrInvalidTimeout = errors.New("invalid network timeout milliseconds")

type IWhois interface {
	Get(string, string) (string, error)
	ParseInfo(string, string) (*WhoisInfo, error)
}

type WhoisInfo struct {
	Domain             string
	DomainAgeInDays    int
	DomainAgeInYears   int
	DomainAgeDate      string
	IsRegisteredDomain bool
}

type Whois struct {
}

func NewWhois() IWhois {
	return &Whois{}
}

// parse age and regisered status
func (x *Whois) ParseInfo(domain, result string) (*WhoisInfo, error) {
	info := &WhoisInfo{}
	info.Domain = domain

	for _, line := range strings.Split(result, "\n") {
		if strings.Contains(line, "Creation Date:") {
			date := strings.Split(line, ": ")[1]
			date = strings.Split(date, "T")[0]
			days, err := common.GetDaysFromDate(date)
			if err != nil {
				return info, err
			}

			info.DomainAgeDate = date
			info.DomainAgeInDays = days
			if days > 0 {
				info.DomainAgeInYears = days / 365
			}
			info.IsRegisteredDomain = true
			break
		}
	}
	return info, nil
}

// whois request with 2 second timeout
func (x *Whois) Get(domain string, d string) (string, error) {
	timeout, err := time.ParseDuration(d)
	if err != nil {
		return "", ErrInvalidTimeout
	}
	return x.GetWithTimeout(domain, timeout)
}

func (x *Whois) GetWithTimeout(domain string, timeout time.Duration) (string, error) {

	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		log.Error().Err(ErrInvalidDomainName).Str("component", "whois").Str("domain", domain).Msg("top level domain")
		return "", ErrInvalidDomainName
	}

	// last part of domain is zome
	zone := parts[len(parts)-1]
	server, ok := servers[zone]

	if !ok {
		log.Error().Err(ErrNoServerForZone).Str("component", "whois").Str("domain", domain).Str("zone", zone).Msg("no server for zone")
		return "", ErrNoServerForZone
	}

	result, err := x.call(server, domain, timeout)
	if err != nil {
		return "", ErrNetwork
	}
	return result, nil
}

func (x *Whois) call(server, domain string, timeout time.Duration) (string, error) {

	tcpserveraddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(server, "43"))
	if err != nil {
		log.Error().Err(err).Str("component", "whois").Str("server", server).Msg("whois server name resolution")
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var dialer net.Dialer
	connection, err := dialer.DialContext(ctx, tcpserveraddr.Network(), tcpserveraddr.String())
	if err != nil {
		log.Error().Err(err).Str("component", "whois").Str("server", server).Str("domain", domain).Msg("whois server connect")
		return "", err
	}
	defer connection.Close()

	connection.Write([]byte(domain + "\r\n"))
	buffer, err := io.ReadAll(connection)
	if err != nil {
		log.Error().Err(err).Str("component", "whois").Str("domain", domain).Msg("whois server read")
		return "", err
	}
	return string(buffer[:]), nil
}
