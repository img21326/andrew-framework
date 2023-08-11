package helper

import (
	"fmt"
	"net/smtp"

	"github.com/spf13/viper"
)

var emailInstance *EmailHelper

type EmailHelper struct {
	host string
	port string
	user string
	pass string
}

func newEmailHelper() *EmailHelper {
	host := viper.GetViper().GetString("EMAIL_HOST")
	port := viper.GetViper().GetString("EMAIL_PORT")
	user := viper.GetViper().GetString("EMAIL_USER")
	pass := viper.GetViper().GetString("EMAIL_PASS")

	if host == "" || port == "" || user == "" || pass == "" {
		panic("Email config error")
	}

	return &EmailHelper{
		host: host,
		port: port,
		user: user,
		pass: pass,
	}
}

func GetEmailHelper() *EmailHelper {
	if emailInstance == nil {
		emailInstance = newEmailHelper()
	}
	return emailInstance
}

type EmailSendOption struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func (e *EmailHelper) SendEmail(option EmailSendOption) error {
	var from string
	if option.From == "" {
		from = e.user
	}

	message := []byte(option.Subject + option.Body)

	auth := smtp.PlainAuth("", e.user, e.pass, e.host)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", e.host, e.port), auth, from, option.To, message)
	return err
}
