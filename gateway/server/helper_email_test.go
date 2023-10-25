// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"testing"

	"github.com/mcnijman/go-emailaddress"
	"github.com/stretchr/testify/assert"
)

func TestSendEmaill(t *testing.T) {
	email := NewEmail(NewMockSender(false))
	recipient, _ := emailaddress.Parse("test@example.com")
	email.Send("test", "this is the body", recipient, false)
	email.Send("test", "this is the body", recipient, true)
}

func TestLoadTemplate(t *testing.T) {
	email := NewEmail(NewMockSender(false))
	body, err := email.loadTemplate("../welcome.template")
	assert.NotNil(t, body)
	assert.Nil(t, err)

	body, err = email.loadTemplate("nope.template")
	assert.Nil(t, body)
	assert.NotNil(t, err)
}

func TestNewSender(t *testing.T) {
	err := NewSender().SendMail("", nil, "", nil, nil)
	assert.NotNil(t, err)
}
