package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	apperrors "github.com/launchventures/team-task-hub-backend/internal/errors"
	"github.com/launchventures/team-task-hub-backend/internal/service"
	"github.com/launchventures/team-task-hub-backend/internal/utils"
)

type taskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *taskHandler {
	return &taskHandler{taskService: taskService}
}

// CreateTask handles POST /api/projects/{project_id}/tasks
func (h *taskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	projectID := chi.URLParam(r, "project_id")

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	task, err := h.taskService.CreateTask(ctx, projectID, userID, req.Title, req.Description, req.Priority, req.AssigneeID, req.DueDate)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	// If assignee is provided, assign the task
	if req.AssigneeID != nil && *req.AssigneeID != "" {
		_ = h.taskService.AssignTask(ctx, task.ID, *req.AssigneeID, userID)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewSuccessResponse(task, "Task created successfully"))
}

// ListTasks handles GET /api/projects/{project_id}/tasks
func (h *taskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	projectID := chi.URLParam(r, "project_id")

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

	// Parse optional filters
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	ctx := context.Background()
	tasks, total, err := h.taskService.ListTasks(ctx, projectID, page, pageSize, status, priority)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewPaginatedResponse(tasks, total, page, pageSize, "Tasks retrieved successfully"))
}

// ListAssignedTasks handles GET /api/tasks/assigned
func (h *taskHandler) ListAssignedTasks(w http.ResponseWriter, r *http.Request) {
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

	// Parse optional filters
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	ctx := context.Background()
	tasks, total, err := h.taskService.ListAssignedTasks(ctx, userID, page, pageSize, status, priority)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewPaginatedResponse(tasks, total, page, pageSize, "Assigned tasks retrieved successfully"))
}

// GetTask handles GET /api/tasks/{task_id}
func (h *taskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	taskID := chi.URLParam(r, "task_id")

	ctx := context.Background()
	task, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(task, "Task retrieved successfully"))
}

// UpdateTask handles PUT /api/projects/{project_id}/tasks/{task_id}
func (h *taskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	taskID := chi.URLParam(r, "task_id")

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()

	// Extract values from pointers (empty string if nil for partial updates)
	title := ""
	if req.Title != nil {
		title = *req.Title
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	status := ""
	if req.Status != nil {
		status = *req.Status
	}

	priority := ""
	if req.Priority != nil {
		priority = *req.Priority
	}

	task, err := h.taskService.UpdateTask(ctx, taskID, title, description, status, priority, req.AssigneeID, req.DueDate)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	// If assignee was set, create assignment record
	if req.AssigneeID != nil && *req.AssigneeID != "" {
		_ = h.taskService.AssignTask(ctx, task.ID, *req.AssigneeID, userID)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(task, "Task updated successfully"))
}

// UpdateTaskStatus handles PATCH /api/projects/{project_id}/tasks/{task_id}/status
func (h *taskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	taskID := chi.URLParam(r, "task_id")

	var req UpdateTaskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	// Get current task to preserve other fields
	task, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	// Update only the status, preserve assignee if set
	updatedTask, err := h.taskService.UpdateTask(ctx, taskID, task.Title, task.Description, req.Status, task.Priority, task.AssigneeID, nil)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(updatedTask, "Task status updated successfully"))
}

// UpdateTaskPriority handles PATCH /api/tasks/{task_id}/priority
func (h *taskHandler) UpdateTaskPriority(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	taskID := chi.URLParam(r, "task_id")

	var req UpdateTaskPriorityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	// Get current task first to preserve other fields
	task, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	// Update only the priority, preserve assignee if set
	updatedTask, err := h.taskService.UpdateTask(ctx, taskID, task.Title, task.Description, task.Status, req.Priority, task.AssigneeID, nil)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(updatedTask, "Task priority updated successfully"))
}

func (h *taskHandler) UpdateTaskAssignee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	taskID := chi.URLParam(r, "task_id")

	var req UpdateTaskAssigneeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(apperrors.NewValidationError(apperrors.ErrInvalidInput, "assignee_id must be a string UUID")))
		return
	}

	ctx := context.Background()
	// If no assignee provided, clear assignment; otherwise set assignee and who assigned
	if req.AssigneeID == nil || *req.AssigneeID == "" {
		if err := h.taskService.UnassignTask(ctx, taskID); err != nil {
			w.WriteHeader(ErrorToStatusCode(err))
			json.NewEncoder(w).Encode(NewErrorResponse(err))
			return
		}
	} else {
		if err := h.taskService.AssignTask(ctx, taskID, *req.AssigneeID, userID); err != nil {
			w.WriteHeader(ErrorToStatusCode(err))
			json.NewEncoder(w).Encode(NewErrorResponse(err))
			return
		}
	}

	// Return fresh task with populated users
	updatedTask, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(updatedTask, "Task assignee updated successfully"))
}

func (h *taskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	taskID := chi.URLParam(r, "task_id")

	var req AssignTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	err = h.taskService.AssignTask(ctx, taskID, req.UserID, userID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(nil, "Task assigned successfully"))
}

// DeleteTask handles DELETE /api/projects/{project_id}/tasks/{task_id}
func (h *taskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	taskID := chi.URLParam(r, "task_id")

	ctx := context.Background()
	err = h.taskService.DeleteTask(ctx, taskID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(nil, "Task deleted successfully"))
}
