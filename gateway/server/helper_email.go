// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"fmt"
	"io"
	"net/smtp"
	"os"

	"github.com/mcnijman/go-emailaddress"
	"github.com/osintami/fingerprintz/log"
)

type IMail interface {
	SendMail(string, smtp.Auth, string, []string, []byte) error
}

type Email struct {
	mailer IMail
}

func NewEmail(mailer IMail) *Email {
	return &Email{mailer: mailer}
}

func (x *Email) Send(subject string, content string, recipient *emailaddress.EmailAddress, isHTML bool) error {
	from := os.Getenv("EMAIL_ALERT_FROM")
	password := os.Getenv("EMAIL_ALERT_PASSWORD")
	smtpHost := os.Getenv("EMAIL_ALERT_SMTP_SERVER")
	smtpPort := os.Getenv("EMAIL_ALERT_SMTP_PORT")

	var msg string
	if isHTML {
		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body := "<html><body><p>" + content + "</p></body></html>"
		msg = fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\n%s%s",
			from, recipient, subject, mime, body)
	} else {
		msg = fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\n\n%s",
			from, recipient, subject, content)
	}

	auth := smtp.PlainAuth(subject, from, password, smtpHost)

	return x.mailer.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		"noreply@osintami.com",
		[]string{recipient.String()},
		[]byte(msg))
}

func (x *Email) loadTemplate(fileName string) ([]byte, error) {
	fh, err := os.Open(fileName)
	if err != nil {
		log.Error().Err(err).Str("file", fileName).Msg("open")
		return nil, err
	}
	defer fh.Close()
	return io.ReadAll(fh)
}

type Sender struct {
}

func NewSender() *Sender {
	return &Sender{}
}

func (x *Sender) SendMail(server string, auth smtp.Auth, from string, recipients []string, content []byte) error {
	return smtp.SendMail(server, auth, from, recipients, content)
}
