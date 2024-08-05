package config

import "net/smtp"

const (
	SmtpHost = "smtp.gmail.com"
	SmtpPort = "587"
	Sender   = "dusong700@gmail.com"
	Password = "cfvp lkgi igov frgh"
	Addr     = SmtpHost + ":" + SmtpPort
)

var Auth smtp.Auth

func init() {
	Auth = smtp.PlainAuth("", Sender, Password, SmtpHost)
}
