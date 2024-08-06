package main

import (
	"OrderManager/config"
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

func testSendEmail() {
	//email := "dusong700@gmail.com"
	email := "728869268@qq.com"
	msg := []byte("测试邮件")

	err := smtp.SendMail(config.Addr, config.Auth, config.Sender, []string{email}, msg)
	if err != nil {
		log.Println("发送失败：", err)
	} else {
		log.Println("发送成功")
	}
}

func emailClock() {
	for {
		now := time.Now()

		// 计算下一次发送时间
		next9AM := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
		next1PM := time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, now.Location())

		if now.After(next9AM) {
			next9AM = next9AM.Add(24 * time.Hour)
		}
		if now.After(next1PM) {
			next1PM = next1PM.Add(24 * time.Hour)
		}

		// 等待到下一个发送时间
		if next9AM.Before(next1PM) {
			time.Sleep(next9AM.Sub(now))
			queryAndSendEmail()

		} else {
			time.Sleep(next1PM.Sub(now))
			queryAndSendEmail()
		}
	}
}

func queryAndSendEmail() {
	type nameAndEmail struct {
		Name  string
		Email string
	}
	var ne []nameAndEmail
	if err := db.Table("user_table").Select("name", "email").Find(&ne).Error; err != nil {
		log.Println(err)
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
		log.Println(err)
	}
	msg := make([]byte, 0)
	for _, content := range contents {
		row := fmt.Sprintf("task_id:%s req_no:%s deadline:%s comment:%s taskTime:%s\n", content.taskId, content.reqNo, content.deadline, content.comment, content.estimatedWorkHour)
		msg = append(msg, []byte(row)...)
	}
	err = smtp.SendMail(config.Addr, config.Auth, config.Sender, []string{email}, msg)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return nil, err
	}
	return ct, nil
}
