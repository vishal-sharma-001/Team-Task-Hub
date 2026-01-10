package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/launchventures/team-task-hub-backend/internal/utils"
)

// AuthMiddleware validates JWT tokens and adds user ID to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"status":"error","error":"Unauthorized","message":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Extract bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"status":"error","error":"Unauthorized","message":"invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		claims, appErr := utils.ValidateToken(token)
		if appErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appErr.StatusCode())
			errResp := map[string]interface{}{
				"status":  "error",
				"error":   appErr.Code,
				"message": appErr.Message,
			}
			json.NewEncoder(w).Encode(errResp)
			return
		}

		// Add user ID and email to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ErrorMiddleware handles and formats API errors consistently, with panic recovery
func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("ERROR: panic recovered in %s %s: %v", r.Method, r.URL.Path, err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				errResp := map[string]interface{}{
					"status":  "error",
					"error":   "InternalServerError",
					"message": "an unexpected error occurred",
				}
				json.NewEncoder(w).Encode(errResp)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		duration := time.Since(startTime).Milliseconds()
		log.Printf("[RESPONSE] %s %s - %dms", r.Method, r.URL.Path, duration)
	})
}
