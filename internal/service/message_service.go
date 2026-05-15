package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Order-Payment-Go-Microservice/message-service/internal/model"
	"github.com/Order-Payment-Go-Microservice/message-service/internal/repository"

	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type MessageService interface {
	SendMessage(chatID, senderID, receiverID uuid.UUID, content, msgType string) (*model.Message, error)
	GetChatHistory(chatID uuid.UUID, limit, offset int) ([]model.Message, error)
	MarkRead(id uuid.UUID) error
	MarkDelivered(id uuid.UUID) error
	DeleteMessage(id uuid.UUID) error
	EditMessage(id uuid.UUID, content string) (*model.Message, error)
	SearchMessages(chatID uuid.UUID, query string) ([]model.Message, error)
	GetMessage(id uuid.UUID) (*model.Message, error)
}

type messageService struct {
	repo               repository.MessageRepository
	notificationClient NotificationClient
	hub                *Hub
	natsConn           *nats.Conn
	redisClient        *redis.Client
}

func NewMessageService(
	repo repository.MessageRepository,
	nc NotificationClient,
	hub *Hub,
	natsConn *nats.Conn,
	rdb *redis.Client,
) MessageService {
	return &messageService{
		repo:               repo,
		notificationClient: nc,
		hub:                hub,
		natsConn:           natsConn,
		redisClient:        rdb,
	}
}

func (s *messageService) SendMessage(chatID, senderID, receiverID uuid.UUID, content, msgType string) (*model.Message, error) {
	if msgType == "" {
		msgType = "text"
	}

	now := time.Now()
	msg := &model.Message{
		ID:          uuid.New(),
		ChatID:      chatID,
		SenderID:    senderID,
		ReceiverID:  receiverID,
		Content:     content,
		MessageType: msgType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Save(msg); err != nil {
		return nil, err
	}

	s.hub.BroadcastMessage(msg)

	if s.notificationClient != nil {
		if err := s.notificationClient.SendPushNotification(receiverID.String(), content); err != nil {
			log.Printf("[Notification] gRPC call failed: %v", err)
		}
	}

	event := map[string]string{
		"user_id": receiverID.String(),
		"title":   "Новое сообщение",
		"message": content,
		"type":    "push",
	}
	eventData, _ := json.Marshal(event)
	if s.natsConn != nil {
		if err := s.natsConn.Publish("notifications", eventData); err != nil {
			log.Printf("[NATS] publish failed: %v", err)
		} else {
			log.Println("[NATS] Event published to notifications subject")
		}
	}

	s.invalidateCache(chatID)
	return msg, nil
}

func (s *messageService) GetChatHistory(chatID uuid.UUID, limit, offset int) ([]model.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	ctx := context.Background()
	cacheKey := "chat_history:" + chatID.String()

	if s.redisClient != nil {
		cached, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var messages []model.Message
			if err := json.Unmarshal([]byte(cached), &messages); err == nil {
				log.Println("[Redis] Chat history loaded from cache")
				return messages, nil
			}
		}
	}

	ptrs, err := s.repo.GetByChatID(chatID, limit, offset)
	if err != nil {
		return nil, err
	}

	messages := make([]model.Message, len(ptrs))
	for i, m := range ptrs {
		messages[i] = *m
	}

	if s.redisClient != nil {
		data, _ := json.Marshal(messages)
		s.redisClient.Set(ctx, cacheKey, data, cacheTTL)
	}

	return messages, nil
}

func (s *messageService) MarkRead(id uuid.UUID) error {
	if err := s.repo.UpdateStatus(id, true, false); err != nil {
		return err
	}
	s.invalidateCacheForMessage(id)
	return nil
}

func (s *messageService) MarkDelivered(id uuid.UUID) error {
	return s.repo.UpdateStatus(id, false, true)
}

func (s *messageService) DeleteMessage(id uuid.UUID) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.invalidateCacheForMessage(id)
	return nil
}

func (s *messageService) EditMessage(id uuid.UUID, content string) (*model.Message, error) {
	if err := s.repo.UpdateContent(id, content); err != nil {
		return nil, err
	}
	msg, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	s.invalidateCache(msg.ChatID)
	return msg, nil
}

func (s *messageService) SearchMessages(chatID uuid.UUID, query string) ([]model.Message, error) {
	ptrs, err := s.repo.Search(chatID, query)
	if err != nil {
		return nil, err
	}
	messages := make([]model.Message, len(ptrs))
	for i, m := range ptrs {
		messages[i] = *m
	}
	return messages, nil
}

func (s *messageService) GetMessage(id uuid.UUID) (*model.Message, error) {
	return s.repo.GetByID(id)
}

func (s *messageService) invalidateCache(chatID uuid.UUID) {
	if s.redisClient == nil {
		return
	}
	ctx := context.Background()
	s.redisClient.Del(ctx, "chat_history:"+chatID.String())
}

func (s *messageService) invalidateCacheForMessage(id uuid.UUID) {
	msg, err := s.repo.GetByID(id)
	if err != nil {
		return
	}
	s.invalidateCache(msg.ChatID)
}
