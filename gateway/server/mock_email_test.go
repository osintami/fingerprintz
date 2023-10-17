// Copyright © 2023 OSINTAMI. This is not yours.
package server

import "net/smtp"

type MockSender struct {
	fail bool
}

func NewMockSender(fail bool) *MockSender {
	return &MockSender{
		fail: fail,
	}
}

func (x *MockSender) SendMail(server string, auth smtp.Auth, from string, recipients []string, content []byte) error {
	if x.fail {
		return ErrBadRequest
	}
	return nil
}
