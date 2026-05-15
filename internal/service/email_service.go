package service

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendChatExport(to, chatID string, content string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	if host == "" || user == "" {
		log.Printf("[SMTP Message Service] Mock export sent for Chat %s to %s", chatID, to)
		return nil
	}

	auth := smtp.PlainAuth("", user, password, host)
	body := fmt.Sprintf("Subject: Chat Export %s\r\n\r\nHere is your chat history:\n%s", chatID, content)

	err := smtp.SendMail(host+":"+port, auth, user, []string{to}, []byte(body))
	return err
}
