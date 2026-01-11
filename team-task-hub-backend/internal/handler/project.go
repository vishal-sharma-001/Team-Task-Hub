package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/launchventures/team-task-hub-backend/internal/service"
	"github.com/launchventures/team-task-hub-backend/internal/utils"
)

type projectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *projectHandler {
	return &projectHandler{projectService: projectService}
}

// CreateProject handles POST /api/projects
func (h *projectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	project, err := h.projectService.CreateProject(ctx, userID, req.Name, req.Description)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewSuccessResponse(project, "Project created successfully"))
}

// ListProjects handles GET /api/projects
func (h *projectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	// Parse pagination parameters
	page := 1
	pageSize := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	ctx := context.Background()
	projects, total, err := h.projectService.ListProjects(ctx, userID, page, pageSize)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewPaginatedResponse(projects, total, page, pageSize, "Projects retrieved successfully"))
}

// GetProject handles GET /api/projects/{project_id}
func (h *projectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	projectID := chi.URLParam(r, "project_id")

	ctx := context.Background()
	project, err := h.projectService.GetProject(ctx, projectID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(project, "Project retrieved successfully"))
}

// UpdateProject handles PUT /api/projects/{project_id}
func (h *projectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	projectID := chi.URLParam(r, "project_id")

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	project, err := h.projectService.UpdateProject(ctx, projectID, req.Name, req.Description)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(project, "Project updated successfully"))
}

// DeleteProject handles DELETE /api/projects/{project_id}
func (h *projectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	projectID := chi.URLParam(r, "project_id")

	ctx := context.Background()
	err = h.projectService.DeleteProject(ctx, projectID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(nil, "Project deleted successfully"))
}
