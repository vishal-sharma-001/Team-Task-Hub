package domain

import "time"

type Task struct {
	ID          int        `json:"id"`
	ProjectID   int        `json:"project_id"`
	AssigneeID  *int       `json:"assignee_id,omitempty"`
	Assignee    *User      `json:"assignee,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
