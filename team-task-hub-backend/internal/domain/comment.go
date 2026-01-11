package domain

import "time"

type Comment struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"task_id"`
	UserID      string    `json:"user_id"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AuthorName  string    `json:"author_name,omitempty"`
	AuthorEmail string    `json:"author_email,omitempty"`
}
