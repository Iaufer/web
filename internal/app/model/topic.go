package model

import "time"

type Topic struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	TopicName   string    `json:"TopicName"`
	Description string    `json:"description,omitempty"`
	Content     string    `json:"content,omitempty"`
	Visibility  bool      `json:"visibility"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
