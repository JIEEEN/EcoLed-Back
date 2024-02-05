package services

import (
	"bytes"
	"crypto/rand"
	"net/smtp"
	"text/template"
)

type EmailServices struct{}

func (srv EmailServices) SendVerifyingEmail(subject string, templatePath string, to []string) (code string, err error) {
	//get the html body
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)

	if err != nil {
		return "", err
	}

	//Generate verification code
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	length := 6
	bytes :=make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = byte(letters[b%byte(len(letters))])
	}
	verificationCode := string(bytes)

	//Execute the template
	t.Execute(&body, struct{ Number string }{Number: verificationCode})

	//Send email
	auth := smtp.PlainAuth(
		"", 
		"hoon30512329@gmail.com", //email sender
		"uzxczxjovdckoiys", // email app password
		"smtp.gmail.com", // smtp server
	)
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	msg := "Subject: " + subject + "\n" + headers + "\n\n" + body.String()
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"hoon30512329@gmail.com",
		to,
		[]byte(msg),
	)
	if err != nil {
		return verificationCode, err
	}

	return verificationCode, nil
}
