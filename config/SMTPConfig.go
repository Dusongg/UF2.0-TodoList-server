package config

import "net/smtp"

const (
	SmtpHost = "smtp.gmail.com"
	SmtpPort = "587"
	Sender   = "dusong700@gmail.com"
	Password = "dusong@041008"
)

var Auth = smtp.PlainAuth("", Sender, Password, SmtpHost)
var Addr = SmtpHost + ":" + SmtpPort
