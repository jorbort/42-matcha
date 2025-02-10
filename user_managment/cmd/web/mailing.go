package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

type EmailSender struct {
	Destiantion   string
	ValidationURI string
}

func (sender *EmailSender) sendValidationEmail(subject, mode, text string) error {
	var htmlContent string
	if mode == "validate" {
		htmlContent = fmt.Sprintf(`<a href="http://localhost:3000/validate?code=%s">Click here to validate your account</a>`, sender.ValidationURI)
	} else if mode == "reset" {
		htmlContent = fmt.Sprintf(`<a href="http://localhost:3000/resetPassword?code=%s">Click here to reset your password</a>`, sender.ValidationURI)
	}
	e := email.NewEmail()
	e.From = "Matcha!! <42pong1992@gmail.com>"
	e.To = []string{sender.Destiantion}
	e.Bcc = []string{"42pong1992@gmail.com"}
	e.Cc = []string{"42pong1992@gmail.com"}
	e.Subject = subject
	e.Text = []byte(text)
	e.HTML = []byte(htmlContent)
	return e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "42pong1992@gmail.com", os.Getenv("EMAIL_APP_PASSWORD"), "smtp.gmail.com"))
}

func (sender *EmailSender) generateValidationURI() []byte {
	b := make([]byte, 8)
	rand.Read(b)
	sender.ValidationURI = hex.EncodeToString(b)
	return []byte(sender.ValidationURI)
}
