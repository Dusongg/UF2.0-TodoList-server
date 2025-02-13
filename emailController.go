package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/smtp"
	"sync"
	"time"
	"github.com/robfig/cron"
)

func testSendEmail() {
	//email := "dusong700@gmail.com"
	email := "728869268@qq.com"
	msg := []byte("测试邮件")
	auth := smtp.PlainAuth("", Conf.SMTP.Sender, Conf.SMTP.SenderPassword, Conf.SMTP.Host)
	err := smtp.SendMail(Conf.SMTP.Host+":"+Conf.SMTP.Port, auth, Conf.SMTP.Sender, []string{email}, msg)
	if err != nil {
		log.Println("发送失败：", err)
	} else {
		log.Println("发送成功")
	}
}

func improvedEmailClock() {
	scheduler := cron.New()
	
	// 添加定时任务
	scheduler.AddFunc("0 9 * * *", queryAndSendEmail)  // 每天9点
	scheduler.AddFunc("0 13 * * *", queryAndSendEmail) // 每天13点
	
	scheduler.Start()
}

func queryAndSendEmail() {
	type nameAndEmail struct {
		Name  string
		Email string
	}
	var ne []nameAndEmail
	if err := db.Table("user_table").Select("name", "email").Find(&ne).Error; err != nil {
		logrus.Warning(err)
		return
	}

	var wg sync.WaitGroup
	for _, item := range ne {
		wg.Add(1)
		name, email := item.Name, item.Email
		go func() {
			wg.Done()
			sendEmail(name, email)
		}()
	}
	wg.Wait()

}

func sendEmail(name string, email string) {
	contents, err := querySendContent(name)
	if err != nil {
		logrus.Warning(err)
	}
	if len(contents) == 0 {
		return
	}
	msg := make([]byte, 0)
	for _, content := range contents {
		row := fmt.Sprintf("task_id:%s req_no:%s deadline:%s comment:%s taskTime:%f\n", content.taskId, content.reqNo, content.deadline, content.comment, content.estimatedWorkHour)
		msg = append(msg, []byte(row)...)
	}
	auth := smtp.PlainAuth("", Conf.SMTP.Sender, Conf.SMTP.SenderPassword, Conf.SMTP.Host)

	err = smtp.SendMail(Conf.SMTP.Host+":"+Conf.SMTP.Port, auth, Conf.SMTP.Sender, []string{email}, msg)
	if err != nil {
		logrus.Warning(err)
	}
}

type sendContent struct {
	taskId            string
	reqNo             string
	deadline          time.Time
	comment           string
	estimatedWorkHour float64
}

func querySendContent(name string) ([]sendContent, error) {

	var ct []sendContent
	if err := db.Table("tasklist_table").Select("task_id", "req_no", "deadline", "comment", "estimated_work_hours").Where("principal = ?", name).Find(&ct).Error; err != nil {
		return nil, err
	}
	return ct, nil
}
