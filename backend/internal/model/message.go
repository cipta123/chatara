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
	Status         MessageStatus `json:"status" db:"status"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
	
	// Joined fields
	Sender *User `json:"sender,omitempty"`
}

type SendMessageRequest struct {
	ConversationID int         `json:"conversation_id" binding:"required"`
	Content        string      `json:"content"`
	Type           MessageType `json:"type" binding:"required"`
	MediaURL       string      `json:"media_url,omitempty"`
}

type TypingIndicator struct {
	ConversationID int    `json:"conversation_id"`
	UserID         int    `json:"user_id"`
	Username       string `json:"username,omitempty"`
	IsTyping       bool   `json:"is_typing"`
}


