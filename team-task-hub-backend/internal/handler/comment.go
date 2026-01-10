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

type commentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *commentHandler {
	return &commentHandler{commentService: commentService}
}

// CreateComment handles POST /api/projects/{project_id}/tasks/{task_id}/comments
func (h *commentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	taskID, err := strconv.Atoi(chi.URLParam(r, "task_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	comment, err := h.commentService.CreateComment(ctx, taskID, userID, req.Content)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewSuccessResponse(comment, "Comment created successfully"))
}

// ListComments handles GET /api/projects/{project_id}/tasks/{task_id}/comments
func (h *commentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	taskID, err := strconv.Atoi(chi.URLParam(r, "task_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
	comments, total, err := h.commentService.ListComments(ctx, taskID, page, pageSize)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewPaginatedResponse(comments, total, page, pageSize, "Comments retrieved successfully"))
}

// UpdateComment handles PUT /api/projects/{project_id}/tasks/{task_id}/comments/{comment_id}
func (h *commentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	commentID, err := strconv.Atoi(chi.URLParam(r, "comment_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	var req UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	// TODO: Add ownership verification (userID must match comment.user_id)
	comment, err := h.commentService.UpdateComment(ctx, commentID, req.Content)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(comment, "Comment updated successfully"))
}

// DeleteComment handles DELETE /api/projects/{project_id}/tasks/{task_id}/comments/{comment_id}
func (h *commentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	commentID, err := strconv.Atoi(chi.URLParam(r, "comment_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	// TODO: Add ownership verification (userID must match comment.user_id)
	err = h.commentService.DeleteComment(ctx, commentID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(nil, "Comment deleted successfully"))
}

// ListRecentComments handles GET /api/comments/recent
func (h *commentHandler) ListRecentComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	// Parse pagination parameters
	page := 1
	pageSize := 10

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
	comments, total, err := h.commentService.ListRecentComments(ctx, page, pageSize)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewPaginatedResponse(comments, total, page, pageSize, "Recent comments retrieved successfully"))
}
