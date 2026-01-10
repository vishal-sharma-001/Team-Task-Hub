package service

import (
	"context"

	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
	"github.com/launchventures/team-task-hub-backend/internal/domain"
	"github.com/launchventures/team-task-hub-backend/internal/repository"
	"github.com/launchventures/team-task-hub-backend/internal/utils"
)

// ProjectService defines project-related business logic operations
type ProjectService interface {
	CreateProject(ctx context.Context, userID int, name, description string) (*domain.Project, error)
	GetProject(ctx context.Context, id int) (*domain.Project, error)
	ListProjects(ctx context.Context, userID, page, pageSize int) ([]domain.Project, int, error)
	UpdateProject(ctx context.Context, id int, name, description string) (*domain.Project, error)
	DeleteProject(ctx context.Context, id int) error
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{projectRepo: projectRepo}
}

// CreateProject creates a new project with validation
func (s *projectService) CreateProject(ctx context.Context, userID int, name, description string) (*domain.Project, error) {
	// Validate project name
	if appErr := utils.ValidateProjectName(name); appErr != nil {
		return nil, appErr
	}

	// Validate description (max 1000 chars)
	if len(description) > 1000 {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "description must not exceed 1000 characters")
	}

	// Create project in database
	project, err := s.projectRepo.CreateProject(ctx, userID, name, description)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// GetProject retrieves a project by ID
func (s *projectService) GetProject(ctx context.Context, id int) (*domain.Project, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid project ID")
	}

	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// ListProjects retrieves all projects for a user with pagination
func (s *projectService) ListProjects(ctx context.Context, userID, page, pageSize int) ([]domain.Project, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	projects, total, err := s.projectRepo.ListProjectsByUserID(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// UpdateProject updates a project with validation
func (s *projectService) UpdateProject(ctx context.Context, id int, name, description string) (*domain.Project, error) {
	// Validate project name
	if appErr := utils.ValidateProjectName(name); appErr != nil {
		return nil, appErr
	}

	// Validate description
	if len(description) > 1000 {
		return nil, apperrors.NewValidationError(apperrors.ErrInvalidInput, "description must not exceed 1000 characters")
	}

	// Update project in database
	project, err := s.projectRepo.UpdateProject(ctx, id, name, description)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteProject deletes a project
func (s *projectService) DeleteProject(ctx context.Context, id int) error {
	if id <= 0 {
		return apperrors.NewValidationError(apperrors.ErrInvalidInput, "invalid project ID")
	}

	err := s.projectRepo.DeleteProject(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
