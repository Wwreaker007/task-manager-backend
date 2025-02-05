package notification

import (
	"context"
	"fmt"
	"log"
	"task-manager-backend/data/common"

	"gopkg.in/gomail.v2"
)

type EmailNotificationServiceManager struct{
	dialer 			*gomail.Dialer
}

func NewEmailNotificationServiceManager(dialer *gomail.Dialer) *EmailNotificationServiceManager{
	return &EmailNotificationServiceManager{
		dialer: dialer,
	}
}

func (nsm *EmailNotificationServiceManager) SendNotification(ctx context.Context, messages []common.NotificationPayload) {
	for _, msg := range messages {
		mailPayload := createEmailMessage(nsm.dialer.Username, msg.Recipient, msg.Subject, msg.Body)
		err := nsm.dialer.DialAndSend(mailPayload)
		if err != nil {
			log.Println("Unable to send email to recipient : ", msg.Recipient)
		}
	}
}

func createEmailMessage(sender string, recipient string, subject string, body string) *gomail.Message{
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", fmt.Sprintf("<p>%s</p> \n", body))
	return m;
}