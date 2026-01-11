package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/launchventures/team-task-hub-backend/internal/domain"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
)

// TaskRepository defines task data access operations
type TaskRepository interface {
	CreateTask(ctx context.Context, projectID, createdByID string, title, description, status, priority string, assigneeID *string, dueDate *time.Time) (*domain.Task, error)
	GetTaskByID(ctx context.Context, id string) (*domain.Task, error)
	ListTasksByProjectID(ctx context.Context, projectID string, limit, offset int, status, priority string) ([]domain.Task, int, error)
	ListTasksByAssignee(ctx context.Context, userID string, limit, offset int, status, priority string) ([]domain.Task, int, error)
	UpdateTask(ctx context.Context, id string, title, description, status, priority string, assigneeID *string, dueDate *time.Time) (*domain.Task, error)
	AssignTaskToUser(ctx context.Context, taskID, userID, assignedByID string) (*domain.TaskAssignment, error)
	UnassignTask(ctx context.Context, taskID string) error
	DeleteTask(ctx context.Context, id string) error
}

type taskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) TaskRepository {
	return &taskRepository{db: db}
}

// CreateTask creates a new task
func (r *taskRepository) CreateTask(ctx context.Context, projectID, createdByID string, title, description, status, priority string, assigneeID *string, dueDate *time.Time) (*domain.Task, error) {
	taskID := uuid.New().String()
	var assignedByID *string
	if createdByID != "" {
		assignedByID = &createdByID
	}

	const insertQuery = `
		INSERT INTO tasks (id, project_id, title, description, status, priority, assignee_id, assigned_by_id, created_by_id, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`

	_, err := r.db.Exec(ctx, insertQuery, taskID, projectID, title, description, status, priority, assigneeID, assignedByID, createdByID, dueDate)
	if err != nil {
		return nil, apperrors.NewDatabaseError("failed to create task", err)
	}

	// Fetch the created task with user objects populated using GetTaskByID
	// which already has the logic to populate Assignee, AssignedBy, and CreatedBy
	return r.GetTaskByID(ctx, taskID)
}

// GetTaskByID retrieves a task by ID
func (r *taskRepository) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	const query = `
		SELECT t.id, t.project_id, t.assignee_id, t.assigned_by_id, t.created_by_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
		       u.id, u.email, u.name, ab.id, ab.email, ab.name, cb.id, cb.email, cb.name
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		LEFT JOIN users ab ON t.assigned_by_id = ab.id
		LEFT JOIN users cb ON t.created_by_id = cb.id
		WHERE t.id = $1
	`

	task := &domain.Task{}
	var userID *string
	var userEmail *string
	var userName *string
	var assignedByUserID *string
	var assignedByUserEmail *string
	var assignedByUserName *string
	var createdByUserID *string
	var createdByUserEmail *string
	var createdByUserName *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.ProjectID,
		&task.AssigneeID,
		&task.AssignedByID,
		&task.CreatedByID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&userID,
		&userEmail,
		&userName,
		&assignedByUserID,
		&assignedByUserEmail,
		&assignedByUserName,
		&createdByUserID,
		&createdByUserEmail,
		&createdByUserName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrTaskNotFound, "task not found")
		}
		return nil, apperrors.NewDatabaseError("failed to get task", err)
	}

	// Populate assignee if available
	if userID != nil && userEmail != nil {
		task.Assignee = &domain.User{
			ID:    *userID,
			Email: *userEmail,
			Name:  *userName,
		}
	}

	// Populate assigned_by user if available
	if assignedByUserID != nil && assignedByUserEmail != nil {
		task.AssignedBy = &domain.User{
			ID:    *assignedByUserID,
			Email: *assignedByUserEmail,
			Name:  *assignedByUserName,
		}
	}

	// Populate created_by user if available
	if createdByUserID != nil && createdByUserEmail != nil {
		task.CreatedBy = &domain.User{
			ID:    *createdByUserID,
			Email: *createdByUserEmail,
			Name:  *createdByUserName,
		}
	}

	return task, nil
}

