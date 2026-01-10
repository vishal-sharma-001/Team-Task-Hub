package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/launchventures/team-task-hub-backend/internal/domain"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
)

// CommentRepository defines comment data access operations
type CommentRepository interface {
	CreateComment(ctx context.Context, taskID, userID int, content string) (*domain.Comment, error)
	GetCommentByID(ctx context.Context, id int) (*domain.Comment, error)
	ListCommentsByTaskID(ctx context.Context, taskID, limit, offset int) ([]domain.Comment, int, error)
	ListRecentComments(ctx context.Context, limit, offset int) ([]domain.Comment, int, error)
	UpdateComment(ctx context.Context, id int, content string) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id int) error
}

type commentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) CommentRepository {
	return &commentRepository{db: db}
}

// CreateComment creates a new comment
func (r *commentRepository) CreateComment(ctx context.Context, taskID, userID int, content string) (*domain.Comment, error) {
	const query = `
		INSERT INTO comments (task_id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, task_id, user_id, content, created_at, updated_at
	`

	comment := &domain.Comment{}
	err := r.db.QueryRow(ctx, query, taskID, userID, content).Scan(
		&comment.ID,
		&comment.TaskID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, apperrors.NewDatabaseError("failed to create comment", err)
	}

	return comment, nil
}

// GetCommentByID retrieves a comment by ID
func (r *commentRepository) GetCommentByID(ctx context.Context, id int) (*domain.Comment, error) {
	const query = `
		SELECT id, task_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE id = $1
	`

	comment := &domain.Comment{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&comment.ID,
		&comment.TaskID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrCommentNotFound, "comment not found")
		}
		return nil, apperrors.NewDatabaseError("failed to get comment", err)
	}

	return comment, nil
}

// ListCommentsByTaskID retrieves all comments for a task with pagination
func (r *commentRepository) ListCommentsByTaskID(ctx context.Context, taskID, limit, offset int) ([]domain.Comment, int, error) {
	countQuery := `SELECT COUNT(*) FROM comments WHERE task_id = $1`
	var total int
	err := r.db.QueryRow(ctx, countQuery, taskID).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("failed to count comments", err)
	}

	const query = `
		SELECT id, task_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE task_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, taskID, limit, offset)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("failed to list comments", err)
	}
	defer rows.Close()

	comments := make([]domain.Comment, 0)
	for rows.Next() {
		var c domain.Comment
		err := rows.Scan(
			&c.ID,
			&c.TaskID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, 0, apperrors.NewDatabaseError("failed to scan comment", err)
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.NewDatabaseError("error iterating comments", err)
	}

	return comments, total, nil
}

// UpdateComment updates a comment
func (r *commentRepository) UpdateComment(ctx context.Context, id int, content string) (*domain.Comment, error) {
	const query = `
		UPDATE comments
		SET content = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, task_id, user_id, content, created_at, updated_at
	`

	comment := &domain.Comment{}
	err := r.db.QueryRow(ctx, query, content, id).Scan(
		&comment.ID,
		&comment.TaskID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrCommentNotFound, "comment not found")
		}
		return nil, apperrors.NewDatabaseError("failed to update comment", err)
	}

	return comment, nil
}

// DeleteComment deletes a comment
func (r *commentRepository) DeleteComment(ctx context.Context, id int) error {
	const query = `DELETE FROM comments WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return apperrors.NewDatabaseError("failed to delete comment", err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.NewNotFoundError(apperrors.ErrCommentNotFound, "comment not found")
	}

	return nil
}

// ListRecentComments retrieves recent comments from all tasks
func (r *commentRepository) ListRecentComments(ctx context.Context, limit, offset int) ([]domain.Comment, int, error) {
	const query = `
		SELECT id, task_id, user_id, content, created_at, updated_at
		FROM comments
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM comments").Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	comments := make([]domain.Comment, 0)
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
