package domain

import "time"

type TaskAssignment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
