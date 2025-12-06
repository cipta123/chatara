package repository

import (
	"database/sql"
	"errors"

	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/pkg/database"
)

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (r *MessageRepository) Create(msg *model.Message) error {
	query := `
		INSERT INTO messages (conversation_id, sender_id, content, type, media_url, reply_to_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := database.DB.QueryRow(
		query,
		msg.ConversationID,
		msg.SenderID,
		msg.Content,
		msg.Type,
		msg.MediaURL,
		msg.ReplyToID,
		msg.Status,
	).Scan(&msg.ID, &msg.CreatedAt, &msg.UpdatedAt)

	return err
}

func (r *MessageRepository) GetByID(id int) (*model.Message, error) {
	msg := &model.Message{}
	var replyToID sql.NullInt64
	var deletedAt sql.NullTime
	query := `
		SELECT id, conversation_id, sender_id, content, type, media_url, reply_to_id, status, created_at, updated_at, deleted_at
		FROM messages
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := database.DB.QueryRow(query, id).Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.SenderID,
		&msg.Content,
		&msg.Type,
		&msg.MediaURL,
		&replyToID,
		&msg.Status,
		&msg.CreatedAt,
		&msg.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}

	if replyToID.Valid {
		replyID := int(replyToID.Int64)
		msg.ReplyToID = &replyID
	}

	return msg, nil
}

func (r *MessageRepository) GetByConversation(conversationID int, limit, offset int) ([]*model.Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, content, type, media_url, reply_to_id, status, created_at, updated_at, deleted_at
		FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := database.DB.Query(query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		msg := &model.Message{}
		var replyToID sql.NullInt64
		var deletedAt sql.NullTime
		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.Content,
			&msg.Type,
			&msg.MediaURL,
			&replyToID,
			&msg.Status,
			&msg.CreatedAt,
			&msg.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if replyToID.Valid {
			replyID := int(replyToID.Int64)
			msg.ReplyToID = &replyID
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *MessageRepository) GetLastMessage(conversationID int) (*model.Message, error) {
	msg := &model.Message{}
	var replyToID sql.NullInt64
	var deletedAt sql.NullTime
	query := `
		SELECT id, conversation_id, sender_id, content, type, media_url, reply_to_id, status, created_at, updated_at, deleted_at
		FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := database.DB.QueryRow(query, conversationID).Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.SenderID,
		&msg.Content,
		&msg.Type,
		&msg.MediaURL,
		&replyToID,
		&msg.Status,
		&msg.CreatedAt,
		&msg.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No messages yet
		}
		return nil, err
	}

	if replyToID.Valid {
		replyID := int(replyToID.Int64)
		msg.ReplyToID = &replyID
	}

	return msg, nil
}

func (r *MessageRepository) UpdateStatus(id int, status model.MessageStatus) error {
	query := `UPDATE messages SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := database.DB.Exec(query, status, id)
	return err
}

func (r *MessageRepository) MarkAsDelivered(conversationID int, userID int) error {
	query := `
		UPDATE messages
		SET status = 'delivered'
		WHERE conversation_id = $1
		AND sender_id != $2
		AND status = 'sent'
	`
	_, err := database.DB.Exec(query, conversationID, userID)
	return err
}

func (r *MessageRepository) MarkAsRead(conversationID int, userID int) error {
	query := `
		UPDATE messages
		SET status = 'read'
		WHERE conversation_id = $1
		AND sender_id != $2
		AND status IN ('sent', 'delivered')
		AND deleted_at IS NULL
	`
	_, err := database.DB.Exec(query, conversationID, userID)
	return err
}

// Delete performs a soft delete by setting deleted_at timestamp
func (r *MessageRepository) Delete(id int, userID int) error {
	query := `
		UPDATE messages
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND sender_id = $2 AND deleted_at IS NULL
	`
	result, err := database.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("message not found or unauthorized")
	}
	
	return nil
}


