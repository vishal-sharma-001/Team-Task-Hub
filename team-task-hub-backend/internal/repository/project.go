package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/launchventures/team-task-hub-backend/internal/domain"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
)

// ProjectRepository defines project data access operations
type ProjectRepository interface {
	CreateProject(ctx context.Context, userID, createdByID string, name, description string) (*domain.Project, error)
	GetProjectByID(ctx context.Context, id string) (*domain.Project, error)
	ListProjectsByUserID(ctx context.Context, userID string, limit, offset int) ([]domain.Project, int, error)
	UpdateProject(ctx context.Context, id string, name, description string) (*domain.Project, error)
	DeleteProject(ctx context.Context, id string) error
}

type projectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) ProjectRepository {
	return &projectRepository{db: db}
}

// CreateProject creates a new project
func (r *projectRepository) CreateProject(ctx context.Context, userID, createdByID string, name, description string) (*domain.Project, error) {
	projectID := uuid.New().String()
	const query = `
			WITH inserted AS (
				INSERT INTO projects (id, user_id, name, description, created_by_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
				RETURNING id, user_id, name, description, created_by_id, created_at, updated_at
			)
			SELECT i.id, i.user_id, i.name, i.description, i.created_by_id, i.created_at, i.updated_at,
			       cb.id, cb.email, cb.name
			FROM inserted i
			LEFT JOIN users cb ON i.created_by_id = cb.id
		`

		project := &domain.Project{}
		var creatorID *string
		var creatorEmail *string
		var creatorName *string
		err := r.db.QueryRow(ctx, query, projectID, userID, name, description, createdByID).Scan(
			&project.ID,
			&project.UserID,
			&project.Name,
			&project.Description,
			&creatorID,
			&project.CreatedAt,
			&project.UpdatedAt,
			&creatorID,
			&creatorEmail,
			&creatorName,
		)

		if err != nil {
			return nil, apperrors.NewDatabaseError("failed to create project", err)
		}

		if creatorID != nil {
			project.CreatedByID = creatorID
			email := ""
			name := ""
			if creatorEmail != nil {
				email = *creatorEmail
			}
			if creatorName != nil {
				name = *creatorName
			}
			project.CreatedBy = &domain.User{
				ID:    *creatorID,
				Email: email,
				Name:  name,
			}
		}

		return project, nil
}

// GetProjectByID retrieves a project by ID
func (r *projectRepository) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	const query = `
		SELECT p.id, p.user_id, p.name, p.description, p.created_by_id, p.created_at, p.updated_at,
		       cb.id, cb.email, cb.name
		FROM projects p
		LEFT JOIN users cb ON p.created_by_id = cb.id
		WHERE p.id = $1
	`

	project := &domain.Project{}
	var createdByID *string
	var createdByEmail *string
	var createdByName *string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&createdByID,
		&project.CreatedAt,
		&project.UpdatedAt,
		&createdByID,
		&createdByEmail,
		&createdByName,
	)

	if createdByID != nil {
		project.CreatedByID = createdByID
		project.CreatedBy = &domain.User{
			ID:    *createdByID,
			Email: *createdByEmail,
			Name:  *createdByName,
		}
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFoundError(apperrors.ErrProjectNotFound, "project not found")
		}
		return nil, apperrors.NewDatabaseError("failed to get project", err)
	}

	return project, nil
}

// ListProjectsByUserID retrieves all projects with pagination (shared across all users)
func (r *projectRepository) ListProjectsByUserID(ctx context.Context, userID string, limit, offset int) ([]domain.Project, int, error) {
	countQuery := `SELECT COUNT(*) FROM projects`
	var total int
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("failed to count projects", err)
	}

	const query = `
		SELECT p.id, p.user_id, p.name, p.description, p.created_by_id, p.created_at, p.updated_at,
		       cb.id, cb.email, cb.name
		FROM projects p
		LEFT JOIN users cb ON p.created_by_id = cb.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, apperrors.NewDatabaseError("failed to list projects", err)
	}
	defer rows.Close()

	projects := make([]domain.Project, 0)
	for rows.Next() {
		var p domain.Project
		var createdByID *string
		var createdByEmail *string
		var createdByName *string
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.Description,
			&createdByID,
			&p.CreatedAt,
			&p.UpdatedAt,
			&createdByID,
			&createdByEmail,
			&createdByName,
		)
		if err != nil {
			return nil, 0, apperrors.NewDatabaseError("failed to scan project", err)
		}

		if createdByID != nil {
			p.CreatedByID = createdByID
			p.CreatedBy = &domain.User{
				ID:    *createdByID,
				Email: *createdByEmail,
				Name:  *createdByName,
			}
		}

		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.NewDatabaseError("error iterating projects", err)
	}

	return projects, total, nil
}

// UpdateProject updates a project
func (r *projectRepository) UpdateProject(ctx context.Context, id string, name, description string) (*domain.Project, error) {
	const query = `
			WITH updated AS (
				UPDATE projects
				SET name = $1, description = $2, updated_at = NOW()
				WHERE id = $3
				RETURNING id, user_id, name, description, created_by_id, created_at, updated_at
			)
			SELECT u.id, u.user_id, u.name, u.description, u.created_by_id, u.created_at, u.updated_at,
			       cb.id, cb.email, cb.name
			FROM updated u
			LEFT JOIN users cb ON u.created_by_id = cb.id
		`

		project := &domain.Project{}
		var creatorID *string
		var creatorEmail *string
		var creatorName *string
		err := r.db.QueryRow(ctx, query, name, description, id).Scan(
			&project.ID,
			&project.UserID,
			&project.Name,
			&project.Description,
			&creatorID,
			&project.CreatedAt,
			&project.UpdatedAt,
			&creatorID,
			&creatorEmail,
			&creatorName,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, apperrors.NewNotFoundError(apperrors.ErrProjectNotFound, "project not found")
			}
			return nil, apperrors.NewDatabaseError("failed to update project", err)
		}

		if creatorID != nil {
			project.CreatedByID = creatorID
			email := ""
			name := ""
			if creatorEmail != nil {
				email = *creatorEmail
			}
			if creatorName != nil {
				name = *creatorName
			}
			project.CreatedBy = &domain.User{
				ID:    *creatorID,
				Email: email,
				Name:  name,
			}
		}

		return project, nil
}

// DeleteProject deletes a project (cascades to tasks)
func (r *projectRepository) DeleteProject(ctx context.Context, id string) error {
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
