package service

import (
	"testing"
	"time"

	"github.com/Order-Payment-Go-Microservice/message-service/internal/model"
	"github.com/google/uuid"
)

type mockMessageRepo struct {
	saved *model.Message
}

func (m *mockMessageRepo) Save(msg *model.Message) error {
	m.saved = msg
	return nil
}
func (m *mockMessageRepo) GetByID(id uuid.UUID) (*model.Message, error) { return m.saved, nil }
func (m *mockMessageRepo) GetByChatID(chatID uuid.UUID, limit, offset int) ([]*model.Message, error) {
	return []*model.Message{m.saved}, nil
}
func (m *mockMessageRepo) UpdateStatus(id uuid.UUID, isRead, isDelivered bool) error { return nil }
func (m *mockMessageRepo) Delete(id uuid.UUID) error                                   { return nil }
func (m *mockMessageRepo) UpdateContent(id uuid.UUID, content string) error            { return nil }
func (m *mockMessageRepo) Search(chatID uuid.UUID, query string) ([]*model.Message, error) {
	return nil, nil
}

type mockNotificationClient struct{}

func (m *mockNotificationClient) SendPushNotification(recipientID, content string) error {
	return nil
}

func TestSendMessage(t *testing.T) {
	repo := &mockMessageRepo{}
	svc := NewMessageService(repo, &mockNotificationClient{}, NewHub(), nil, nil)

	chatID := uuid.New()
	senderID := uuid.New()
	receiverID := uuid.New()

	msg, err := svc.SendMessage(chatID, senderID, receiverID, "Hello", "text")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	if msg.Content != "Hello" {
		t.Errorf("expected Hello, got %s", msg.Content)
	}
	if repo.saved == nil {
		t.Fatal("message was not saved")
	}
}

func TestGetChatHistory(t *testing.T) {
	repo := &mockMessageRepo{
		saved: &model.Message{
			ID: uuid.New(), ChatID: uuid.New(), Content: "Hi",
			CreatedAt: time.Now(), UpdatedAt: time.Now(),
		},
	}
	svc := NewMessageService(repo, nil, NewHub(), nil, nil)

	msgs, err := svc.GetChatHistory(repo.saved.ChatID, 10, 0)
	if err != nil {
		t.Fatalf("GetChatHistory failed: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
}
