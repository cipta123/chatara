package repository

import (
	"database/sql"
	"errors"

	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/pkg/database"
)

type ConversationRepository struct{}

func NewConversationRepository() *ConversationRepository {
	return &ConversationRepository{}
}

func (r *ConversationRepository) Create(conv *model.Conversation) error {
	query := `
		INSERT INTO conversations (type, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := database.DB.QueryRow(query, conv.Type).Scan(
		&conv.ID,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)
	return err
}

func (r *ConversationRepository) GetByID(id int) (*model.Conversation, error) {
	conv := &model.Conversation{}
	query := `
		SELECT id, type, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`
	err := database.DB.QueryRow(query, id).Scan(
		&conv.ID,
		&conv.Type,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("conversation not found")
		}
		return nil, err
	}

	return conv, nil
}

func (r *ConversationRepository) FindDirectConversation(userID1, userID2 int) (*model.Conversation, error) {
	conv := &model.Conversation{}
	query := `
		SELECT DISTINCT c.id, c.type, c.created_at, c.updated_at
		FROM conversations c
		INNER JOIN participants p1 ON c.id = p1.conversation_id
		INNER JOIN participants p2 ON c.id = p2.conversation_id
		WHERE c.type = 'direct'
		AND p1.user_id = $1 AND p2.user_id = $2
		AND (SELECT COUNT(*) FROM participants WHERE conversation_id = c.id) = 2
		LIMIT 1
	`
	err := database.DB.QueryRow(query, userID1, userID2).Scan(
		&conv.ID,
		&conv.Type,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No conversation found, not an error
		}
		return nil, err
	}

	return conv, nil
}

func (r *ConversationRepository) GetUserConversations(userID int) ([]*model.ConversationWithParticipant, error) {
	query := `
		SELECT 
			c.id,
			c.type,
			c.created_at,
			c.updated_at
		FROM conversations c
		INNER JOIN participants p ON c.id = p.conversation_id
		WHERE p.user_id = $1
		ORDER BY c.updated_at DESC
	`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*model.ConversationWithParticipant
	for rows.Next() {
		conv := &model.ConversationWithParticipant{}
		err := rows.Scan(
			&conv.ID,
			&conv.Type,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return []*model.ConversationWithParticipant{}, err
		}
		conversations = append(conversations, conv)
	}

	// Always return at least an empty slice, never nil
	if conversations == nil {
		conversations = []*model.ConversationWithParticipant{}
	}
	return conversations, nil
}

func (r *ConversationRepository) AddParticipant(conversationID, userID int) error {
	query := `
		INSERT INTO participants (conversation_id, user_id, joined_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (conversation_id, user_id) DO NOTHING
	`
	_, err := database.DB.Exec(query, conversationID, userID)
	return err
}

func (r *ConversationRepository) GetParticipants(conversationID int) ([]int, error) {
	query := `
		SELECT user_id
		FROM participants
		WHERE conversation_id = $1
	`
	rows, err := database.DB.Query(query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func (r *ConversationRepository) UpdateTimestamp(conversationID int) error {
	query := `UPDATE conversations SET updated_at = NOW() WHERE id = $1`
	_, err := database.DB.Exec(query, conversationID)
	return err
}