// ListTasksByProjectID retrieves all tasks for a project with optional filters
func (r *taskRepository) ListTasksByProjectID(ctx context.Context, projectID string, limit, offset int, status, priority string) ([]domain.Task, int, error) {
	// Build dynamic query with optional filters
	whereClause := "WHERE t.project_id = $1"
	args := []interface{}{projectID}

	if status != "" {
		whereClause += " AND t.status = $" + strconv.Itoa(len(args)+1)
		args = append(args, status)
	}

	if priority != "" {
		whereClause += " AND t.priority = $" + strconv.Itoa(len(args)+1)
		args = append(args, priority)
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM tasks t " + whereClause
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("failed to count tasks", err)
	}

	// Get paginated results with user data
	query := `SELECT t.id, t.project_id, t.assignee_id, t.assigned_by_id, t.created_by_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
	       u.id, u.email, u.name, ab.id, ab.email, ab.name, cb.id, cb.email, cb.name
	FROM tasks t
	LEFT JOIN users u ON t.assignee_id = u.id
	LEFT JOIN users ab ON t.assigned_by_id = ab.id
	LEFT JOIN users cb ON t.created_by_id = cb.id
	` + whereClause
	query += " ORDER BY t.id DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("failed to list tasks", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var t domain.Task
		var userID *string
		var userEmail *string
		var userName *string
		var assignedByUserID *string
		var assignedByUserEmail *string
		var assignedByUserName *string
		var createdByUserID *string
		var createdByUserEmail *string
		var createdByUserName *string

		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.AssigneeID,
			&t.AssignedByID,
			&t.CreatedByID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
			&userID,
			&userEmail,
			&userName,
			&assignedByUserID,
			&assignedByUserEmail,
			&assignedByUserName,
			&createdByUserID,
			&createdByUserEmail,
			&createdByUserName,
		)
		if err != nil {
			return nil, 0, apperrors.NewDatabaseError("failed to scan task", err)
		}

		// Populate assignee if available
		if userID != nil && userEmail != nil {
			t.Assignee = &domain.User{
				ID:    *userID,
				Email: *userEmail,
				Name:  *userName,
			}
		}

		// Populate assigned_by user if available
		if assignedByUserID != nil && assignedByUserEmail != nil {
			t.AssignedBy = &domain.User{
				ID:    *assignedByUserID,
				Email: *assignedByUserEmail,
				Name:  *assignedByUserName,
			}
		}

		// Populate created_by user if available
		if createdByUserID != nil && createdByUserEmail != nil {
			t.CreatedBy = &domain.User{
				ID:    *createdByUserID,
				Email: *createdByUserEmail,
				Name:  *createdByUserName,
			}
		}

		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.NewDatabaseError("error iterating tasks", err)
	}

	return tasks, total, nil
}

// ListTasksByAssignee retrieves all tasks assigned to a user
func (r *taskRepository) ListTasksByAssignee(ctx context.Context, userID string, limit, offset int, status, priority string) ([]domain.Task, int, error) {
	query := `
		SELECT t.id, t.project_id, t.assignee_id, t.assigned_by_id, t.created_by_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
		       u.id, u.email, u.name, ab.id, ab.email, ab.name, cb.id, cb.email, cb.name
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		LEFT JOIN users ab ON t.assigned_by_id = ab.id
		LEFT JOIN users cb ON t.created_by_id = cb.id
		WHERE t.assignee_id = $1
	`
	args := []interface{}{userID}

	if status != "" {
		query += " AND t.status = $" + strconv.Itoa(len(args)+1)
		args = append(args, status)
	}

	if priority != "" {
		query += " AND t.priority = $" + strconv.Itoa(len(args)+1)
		args = append(args, priority)
	}

	countQuery := `SELECT COUNT(*) FROM tasks WHERE assignee_id = $1`
	countArgs := []interface{}{userID}
	if status != "" {
		countQuery += " AND status = $" + strconv.Itoa(len(countArgs)+1)
		countArgs = append(countArgs, status)
	}
	if priority != "" {
		countQuery += " AND priority = $" + strconv.Itoa(len(countArgs)+1)
		countArgs = append(countArgs, priority)
	}

	var count int
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&count); err != nil {
		return nil, 0, err
	}

	query += " ORDER BY t.id DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var t domain.Task
		var assigneeUserID *string
		var userEmail *string
		var userName *string
		var assignedByUserID *string
		var assignedByUserEmail *string
		var assignedByUserName *string
		var createdByUserID *string
		var createdByUserEmail *string
		var createdByUserName *string

		if err := rows.Scan(
			&t.ID, &t.ProjectID, &t.AssigneeID, &t.AssignedByID, &t.CreatedByID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
			&assigneeUserID, &userEmail, &userName, &assignedByUserID, &assignedByUserEmail, &assignedByUserName, &createdByUserID, &createdByUserEmail, &createdByUserName,
		); err != nil {
			return nil, 0, err
		}

		// Populate assignee if available
		if assigneeUserID != nil && userEmail != nil {
			t.Assignee = &domain.User{
				ID:    *assigneeUserID,
				Email: *userEmail,
				Name:  *userName,
			}
		}

		// Populate assigned_by user if available
		if assignedByUserID != nil && assignedByUserEmail != nil {
			t.AssignedBy = &domain.User{
				ID:    *assignedByUserID,
				Email: *assignedByUserEmail,
				Name:  *assignedByUserName,
			}
		}

		// Populate created_by user if available
		if createdByUserID != nil && createdByUserEmail != nil {
			t.CreatedBy = &domain.User{
				ID:    *createdByUserID,
				Email: *createdByUserEmail,
				Name:  *createdByUserName,
			}
		}

		tasks = append(tasks, t)
	}

	return tasks, count, nil
}

