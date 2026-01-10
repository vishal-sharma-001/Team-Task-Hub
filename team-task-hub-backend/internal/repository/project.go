package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
	"github.com/launchventures/team-task-hub-backend/internal/domain"
)

// ProjectRepository defines project data access operations
type ProjectRepository interface {
	CreateProject(ctx context.Context, userID int, name, description string) (*domain.Project, error)
	GetProjectByID(ctx context.Context, id int) (*domain.Project, error)
	ListProjectsByUserID(ctx context.Context, userID, limit, offset int) ([]domain.Project, int, error)
	UpdateProject(ctx context.Context, id int, name, description string) (*domain.Project, error)
	DeleteProject(ctx context.Context, id int) error
}

type projectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) ProjectRepository {
	return &projectRepository{db: db}
}

// CreateProject creates a new project
func (r *projectRepository) CreateProject(ctx context.Context, userID int, name, description string) (*domain.Project, error) {
	const query = `
		INSERT INTO projects (user_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, user_id, name, description, created_at, updated_at
	`

	project := &domain.Project{}
	err := r.db.QueryRow(ctx, query, userID, name, description).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, apperrors.NewDatabaseError("failed to create project", err)
	}

	return project, nil
}

// GetProjectByID retrieves a project by ID
func (r *projectRepository) GetProjectByID(ctx context.Context, id int) (*domain.Project, error) {
	const query = `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	project := &domain.Project{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrProjectNotFound, "project not found")
		}
		return nil, apperrors.NewDatabaseError("failed to get project", err)
	}

	return project, nil
}

// ListProjectsByUserID retrieves all projects for a user with pagination
func (r *projectRepository) ListProjectsByUserID(ctx context.Context, userID, limit, offset int) ([]domain.Project, int, error) {
	countQuery := `SELECT COUNT(*) FROM projects WHERE user_id = $1`
	var total int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("failed to count projects", err)
	}

	const query = `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("failed to list projects", err)
	}
	defer rows.Close()

	projects := make([]domain.Project, 0)
	for rows.Next() {
		var p domain.Project
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, apperrors.NewDatabaseError("failed to scan project", err)
		}
		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.NewDatabaseError("error iterating projects", err)
	}

	return projects, total, nil
}

// UpdateProject updates a project
func (r *projectRepository) UpdateProject(ctx context.Context, id int, name, description string) (*domain.Project, error) {
	const query = `
		UPDATE projects
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, user_id, name, description, created_at, updated_at
	`

	project := &domain.Project{}
	err := r.db.QueryRow(ctx, query, name, description, id).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrProjectNotFound, "project not found")
		}
		return nil, apperrors.NewDatabaseError("failed to update project", err)
	}

	return project, nil
}

// DeleteProject deletes a project (cascades to tasks)
func (r *projectRepository) DeleteProject(ctx context.Context, id int) error {
	const query = `DELETE FROM projects WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return apperrors.NewDatabaseError("failed to delete project", err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.NewNotFoundError(apperrors.ErrProjectNotFound, "project not found")
	}

	return nil
}
