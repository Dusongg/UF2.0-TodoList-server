package main

import (
	"OrderManager/config"
	"fmt"
	"log"
	"net/smtp"
	"time"
)

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
	rows, err := db.Query("select name, email from user_table")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name, email string
		err := rows.Scan(&name, &email)
		if err != nil {
			log.Println(err)
		}
		sendEmail(name, email)
	}

}

func sendEmail(name string, email string) {
	contents, err := querySendContent(name)
	if err != nil {
		log.Println(err)
	}
	msg := make([]byte, 0)
	for _, content := range contents {
		row := fmt.Sprintf("task_id:%s req_no:%s deadline:%s comment:%s taskTime:%s\n", content.taskId, content.reqNo, content.deadline, content.comment, content.taskTime)
		msg = append(msg, []byte(row)...)
	}
	err = smtp.SendMail(config.Addr, config.Auth, config.Sender, []string{email}, msg)
	if err != nil {
		log.Println(err)
	}
}

type sendContent struct {
	taskId   string
	reqNo    string
	deadline string
	taskTime float64
	comment  string
}

func querySendContent(name string) ([]*sendContent, error) {
	rows, err := db.Query("select task_id, req_no, deadline, comment, estimated_work_hours from tasklist_table where  principal = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contents := make([]*sendContent, 0)
	for rows.Next() {
		content := &sendContent{}
		err = rows.Scan(&content.taskId, &content.reqNo, &content.deadline, &content.comment, &content.taskTime)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}
