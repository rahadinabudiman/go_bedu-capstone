package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"go_bedu/initializers"
	"go_bedu/models"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// 👇 Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *models.Administrator, data *EmailData) {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	fromName := config.FromName
	to := user.Email

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.ExecuteTemplate(&body, "verificationCode.html", &data)

	mailer := gomail.NewMessage()
	mailer.SetAddressHeader("From", from, fromName)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", data.Subject)
	mailer.SetBody("text/html", body.String())
	mailer.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	dialer := gomail.NewDialer(
		smtpHost,
		smtpPort,
		smtpUser,
		smtpPass,
	)

	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// // Send Email
	// if err := d.DialAndSend(m); err != nil {
	// 	log.Fatal("Could not send email: ", err)
	// }

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal("Could not send email: ", err.Error())
	}

	log.Println("Mail sent!")
}

func SendEmailUser(user *models.User, data *EmailData) {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	fromName := config.FromName
	to := user.Email

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.ExecuteTemplate(&body, "verificationCode.html", &data)

	mailer := gomail.NewMessage()
	mailer.SetAddressHeader("From", from, fromName)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", data.Subject)
	mailer.SetBody("text/html", body.String())
	mailer.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	dialer := gomail.NewDialer(
		smtpHost,
		smtpPort,
		smtpUser,
		smtpPass,
	)

	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// // Send Email
	// if err := d.DialAndSend(m); err != nil {
	// 	log.Fatal("Could not send email: ", err)
	// }
	fmt.Println("Ini dialer: ", dialer)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal("Could not send email: ", err.Error())
	}

	log.Println("Mail sent!")
}
