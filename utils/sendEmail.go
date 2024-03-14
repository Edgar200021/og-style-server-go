package utils

import (
	"gopkg.in/gomail.v2"
	"os"
)

func SendEmail(htmlString, subject string, to ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "easatryan2000@gmail.com")
	m.SetHeader("To", to...)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlString)

	d := gomail.NewDialer("smtp.gmail.com", 465, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
