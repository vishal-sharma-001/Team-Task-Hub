package domain

import "time"

type Comment struct {
	ID          int       `json:"id"`
	TaskID      int       `json:"task_id"`
	UserID      int       `json:"user_id"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AuthorName  string    `json:"author_name,omitempty"`
	AuthorEmail string    `json:"author_email,omitempty"`
}
