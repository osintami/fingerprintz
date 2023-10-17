// Copyright © 2023 OSINTAMI. This is not yours.
package server

type MockJWTSigner struct {
	fingerprint *Fingerprint
	decodeFail  bool
	encodeFail  bool
}

func NewMockJWTSigner(decodeFail, encodeFail bool) IJWTSigner {
	return &MockJWTSigner{
		decodeFail: decodeFail,
		encodeFail: encodeFail}
}

func (x *MockJWTSigner) SignJWT(fingerprint *Fingerprint) (string, error) {
	if !x.encodeFail {
		x.fingerprint = fingerprint
		return "1234567890", nil
	} else {
		return "", ErrFingerprintSmudged
	}
}

func (x *MockJWTSigner) DecodeJWT(signature string) (*Fingerprint, error) {
	if !x.decodeFail {
		if x.fingerprint == nil {
			x.fingerprint = &Fingerprint{}
		}
		return x.fingerprint, nil
	} else {
		return nil, ErrFingerprintSmudged
	}
}
