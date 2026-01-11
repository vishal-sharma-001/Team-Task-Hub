package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/launchventures/team-task-hub-backend/internal/config"
	"github.com/launchventures/team-task-hub-backend/internal/handler"
	appMiddleware "github.com/launchventures/team-task-hub-backend/internal/middleware"
	"github.com/launchventures/team-task-hub-backend/internal/repository"
	"github.com/launchventures/team-task-hub-backend/internal/service"
)

// App represents the application
type App struct {
	DB     *pgxpool.Pool
	Config *config.Config
	Router *chi.Mux
}

func New(cfg *config.Config) (*App, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	app := &App{
		DB:     pool,
		Config: cfg,
		Router: chi.NewRouter(),
	}

	app.setupRoutes()
	return app, nil
}

func (a *App) setupRoutes() {
	// Global middleware - order matters!
	a.Router.Use(appMiddleware.ErrorMiddleware)   // Error handling and panic recovery
	a.Router.Use(middleware.Logger)               // Chi's built-in logger
	a.Router.Use(appMiddleware.LoggingMiddleware) // Custom request/response logging
	a.Router.Use(middleware.Recoverer)            // Chi's built-in recoverer
	a.Router.Use(corsMiddleware)                  // CORS support

	// Health check (public endpoint)
	a.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(a.DB)
	projectRepo := repository.NewProjectRepository(a.DB)
	taskRepo := repository.NewTaskRepository(a.DB)
	commentRepo := repository.NewCommentRepository(a.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	taskService := service.NewTaskService(taskRepo)
	commentService := service.NewCommentService(commentRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)
	commentHandler := handler.NewCommentHandler(commentService)

	// Public auth routes (no authentication required)
	a.Router.Post("/api/auth/signup", userHandler.SignUp)
	a.Router.Post("/api/auth/login", userHandler.Login)

	// Protected routes (authentication required)
	a.Router.Group(func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware)

		// User routes
		r.Get("/api/auth/me", userHandler.GetProfile)
		r.Put("/api/auth/me", userHandler.UpdateProfile)
		r.Get("/api/users", userHandler.ListUsers)

		// Project routes
		r.Post("/api/projects", projectHandler.CreateProject)
		r.Get("/api/projects", projectHandler.ListProjects)
		r.Get("/api/projects/{project_id}", projectHandler.GetProject)
		r.Put("/api/projects/{project_id}", projectHandler.UpdateProject)
		r.Delete("/api/projects/{project_id}", projectHandler.DeleteProject)

		// Task routes
		r.Post("/api/projects/{project_id}/tasks", taskHandler.CreateTask)
		r.Get("/api/projects/{project_id}/tasks", taskHandler.ListTasks)
		r.Get("/api/tasks/assigned", taskHandler.ListAssignedTasks)
		r.Get("/api/tasks/{task_id}", taskHandler.GetTask)
		r.Put("/api/projects/{project_id}/tasks/{task_id}", taskHandler.UpdateTask)
		r.Put("/api/tasks/{task_id}", taskHandler.UpdateTask)
		r.Patch("/api/projects/{project_id}/tasks/{task_id}/status", taskHandler.UpdateTaskStatus)
		r.Patch("/api/tasks/{task_id}/status", taskHandler.UpdateTaskStatus)
		r.Patch("/api/tasks/{task_id}/priority", taskHandler.UpdateTaskPriority)
		r.Patch("/api/tasks/{task_id}/assignee", taskHandler.UpdateTaskAssignee)
		r.Post("/api/projects/{project_id}/tasks/{task_id}/assign", taskHandler.AssignTask)
		r.Post("/api/tasks/{task_id}/assign", taskHandler.AssignTask)
		r.Delete("/api/projects/{project_id}/tasks/{task_id}", taskHandler.DeleteTask)
		r.Delete("/api/tasks/{task_id}", taskHandler.DeleteTask)

		// Comment routes
		r.Post("/api/projects/{project_id}/tasks/{task_id}/comments", commentHandler.CreateComment)
		r.Post("/api/tasks/{task_id}/comments", commentHandler.CreateComment)
		r.Get("/api/projects/{project_id}/tasks/{task_id}/comments", commentHandler.ListComments)
		r.Get("/api/tasks/{task_id}/comments", commentHandler.ListComments)
		r.Get("/api/comments/recent", commentHandler.ListRecentComments)
		r.Put("/api/projects/{project_id}/tasks/{task_id}/comments/{comment_id}", commentHandler.UpdateComment)
		r.Put("/api/comments/{comment_id}", commentHandler.UpdateComment)
		r.Delete("/api/projects/{project_id}/tasks/{task_id}/comments/{comment_id}", commentHandler.DeleteComment)
		r.Delete("/api/comments/{comment_id}", commentHandler.DeleteComment)
	})
}

func (a *App) Close() error {
	a.DB.Close()
	return nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
