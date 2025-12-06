package model

import (
	"time"
)

type Conversation struct {
	ID        int       `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"` // "direct" or "group"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Participant struct {
	ID             int       `json:"id" db:"id"`
	ConversationID int       `json:"conversation_id" db:"conversation_id"`
	UserID         int       `json:"user_id" db:"user_id"`
	JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
}

type ConversationWithParticipant struct {
	Conversation
	Participants []User `json:"participants"`
	LastMessage  *Message `json:"last_message,omitempty"`
	UnreadCount  int      `json:"unread_count,omitempty"`
}


