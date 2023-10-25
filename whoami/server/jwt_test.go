// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	jwt := NewJWTSigner(NewJWT("test-jwt-secret"))
	info := &Fingerprint{}
	info.LastSeenAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	info.EHash = "test-email-hash"
	info.Latitude = 30.000
	info.Longitude = 97.000

	info.City = "test-city"
	info.Country = "test-country"
	info.IpAddr = "1.2.3.4"
	info.UserAgent = "test-user-agent"
	info.DeviceId = "test-device-id"
	info.NetworkId = "0x0020"
	info.PartnerId = "test-partner-id"

	signature, err := jwt.SignJWT(info)
	assert.Nil(t, err)

	out, err := jwt.DecodeJWT(signature)
	assert.Nil(t, err)
	assert.Equal(t, info.LastSeenAt, out.LastSeenAt)
}

func TestJWTDecodeCorrupt(t *testing.T) {
	jwt := NewJWTSigner(NewJWT(""))
	out, err := jwt.DecodeJWT("header.payload.signature")
	assert.Equal(t, ErrFingerprintSmudged, err)
	assert.Nil(t, out)
}
