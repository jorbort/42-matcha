package main

import (
	"github.com/jordan-wright/email"
	"os"
	"net/smtp"
	"fmt"
	"crypto/rand"
    "encoding/hex"

)

type EmailSender struct {
	destiantion string
	validationURI string
}

func (sender *EmailSender)sendValidationEmail(subject, mode , text string) error {
	if mode == "validate" {
		htmlContent := fmt.Sprintf(`<a href="http://localhost:3000/validate?code=%s">Click here to validate your account</a>`, sender.validationURI)
	} else if mode == "reset" {
		htmlContent := fmt.Sprintf(`<a href="http://localhost:3000/resetPassword?code=%s">Click here to reset your password</a>`, sender.validationURI)
	}
	e := email.NewEmail()
	e.From = "Matcha!! <42pong1992@gmail.com>"
	e.To = []string{sender.destiantion}
	e.Bcc = []string{"42pong1992@gmail.com"}
	e.Cc = []string{"42pong1992@gmail.com"}
	e.Subject = subject
	e.Text = []byte(text)
	e.HTML = []byte(htmlContent)
	return	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "42pong1992@gmail.com", os.Getenv("EMAIL_APP_PASSWORD"), "smtp.gmail.com"))
}

func (sender *EmailSender)generateValidationURI() []byte {
	b := make([]byte, 8)
	rand.Read(b)
	sender.validationURI = hex.EncodeToString(b)
	return []byte(sender.validationURI)
}
