package repository

import (
	"database/sql"
	"github.com/Order-Payment-Go-Microservice/message-service/internal/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type MessageRepository interface {
	Save(msg *model.Message) error
	GetByID(id uuid.UUID) (*model.Message, error)
	GetByChatID(chatID uuid.UUID, limit, offset int) ([]*model.Message, error)
	UpdateStatus(id uuid.UUID, isRead, isDelivered bool) error
}

type postgresMessageRepository struct {
	db *sql.DB
}

func NewPostgresMessageRepository(db *sql.DB) MessageRepository {
	return &postgresMessageRepository{db: db}
}

func (r *postgresMessageRepository) Save(msg *model.Message) error {
	query := `INSERT INTO messages (id, chat_id, sender_id, receiver_id, content, message_type, is_read, is_delivered, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(query, msg.ID, msg.ChatID, msg.SenderID, msg.ReceiverID, msg.Content, msg.MessageType, msg.IsRead, msg.IsDelivered, msg.CreatedAt, msg.UpdatedAt)
	return err
}

func (r *postgresMessageRepository) GetByID(id uuid.UUID) (*model.Message, error) {
	query := `SELECT id, chat_id, sender_id, receiver_id, content, message_type, is_read, is_delivered, created_at, updated_at FROM messages WHERE id = $1`
	row := r.db.QueryRow(query, id)
	msg := &model.Message{}
	err := row.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.MessageType, &msg.IsRead, &msg.IsDelivered, &msg.CreatedAt, &msg.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (r *postgresMessageRepository) GetByChatID(chatID uuid.UUID, limit, offset int) ([]*model.Message, error) {
	query := `SELECT id, chat_id, sender_id, receiver_id, content, message_type, is_read, is_delivered, created_at, updated_at 
              FROM messages WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		msg := &model.Message{}
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.MessageType, &msg.IsRead, &msg.IsDelivered, &msg.CreatedAt, &msg.UpdatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *postgresMessageRepository) UpdateStatus(id uuid.UUID, isRead, isDelivered bool) error {
	query := `UPDATE messages SET is_read = $1, is_delivered = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(query, isRead, isDelivered, id)
	return err
}
