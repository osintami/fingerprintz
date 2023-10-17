// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJWT interface {
	Sign(claims jwt.Claims) (string, error)
	Parse(signature string) (*jwt.Token, error)
}

type JWT struct {
	secret string
}

func NewJWT(secret string) IJWT {
	return &JWT{}
}

func (x *JWT) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(x.secret))
}

func (x *JWT) Parse(signature string) (*jwt.Token, error) {
	return jwt.Parse(signature, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrFingerprintSmudged
		}
		return []byte(x.secret), nil
	})
}

type IJWTSigner interface {
	SignJWT(fingerprint *Fingerprint) (string, error)
	DecodeJWT(signature string) (*Fingerprint, error)
}

type JWTSigner struct {
	jwt IJWT
}

func NewJWTSigner(jwt IJWT) IJWTSigner {
	return &JWTSigner{jwt: jwt}
}

func (x *JWTSigner) SignJWT(fingerprint *Fingerprint) (string, error) {
	claims := FingerprintClaims{
		fingerprint.LastSeenAt,
		fingerprint.EHash,
		fingerprint.Latitude,
		fingerprint.Longitude,
		fingerprint.City,
		fingerprint.Country,
		fingerprint.IpAddr,
		fingerprint.UserAgent,
		fingerprint.DeviceId,
		fingerprint.NetworkId,
		fingerprint.PartnerId,
		"1.0.0",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "osintami",
			Subject:   "fingerprint",
		},
	}
	return x.jwt.Sign(claims)
}

func (x *JWTSigner) DecodeJWT(signature string) (*Fingerprint, error) {
	if token, err := x.jwt.Parse(signature); err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fingerprint := &Fingerprint{}
			fingerprint.LastSeenAt = claims["time"].(string)
			fingerprint.EHash = claims["ehash"].(string)
			fingerprint.Latitude = claims["latitude"].(float64)
			fingerprint.Longitude = claims["longitude"].(float64)
			fingerprint.City = claims["city"].(string)
			fingerprint.Country = claims["country"].(string)
			fingerprint.IpAddr = claims["ip"].(string)
			fingerprint.UserAgent = claims["ua"].(string)
			fingerprint.DeviceId = claims["hwid"].(string)
			fingerprint.NetworkId = claims["nid"].(string)
			fingerprint.PartnerId = claims["pid"].(string)
			return fingerprint, nil
		}
	}
	return nil, ErrFingerprintSmudged
}
