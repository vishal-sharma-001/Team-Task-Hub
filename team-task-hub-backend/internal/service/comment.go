package service

import (
	"context"
	"strings"

	"github.com/launchventures/team-task-hub-backend/internal/domain"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
	"github.com/launchventures/team-task-hub-backend/internal/repository"
)

// CommentService defines comment-related business logic operations
type CommentService interface {
	CreateComment(ctx context.Context, taskID, userID string, content string) (*domain.Comment, error)
	GetComment(ctx context.Context, id string) (*domain.Comment, error)
	ListComments(ctx context.Context, taskID string, page, pageSize int) ([]domain.Comment, int, error)
	ListRecentComments(ctx context.Context, page, pageSize int) ([]domain.Comment, int, error)
	UpdateComment(ctx context.Context, id string, content string) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id string) error
}

type commentService struct {
	commentRepo repository.CommentRepository
}

func NewCommentService(commentRepo repository.CommentRepository) CommentService {
	return &commentService{commentRepo: commentRepo}
}

// CreateComment creates a new comment with validation
func (s *commentService) CreateComment(ctx context.Context, taskID, userID string, content string) (*domain.Comment, error) {
	// Validate task ID and user ID
	if taskID == "" || userID == "" {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid task ID or user ID")
	}

	// Validate content
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "comment content cannot be empty")
	}

	if len(content) > 3000 {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "comment content must not exceed 3000 characters")
	}

	// Create comment in database
	comment, err := s.commentRepo.CreateComment(ctx, taskID, userID, content)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetComment retrieves a comment by ID
func (s *commentService) GetComment(ctx context.Context, id string) (*domain.Comment, error) {
	if id == "" {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid comment ID")
	}

	comment, err := s.commentRepo.GetCommentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// ListComments retrieves all comments for a task with pagination
func (s *commentService) ListComments(ctx context.Context, taskID string, page, pageSize int) ([]domain.Comment, int, error) {
	// Validate task ID
	if taskID == "" {
		return nil, 0, apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid task ID")
	}

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	comments, total, err := s.commentRepo.ListCommentsByTaskID(ctx, taskID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListRecentComments retrieves recent comments from all tasks with pagination
func (s *commentService) ListRecentComments(ctx context.Context, page, pageSize int) ([]domain.Comment, int, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	comments, total, err := s.commentRepo.ListRecentComments(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// UpdateComment updates a comment with validation
func (s *commentService) UpdateComment(ctx context.Context, id string, content string) (*domain.Comment, error) {
	// Validate comment ID
	if id == "" {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid comment ID")
	}

	// Validate content
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "comment content cannot be empty")
	}

	if len(content) > 3000 {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "comment content must not exceed 3000 characters")
	}

	// Update comment in database
	comment, err := s.commentRepo.UpdateComment(ctx, id, content)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// DeleteComment deletes a comment
func (s *commentService) DeleteComment(ctx context.Context, id string) error {
	if id == "" {
		return apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid comment ID")
	}

	err := s.commentRepo.DeleteComment(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
