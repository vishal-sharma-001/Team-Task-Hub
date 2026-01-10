package utils

import (
	"regexp"
	"strings"

	"github.com/launchventures/team-task-hub-backend/internal/errors"
)

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) *errors.AppError {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.NewValidationError(errors.ErrInvalidEmail, "email cannot be empty")
	}

	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.NewValidationError(errors.ErrInvalidEmail, "invalid email format")
	}

	if len(email) > 255 {
		return errors.NewValidationError(errors.ErrInvalidEmail, "email is too long")
	}

	return nil
}

// ValidatePassword checks if password meets minimum requirements
func ValidatePassword(password string) *errors.AppError {
	if len(password) < 6 {
		return errors.NewValidationError(errors.ErrWeakPassword, "password must be at least 6 characters")
	}

	if len(password) > 128 {
		return errors.NewValidationError(errors.ErrWeakPassword, "password is too long")
	}

	return nil
}

// ValidateProjectName checks if project name is valid
func ValidateProjectName(name string) *errors.AppError {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.NewValidationError(errors.ErrEmptyName, "project name cannot be empty")
	}

	if len(name) > 255 {
		return errors.NewValidationError(errors.ErrEmptyName, "project name is too long")
	}

	return nil
}

// ValidateTaskTitle checks if task title is valid
func ValidateTaskTitle(title string) *errors.AppError {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.NewValidationError(errors.ErrEmptyTitle, "task title cannot be empty")
	}

	if len(title) > 255 {
		return errors.NewValidationError(errors.ErrEmptyTitle, "task title is too long")
	}

	return nil
}

// ValidateTaskStatus checks if status is valid
func ValidateTaskStatus(status string) *errors.AppError {
	validStatuses := map[string]bool{
		"OPEN":        true,
		"IN_PROGRESS": true,
		"DONE":        true,
	}

	if !validStatuses[status] {
		return errors.NewValidationError(errors.ErrInvalidStatus, "invalid task status")
	}

	return nil
}

// ValidateStatus is an alias for ValidateTaskStatus
func ValidateStatus(status string) *errors.AppError {
	return ValidateTaskStatus(status)
}

// ValidateTaskPriority checks if priority is valid
func ValidateTaskPriority(priority string) *errors.AppError {
	validPriorities := map[string]bool{
		"LOW":    true,
		"MEDIUM": true,
		"HIGH":   true,
	}

	if !validPriorities[priority] {
		return errors.NewValidationError(errors.ErrInvalidPriority, "invalid task priority")
	}

	return nil
}

// ValidatePriority is an alias for ValidateTaskPriority
func ValidatePriority(priority string) *errors.AppError {
	return ValidateTaskPriority(priority)
}

// ValidateTaskStatusTransition checks if status transition is valid
func ValidateTaskStatusTransition(currentStatus, newStatus string) *errors.AppError {
	validTransitions := map[string]map[string]bool{
		"OPEN": {
			"IN_PROGRESS": true,
			"DONE":        true,
			"OPEN":        true,
		},
		"IN_PROGRESS": {
			"DONE":        true,
			"OPEN":        true,
			"IN_PROGRESS": true,
		},
		"DONE": {
			"OPEN":        true,
			"IN_PROGRESS": true,
			"DONE":        true,
		},
	}

	if transitions, exists := validTransitions[currentStatus]; !exists || !transitions[newStatus] {
		return errors.NewValidationError(errors.ErrInvalidTransition, "invalid status transition")
	}

	return nil
}

// ValidateCommentContent checks if comment content is valid
func ValidateCommentContent(content string) *errors.AppError {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.NewValidationError(errors.ErrEmptyContent, "comment content cannot be empty")
	}

	if len(content) > 5000 {
		return errors.NewValidationError(errors.ErrEmptyContent, "comment content is too long")
	}

	return nil
}
