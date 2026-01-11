package domain

import "time"

type TaskAssignment struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
