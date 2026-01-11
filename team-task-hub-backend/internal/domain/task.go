package domain

import "time"

type Task struct {
	ID           string     `json:"id"`
	ProjectID    string     `json:"project_id"`
	AssigneeID   *string    `json:"assignee_id,omitempty"`
	Assignee     *User      `json:"assignee,omitempty"`
	AssignedByID *string    `json:"assigned_by_id,omitempty"`
	AssignedBy   *User      `json:"assigned_by,omitempty"`
	CreatedByID  *string    `json:"created_by_id,omitempty"`
	CreatedBy    *User      `json:"created_by,omitempty"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
