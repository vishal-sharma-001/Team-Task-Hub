package service

import (
	"context"
	"time"

	"github.com/launchventures/team-task-hub-backend/internal/domain"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
	"github.com/launchventures/team-task-hub-backend/internal/repository"
	"github.com/launchventures/team-task-hub-backend/internal/utils"
)

// TaskService defines task-related business logic operations
type TaskService interface {
	CreateTask(ctx context.Context, projectID, createdByID string, title, description, priority string, assigneeID *string, dueDate *time.Time) (*domain.Task, error)
	GetTask(ctx context.Context, id string) (*domain.Task, error)
	ListTasks(ctx context.Context, projectID string, page, pageSize int, status, priority string) ([]domain.Task, int, error)
	ListAssignedTasks(ctx context.Context, userID string, page, pageSize int, status, priority string) ([]domain.Task, int, error)
	UpdateTask(ctx context.Context, id string, title, description, status, priority string, assigneeID *string, dueDate *time.Time) (*domain.Task, error)
	AssignTask(ctx context.Context, taskID, userID, assignedByID string) error
	UnassignTask(ctx context.Context, taskID string) error
	DeleteTask(ctx context.Context, id string) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

// CreateTask creates a new task with validation
func (s *taskService) CreateTask(ctx context.Context, projectID, createdByID string, title, description, priority string, assigneeID *string, dueDate *time.Time) (*domain.Task, error) {
	// Validate task title
	if appErr := utils.ValidateTaskTitle(title); appErr != nil {
		return nil, appErr
	}

	// Validate description (max 2000 chars)
	if len(description) > 2000 {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "description must not exceed 2000 characters")
	}

	// Validate priority
	if appErr := utils.ValidatePriority(priority); appErr != nil {
		return nil, appErr
	}

	// Tasks are created with "OPEN" status by default
	status := "OPEN"

	// Create task in database
	task, err := s.taskRepo.CreateTask(ctx, projectID, createdByID, title, description, status, priority, assigneeID, dueDate)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetTask retrieves a task by ID
func (s *taskService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	if id == "" {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid task ID")
	}

	task, err := s.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// ListTasks retrieves all tasks for a project with optional filters and pagination
func (s *taskService) ListTasks(ctx context.Context, projectID string, page, pageSize int, status, priority string) ([]domain.Task, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Validate status if provided
	if status != "" {
		if appErr := utils.ValidateStatus(status); appErr != nil {
			return nil, 0, appErr
		}
	}

	// Validate priority if provided
	if priority != "" {
		if appErr := utils.ValidatePriority(priority); appErr != nil {
			return nil, 0, appErr
		}
	}

	offset := (page - 1) * pageSize

	tasks, total, err := s.taskRepo.ListTasksByProjectID(ctx, projectID, pageSize, offset, status, priority)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// ListAssignedTasks retrieves all tasks assigned to a user with optional filters and pagination
func (s *taskService) ListAssignedTasks(ctx context.Context, userID string, page, pageSize int, status, priority string) ([]domain.Task, int, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	tasks, total, err := s.taskRepo.ListTasksByAssignee(ctx, userID, pageSize, offset, status, priority)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// UpdateTask updates a task with validation
func (s *taskService) UpdateTask(ctx context.Context, id string, title, description, status, priority string, assigneeID *string, dueDate *time.Time) (*domain.Task, error) {
	// Get current task first to support partial updates
	currentTask, err := s.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Use current values if fields are empty (partial update)
	if title == "" {
		title = currentTask.Title
	} else {
		// Only validate title if it's being updated and is not empty
		if appErr := utils.ValidateTaskTitle(title); appErr != nil {
			return nil, appErr
		}
	}

	if description == "" {
		description = currentTask.Description
	} else {
		// Only validate description if it's being updated
		if len(description) > 2000 {
			return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "description must not exceed 2000 characters")
		}
	}

	if status == "" {
		status = currentTask.Status
	} else {
		// Only validate status if it's being updated
		if appErr := utils.ValidateStatus(status); appErr != nil {
			return nil, appErr
		}
	}

	if priority == "" {
		priority = currentTask.Priority
	} else {
		// Only validate priority if it's being updated
		if appErr := utils.ValidatePriority(priority); appErr != nil {
			return nil, appErr
		}
	}

	// Handle due_date - preserve current value if not provided
	var finalDueDate *time.Time
	if dueDate != nil {
		finalDueDate = dueDate
	} else {
		finalDueDate = currentTask.DueDate
	}

	// Update task in database
	task, err := s.taskRepo.UpdateTask(ctx, id, title, description, status, priority, assigneeID, finalDueDate)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// AssignTask assigns a task to a user
func (s *taskService) AssignTask(ctx context.Context, taskID, userID, assignedByID string) error {
	if taskID == "" || userID == "" || assignedByID == "" {
		return apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid task ID, user ID, or assigned by ID")
	}

	_, err := s.taskRepo.AssignTaskToUser(ctx, taskID, userID, assignedByID)
	if err != nil {
		return err
	}

	return nil
}

// UnassignTask clears the assignee and assigned_by for a task
func (s *taskService) UnassignTask(ctx context.Context, taskID string) error {
	if taskID == "" {
		return apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid task ID")
	}

	return s.taskRepo.UnassignTask(ctx, taskID)
}

// DeleteTask deletes a task
func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid task ID")
	}

	err := s.taskRepo.DeleteTask(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
