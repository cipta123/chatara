package model

import (
	"time"
)

type Media struct {
	ID             int       `json:"id" db:"id"`
	MessageID      int       `json:"message_id" db:"message_id"`
	FileName       string    `json:"file_name" db:"file_name"`
	FileSize       int64     `json:"file_size" db:"file_size"`
	MimeType       string    `json:"mime_type" db:"mime_type"`
	FilePath       string    `json:"file_path" db:"file_path"`
	URL            string    `json:"url" db:"url"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}


