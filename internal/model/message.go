package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          uuid.UUID `json:"id"`
	ChatID      uuid.UUID `json:"chat_id"`
	SenderID    uuid.UUID `json:"sender_id"`
	ReceiverID  uuid.UUID `json:"receiver_id"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	IsRead      bool      `json:"is_read"`
	IsDelivered bool      `json:"is_delivered"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SendMessageRequest struct {
	ChatID      string `json:"chat_id" binding:"required"`
	SenderID    string `json:"sender_id" binding:"required"`
	ReceiverID  string `json:"receiver_id" binding:"required"`
	Content     string `json:"content" binding:"required"`
	MessageType string `json:"message_type"`
}
