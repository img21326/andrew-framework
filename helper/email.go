package helper

import (
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

var emailInstance *EmailHelper

type EmailHelper struct {
	host string
	port int
	user string
	pass string
}

func newEmailHelper() *EmailHelper {
	host := viper.GetViper().GetString("EMAIL_HOST")
	port := viper.GetViper().GetInt("EMAIL_PORT")
	user := viper.GetViper().GetString("EMAIL_USER")
	pass := viper.GetViper().GetString("EMAIL_PASS")

	if host == "" || port == 0 || user == "" || pass == "" {
		return nil
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

func NewEmailInstance(host string, port int, user string, pass string) *EmailHelper {
	return &EmailHelper{
		host: host,
		port: port,
		user: user,
		pass: pass,
	}
}

type EmailSendOption struct {
	From    string
	To      []string
	Subject string
	IsHtml  bool
	Body    string
}

func (e *EmailHelper) SendEmail(option EmailSendOption) error {
	var from string
	if option.From == "" {
		from = e.user
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", option.To...)
	msg.SetHeader("Subject", option.Subject)
	// text/html for a html email
	if option.IsHtml {
		msg.SetBody("text/html", option.Body)
	} else {
		msg.SetBody("text/plain", option.Body)
	}

	n := gomail.NewDialer(e.host, e.port, from, e.pass)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
