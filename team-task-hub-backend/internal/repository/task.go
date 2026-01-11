package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/launchventures/team-task-hub-backend/internal/domain"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
)

// TaskRepository defines task data access operations
type TaskRepository interface {
	CreateTask(ctx context.Context, projectID, createdByID int, title, description, status, priority string, assigneeID *int, dueDate *time.Time) (*domain.Task, error)
	GetTaskByID(ctx context.Context, id int) (*domain.Task, error)
	ListTasksByProjectID(ctx context.Context, projectID, limit, offset int, status, priority string) ([]domain.Task, int, error)
	ListTasksByAssignee(ctx context.Context, userID, limit, offset int, status, priority string) ([]domain.Task, int, error)
	UpdateTask(ctx context.Context, id int, title, description, status, priority string, assigneeID *int, dueDate *time.Time) (*domain.Task, error)
	AssignTaskToUser(ctx context.Context, taskID, userID int) (*domain.TaskAssignment, error)
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) TaskRepository {
	return &taskRepository{db: db}
}

// CreateTask creates a new task
func (r *taskRepository) CreateTask(ctx context.Context, projectID, createdByID int, title, description, status, priority string, assigneeID *int, dueDate *time.Time) (*domain.Task, error) {
	const query = `
		INSERT INTO tasks (project_id, title, description, status, priority, assignee_id, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, project_id, assignee_id, title, description, status, priority, due_date, created_at, updated_at
	`

	task := &domain.Task{}
	err := r.db.QueryRow(ctx, query, projectID, title, description, status, priority, assigneeID, dueDate).Scan(
		&task.ID,
		&task.ProjectID,
		&task.AssigneeID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, apperrors.NewDatabaseError("failed to create task", err)
	}

	return task, nil
}

// GetTaskByID retrieves a task by ID
func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (*domain.Task, error) {
	const query = `
		SELECT t.id, t.project_id, t.assignee_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
		       u.id, u.email
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		WHERE t.id = $1
	`

	task := &domain.Task{}
	var userID *int
	var userEmail *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.ProjectID,
		&task.AssigneeID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&userID,
		&userEmail,
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
		}
	}

	return task, nil
}

// ListTasksByProjectID retrieves all tasks for a project with optional filters
func (r *taskRepository) ListTasksByProjectID(ctx context.Context, projectID, limit, offset int, status, priority string) ([]domain.Task, int, error) {
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
	query := `SELECT t.id, t.project_id, t.assignee_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
	       u.id, u.email
	FROM tasks t
	LEFT JOIN users u ON t.assignee_id = u.id
	` + whereClause
	query += " ORDER BY t.created_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("failed to list tasks", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var t domain.Task
		var userID *int
		var userEmail *string

		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.AssigneeID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
			&userID,
			&userEmail,
		)
		if err != nil {
			return nil, 0, apperrors.NewDatabaseError("failed to scan task", err)
		}

		// Populate assignee if available
		if userID != nil && userEmail != nil {
			t.Assignee = &domain.User{
				ID:    *userID,
				Email: *userEmail,
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
func (r *taskRepository) ListTasksByAssignee(ctx context.Context, userID, limit, offset int, status, priority string) ([]domain.Task, int, error) {
	query := `
		SELECT t.id, t.project_id, t.assignee_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
		       u.id, u.email
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
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
		countQuery += " AND t.status = $" + strconv.Itoa(len(countArgs)+1)
		countArgs = append(countArgs, status)
	}
	if priority != "" {
		countQuery += " AND t.priority = $" + strconv.Itoa(len(countArgs)+1)
		countArgs = append(countArgs, priority)
	}

	var count int
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&count); err != nil {
		return nil, 0, err
	}

	query += " ORDER BY t.updated_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var t domain.Task
		var assigneeUserID *int
		var userEmail *string

		if err := rows.Scan(
			&t.ID, &t.ProjectID, &t.AssigneeID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
			&assigneeUserID, &userEmail,
		); err != nil {
			return nil, 0, err
		}

		// Populate assignee if available
		if assigneeUserID != nil && userEmail != nil {
			t.Assignee = &domain.User{
				ID:    *assigneeUserID,
				Email: *userEmail,
			}
		}

		tasks = append(tasks, t)
	}

	return tasks, count, nil
}

// UpdateTask updates a task
func (r *taskRepository) UpdateTask(ctx context.Context, id int, title, description, status, priority string, assigneeID *int, dueDate *time.Time) (*domain.Task, error) {
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

	// Fetch and return the updated task with assignee
	const selectQuery = `
		SELECT t.id, t.project_id, t.assignee_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
		       u.id, u.email
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		WHERE t.id = $1
	`

	task := &domain.Task{}
	var userID *int
	var userEmail *string

	err := r.db.QueryRow(ctx, selectQuery, id).Scan(
		&task.ID,
		&task.ProjectID,
		&task.AssigneeID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&userID,
		&userEmail,
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

	return task, nil
}

// AssignTaskToUser assigns a task to a user by updating assignee_id
func (r *taskRepository) AssignTaskToUser(ctx context.Context, taskID, userID int) (*domain.TaskAssignment, error) {
	const query = `
		UPDATE tasks
		SET assignee_id = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id
	`

	assignment := &domain.TaskAssignment{}
	var taskIDReturned int
	err := r.db.QueryRow(ctx, query, taskID, userID).Scan(&taskIDReturned)

	if err != nil {
		return nil, apperrors.NewDatabaseError("failed to assign task", err)
	}

	assignment.TaskID = taskIDReturned
	assignment.UserID = userID
	return assignment, nil
}

// DeleteTask deletes a task
func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
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
