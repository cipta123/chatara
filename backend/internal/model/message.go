package model

import (
	"time"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

type MessageStatus string

const (
	MessageStatusSent     MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead     MessageStatus = "read"
)

type Message struct {
	ID             int          `json:"id" db:"id"`
	ConversationID int          `json:"conversation_id" db:"conversation_id"`
	SenderID       int          `json:"sender_id" db:"sender_id"`
	Content        string       `json:"content" db:"content"`
	Type           MessageType  `json:"type" db:"type"`
	MediaURL       string       `json:"media_url,omitempty" db:"media_url"`
	ReplyToID      *int         `json:"reply_to_id,omitempty" db:"reply_to_id"`
	Status         MessageStatus `json:"status" db:"status"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
	
	// Joined fields
	Sender   *User    `json:"sender,omitempty"`
	ReplyTo  *Message `json:"reply_to,omitempty"` // The message being replied to
}

type SendMessageRequest struct {
	ConversationID int         `json:"conversation_id" binding:"required"`
	Content        string      `json:"content"`
	Type           MessageType `json:"type" binding:"required"`
	MediaURL       string      `json:"media_url,omitempty"`
	ReplyToID      *int        `json:"reply_to_id,omitempty"`
}

type TypingIndicator struct {
	ConversationID int    `json:"conversation_id"`
	UserID         int    `json:"user_id"`
	Username       string `json:"username,omitempty"`
	IsTyping       bool   `json:"is_typing"`
}


