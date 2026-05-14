package service

import (
	"github.com/Order-Payment-Go-Microservice/message-service/internal/model"
	"github.com/Order-Payment-Go-Microservice/message-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type NotificationClient interface {
	SendPushNotification(recipientID, content string) error
}

type MessageService interface {
	SendMessage(req *model.SendMessageRequest) (*model.Message, error)
	GetHistory(chatID string, limit, offset int) ([]*model.Message, error)
	MarkRead(id string) error
}

type messageService struct {
	repo               repository.MessageRepository
	notificationClient NotificationClient
	hub                *Hub
}

func NewMessageService(repo repository.MessageRepository, nc NotificationClient, hub *Hub) MessageService {
	return &messageService{
		repo:               repo,
		notificationClient: nc,
		hub:                hub,
	}
}

func (s *messageService) SendMessage(req *model.SendMessageRequest) (*model.Message, error) {
	chatID, _ := uuid.Parse(req.ChatID)
	senderID, _ := uuid.Parse(req.SenderID)
	receiverID, _ := uuid.Parse(req.ReceiverID)

	msg := &model.Message{
		ID:          uuid.New(),
		ChatID:      chatID,
		SenderID:    senderID,
		ReceiverID:  receiverID,
		Content:     req.Content,
		MessageType: req.MessageType,
		IsRead:      false,
		IsDelivered: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Save(msg); err != nil {
		return nil, err
	}

	go s.notificationClient.SendPushNotification(req.ReceiverID, req.Content)

	s.hub.BroadcastMessage(msg)

	return msg, nil
}

func (s *messageService) GetHistory(chatIDStr string, limit, offset int) ([]*model.Message, error) {
	chatID, _ := uuid.Parse(chatIDStr)
	return s.repo.GetByChatID(chatID, limit, offset)
}

func (s *messageService) MarkRead(idStr string) error {
	id, _ := uuid.Parse(idStr)
	return s.repo.UpdateStatus(id, true, true)
}