// UpdateTask updates a task
func (r *taskRepository) UpdateTask(ctx context.Context, id string, title, description, status, priority string, assigneeID *string, dueDate *time.Time) (*domain.Task, error) {
	// Update the task
	if assigneeID != nil && dueDate != nil {
		_, err := r.db.Exec(ctx, `
			UPDATE tasks
			SET title = $1, description = $2, status = $3, priority = $4, assignee_id = $5, due_date = $6, updated_at = NOW()
			WHERE id = $7
		`, title, description, status, priority, *assigneeID, *dueDate, id)
		if err != nil {
			return nil, apperrors.NewDatabaseError("failed to update task", err)
		}
	} else if assigneeID != nil {
		_, err := r.db.Exec(ctx, `
			UPDATE tasks
			SET title = $1, description = $2, status = $3, priority = $4, assignee_id = $5, updated_at = NOW()
			WHERE id = $6
		`, title, description, status, priority, *assigneeID, id)
		if err != nil {
			return nil, apperrors.NewDatabaseError("failed to update task", err)
		}
	} else if dueDate != nil {
		_, err := r.db.Exec(ctx, `
			UPDATE tasks
			SET title = $1, description = $2, status = $3, priority = $4, due_date = $5, updated_at = NOW()
			WHERE id = $6
		`, title, description, status, priority, *dueDate, id)
		if err != nil {
			return nil, apperrors.NewDatabaseError("failed to update task", err)
		}
	} else {
		_, err := r.db.Exec(ctx, `
			UPDATE tasks
			SET title = $1, description = $2, status = $3, priority = $4, updated_at = NOW()
			WHERE id = $5
		`, title, description, status, priority, id)
		if err != nil {
			return nil, apperrors.NewDatabaseError("failed to update task", err)
		}
	}

	// Fetch and return the updated task with assignee, assigned_by, and created_by
	const selectQuery = `
		SELECT t.id, t.project_id, t.assignee_id, t.assigned_by_id, t.created_by_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
		       u.id, u.email, ab.id, ab.email, cb.id, cb.email
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		LEFT JOIN users ab ON t.assigned_by_id = ab.id
		LEFT JOIN users cb ON t.created_by_id = cb.id
		WHERE t.id = $1
	`

	task := &domain.Task{}
	var userID *string
	var userEmail *string
	var assignedByUserID *string
	var assignedByUserEmail *string
	var createdByUserID *string
	var createdByUserEmail *string

	err := r.db.QueryRow(ctx, selectQuery, id).Scan(
		&task.ID,
		&task.ProjectID,
		&task.AssigneeID,
		&task.AssignedByID,
		&task.CreatedByID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&userID,
		&userEmail,
		&assignedByUserID,
		&assignedByUserEmail,
		&createdByUserID,
		&createdByUserEmail,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrTaskNotFound, "task not found")
		}
		return nil, apperrors.NewDatabaseError("failed to fetch updated task", err)
	}

	// Populate assignee if available
	if userID != nil && userEmail != nil {
		task.Assignee = &domain.User{
			ID:    *userID,
			Email: *userEmail,
		}
	}

	// Populate assigned_by user if available
	if assignedByUserID != nil && assignedByUserEmail != nil {
		task.AssignedBy = &domain.User{
			ID:    *assignedByUserID,
			Email: *assignedByUserEmail,
		}
	}

	// Populate created_by user if available
	if createdByUserID != nil && createdByUserEmail != nil {
		task.CreatedBy = &domain.User{
			ID:    *createdByUserID,
			Email: *createdByUserEmail,
		}
	}

	return task, nil
}

// AssignTaskToUser assigns a task to a user by updating assignee_id and assigned_by_id
func (r *taskRepository) AssignTaskToUser(ctx context.Context, taskID, userID, assignedByID string) (*domain.TaskAssignment, error) {
	const query = `
		UPDATE tasks
		SET assignee_id = $2, assigned_by_id = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id
	`

	assignment := &domain.TaskAssignment{}
	var taskIDReturned string
	err := r.db.QueryRow(ctx, query, taskID, userID, assignedByID).Scan(&taskIDReturned)

	if err != nil {
		return nil, apperrors.NewDatabaseError("failed to assign task", err)
	}

	assignment.TaskID = taskIDReturned
	assignment.UserID = userID
	return assignment, nil
}

// UnassignTask clears assignee and assigned_by for a task
func (r *taskRepository) UnassignTask(ctx context.Context, taskID string) error {
	const query = `
		UPDATE tasks
		SET assignee_id = NULL, assigned_by_id = NULL, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, taskID)
	if err != nil {
		return apperrors.NewDatabaseError("failed to unassign task", err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.NewNotFoundError(apperrors.ErrTaskNotFound, "task not found")
	}

	return nil
}

// DeleteTask deletes a task
func (r *taskRepository) DeleteTask(ctx context.Context, id string) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return apperrors.NewDatabaseError("failed to delete task", err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.NewNotFoundError(apperrors.ErrTaskNotFound, "task not found")
	}

	return nil
}
