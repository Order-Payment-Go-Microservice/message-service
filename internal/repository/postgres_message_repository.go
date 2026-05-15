package repository

import (
	"database/sql"
	"fmt"

	"github.com/Order-Payment-Go-Microservice/message-service/internal/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type MessageRepository interface {
	Save(msg *model.Message) error
	GetByID(id uuid.UUID) (*model.Message, error)
	GetByChatID(chatID uuid.UUID, limit, offset int) ([]*model.Message, error)
	UpdateStatus(id uuid.UUID, isRead, isDelivered bool) error
	Delete(id uuid.UUID) error
	UpdateContent(id uuid.UUID, content string) error
	Search(chatID uuid.UUID, query string) ([]*model.Message, error)
}

type postgresMessageRepository struct {
	db *sql.DB
}

func NewPostgresMessageRepository(db *sql.DB) MessageRepository {
	return &postgresMessageRepository{db: db}
}

func (r *postgresMessageRepository) Save(msg *model.Message) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO messages (id, chat_id, sender_id, receiver_id, content, message_type, is_read, is_delivered, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tx.Exec(query,
		msg.ID, msg.ChatID, msg.SenderID, msg.ReceiverID,
		msg.Content, msg.MessageType, msg.IsRead, msg.IsDelivered,
		msg.CreatedAt, msg.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}

	return tx.Commit()
}

func (r *postgresMessageRepository) GetByID(id uuid.UUID) (*model.Message, error) {
	query := `SELECT id, chat_id, sender_id, receiver_id, content, message_type, is_read, is_delivered, created_at, updated_at 
              FROM messages WHERE id = $1`
	row := r.db.QueryRow(query, id)
	msg := &model.Message{}
	err := row.Scan(
		&msg.ID, &msg.ChatID, &msg.SenderID, &msg.ReceiverID,
		&msg.Content, &msg.MessageType, &msg.IsRead, &msg.IsDelivered,
		&msg.CreatedAt, &msg.UpdatedAt,
	)
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
		err := rows.Scan(
			&msg.ID, &msg.ChatID, &msg.SenderID, &msg.ReceiverID,
			&msg.Content, &msg.MessageType, &msg.IsRead, &msg.IsDelivered,
			&msg.CreatedAt, &msg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *postgresMessageRepository) UpdateStatus(id uuid.UUID, isRead, isDelivered bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE messages SET is_read = $1, is_delivered = $2, updated_at = NOW() WHERE id = $3`
	if _, err = tx.Exec(query, isRead, isDelivered, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *postgresMessageRepository) Delete(id uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`DELETE FROM messages WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *postgresMessageRepository) UpdateContent(id uuid.UUID, content string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE messages SET content = $1, updated_at = NOW() WHERE id = $2`
	if _, err = tx.Exec(query, content, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *postgresMessageRepository) Search(chatID uuid.UUID, query string) ([]*model.Message, error) {
	sqlQuery := `SELECT id, chat_id, sender_id, receiver_id, content, message_type, is_read, is_delivered, created_at, updated_at 
                 FROM messages WHERE chat_id = $1 AND content ILIKE '%' || $2 || '%' ORDER BY created_at DESC`
	rows, err := r.db.Query(sqlQuery, chatID, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		msg := &model.Message{}
		if err := rows.Scan(
			&msg.ID, &msg.ChatID, &msg.SenderID, &msg.ReceiverID,
			&msg.Content, &msg.MessageType, &msg.IsRead, &msg.IsDelivered,
			&msg.CreatedAt, &msg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
