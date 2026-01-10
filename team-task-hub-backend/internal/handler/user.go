package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/launchventures/team-task-hub-backend/internal/service"
	"github.com/launchventures/team-task-hub-backend/internal/utils"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{userService: userService}
}

// SignUp handles POST /api/auth/signup
func (h *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	user, token, err := h.userService.SignUp(ctx, req.Email, req.Password)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	authResp := AuthResponse{
		User:  user,
		Token: token,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewSuccessResponse(authResp, "User registered successfully"))
}

// Login handles POST /api/auth/login
func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	user, token, err := h.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	authResp := AuthResponse{
		User:  user,
		Token: token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(authResp, "Login successful"))
}

// GetProfile handles GET /api/auth/me
func (h *userHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from context (set by auth middleware)
	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	user, err := h.userService.GetProfile(ctx, userID)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(user, "Profile retrieved successfully"))
}

// UpdateProfile handles PUT /api/auth/me
func (h *userHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from context (set by auth middleware)
	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	user, err := h.userService.UpdateProfile(ctx, userID, req.Name)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(user, "Profile updated successfully"))
}

// ListUsers handles GET /api/users
func (h *userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verify user is authenticated
	_, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	ctx := context.Background()
	users, err := h.userService.ListUsers(ctx)
	if err != nil {
		w.WriteHeader(ErrorToStatusCode(err))
		json.NewEncoder(w).Encode(NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewSuccessResponse(users, "Users retrieved successfully"))
}
