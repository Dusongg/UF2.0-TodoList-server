package config

import (
	"net/smtp"
	"time"
)

const (
	SmtpHost = "smtp.gmail.com"
	SmtpPort = "587"
	Sender   = "dusong700@gmail.com"
	Password = "cfvp lkgi igov frgh"
	Addr     = SmtpHost + ":" + SmtpPort
)

var Auth smtp.Auth
var now = time.Now()

// TimePoint1 上午九点发送一次
var TimePoint1 = time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())

// TimePoint2 下午一点发送一次
var TimePoint2 = time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, now.Location())

func init() {
	Auth = smtp.PlainAuth("", Sender, Password, SmtpHost)
}
