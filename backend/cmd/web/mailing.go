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

func (sender *EmailSender)sendValidationEmail() error {
	htmlContent := fmt.Sprintf(`<a href="%s">Click here to validate your account</a>`, sender.validationURI) 
	
	e := email.NewEmail()
	e.From = "Jordan Wright <42pong1992@gmail.com>"
	e.To = []string{sender.destiantion}
	e.Bcc = []string{"42pong1992@gmail.com"}
	e.Cc = []string{"42pong1992@gmail.com"}
	e.Subject = "Validate your account on Matcha!!"
	e.Text = []byte("Click on the link below to validate your account")
	e.HTML = []byte(htmlContent)
	return	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "42pong1992@gmail.com", os.Getenv("EMAIL_APP_PASSWORD"), "smtp.gmail.com"))
}

func (sender *EmailSender)generateValidationURI() string {
	b := make([]byte, 16)
	rand.Read(b)
	sender.validationURI = hex.EncodeToString(b)
	return sender.validationURI
}