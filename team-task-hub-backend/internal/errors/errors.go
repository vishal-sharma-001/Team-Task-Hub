package errors

import "fmt"

type ErrorCode string

const (
	// Validation errors
	ErrInvalidInput    ErrorCode = "invalid_input"
	ErrInvalidEmail    ErrorCode = "invalid_email"
	ErrWeakPassword    ErrorCode = "weak_password"
	ErrEmptyTitle      ErrorCode = "empty_title"
	ErrEmptyName       ErrorCode = "empty_name"
	ErrEmptyContent    ErrorCode = "empty_content"
	ErrInvalidStatus   ErrorCode = "invalid_status"
	ErrInvalidPriority ErrorCode = "invalid_priority"

	// Authentication/Authorization errors
	ErrUnauthorized    ErrorCode = "unauthorized"
	ErrForbidden       ErrorCode = "forbidden"
	ErrInvalidToken    ErrorCode = "invalid_token"
	ErrTokenExpired    ErrorCode = "token_expired"
	ErrInvalidPassword ErrorCode = "invalid_password"

	// Resource errors
	ErrUserNotFound    ErrorCode = "user_not_found"
	ErrProjectNotFound ErrorCode = "project_not_found"
	ErrTaskNotFound    ErrorCode = "task_not_found"
	ErrCommentNotFound ErrorCode = "comment_not_found"

	// Conflict errors
	ErrEmailExists       ErrorCode = "email_already_exists"
	ErrInvalidTransition ErrorCode = "invalid_status_transition"

	// Database/Server errors
	ErrInternal      ErrorCode = "internal_server_error"
	ErrDatabaseError ErrorCode = "database_error"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// HTTP Status Code mapping
func (e *AppError) StatusCode() int {
	switch e.Code {
	case ErrInvalidEmail, ErrWeakPassword, ErrEmptyTitle, ErrEmptyName, ErrEmptyContent, ErrInvalidStatus, ErrInvalidPriority:
		return 400
	case ErrUnauthorized, ErrInvalidToken, ErrTokenExpired, ErrInvalidPassword:
		return 401
	case ErrForbidden:
		return 403
	case ErrUserNotFound, ErrProjectNotFound, ErrTaskNotFound, ErrCommentNotFound:
		return 404
	case ErrEmailExists, ErrInvalidTransition:
		return 409
	default:
		return 500
	}
}

// Constructor functions
func NewValidationError(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewAuthError(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewNotFoundError(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewConflictError(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{Code: ErrInternal, Message: message, Err: err}
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{Code: ErrDatabaseError, Message: message, Err: err}
}
