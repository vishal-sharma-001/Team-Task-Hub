package handler

import (
	"net/http"

	"github.com/launchventures/team-task-hub-backend/internal/domain"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
)

// Response types for API
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type PaginatedResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Pages   int         `json:"pages"`
	Message string      `json:"message,omitempty"`
}

// DTO for signup/login requests
type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
	Name string `json:"name" validate:"max=255"`
}

type AuthResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

// DTO for project requests
type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=1000"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=1000"`
}

// DTO for task requests
type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=200"`
	Description string `json:"description" validate:"max=2000"`
	Priority    string `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH"`
	AssigneeID  *int   `json:"assignee_id"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=3,max=200"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	Status      *string `json:"status" validate:"omitempty,oneof=OPEN IN_PROGRESS DONE"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=LOW MEDIUM HIGH"`
	AssigneeID  *int    `json:"assignee_id"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=OPEN IN_PROGRESS DONE"`
}

type UpdateTaskPriorityRequest struct {
	Priority string `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH"`
}

type UpdateTaskAssigneeRequest struct {
	AssigneeID *int `json:"assignee_id"`
}

type AssignTaskRequest struct {
	UserID int `json:"user_id" validate:"required,min=1"`
}

// DTO for comment requests
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=3000"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=3000"`
}

// ErrorToStatusCode maps AppError to HTTP status codes
func ErrorToStatusCode(err error) int {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.StatusCode()
	}
	return http.StatusInternalServerError
}

// NewErrorResponse creates an error response
func NewErrorResponse(err error) ErrorResponse {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return ErrorResponse{
			Status:  "error",
			Error:   string(appErr.Code),
			Message: appErr.Message,
			Code:    string(appErr.Code),
		}
	}
	return ErrorResponse{
		Status:  "error",
		Error:   "InternalServerError",
		Message: "An unexpected error occurred",
	}
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}, message string) SuccessResponse {
	return SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, total, page, pageSize int, message string) PaginatedResponse {
	pages := (total + pageSize - 1) / pageSize
	if pages == 0 {
		pages = 1
	}
	return PaginatedResponse{
		Status:  "success",
		Data:    data,
		Total:   total,
		Page:    page,
		Pages:   pages,
		Message: message,
	}
}
