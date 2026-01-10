# Team Task Hub - Backend

A REST API built with Go, Chi, and PostgreSQL for task management.

## 📋 Overview

The backend provides APIs for:
- User authentication (signup, login)
- Project management (CRUD)
- Task management (CRUD, filtering, assignment)
- Comments on tasks
- Dashboard with task summary

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- PostgreSQL 15+
- Make (optional)

### Setup

1. Update `.env` with your database credentials (already configured):
```
DB_HOST=localhost
DB_PORT=5435
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=team_task_hub
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-change-in-production
SERVER_PORT=8080
```

2. Build the application:
```bash
make build
# or
go build -o bin/team-task-hub ./cmd/team-task-hub/main.go
```

3. Run the server (migrations run automatically):
```bash
make run
# or
./bin/team-task-hub
```

The server will start on `http://localhost:8080`

## 📁 Project Structure

```
team-task-hub-backend/
├── cmd/
│   └── team-task-hub/
│       └── main.go                # Application entry point
│
├── internal/                      # Private packages
│   ├── app/
│   │   └── app.go                # App initialization & router setup
│   │
│   ├── config/
│   │   └── config.go             # Configuration from env variables
│   │
│   ├── domain/                   # Domain entities (one per file)
│   │   ├── user.go
│   │   ├── project.go
│   │   ├── task.go
│   │   ├── comment.go
│   │   └── task_assignment.go
│   │
│   ├── repository/               # Data access layer (one per entity)
│   │   ├── user.go
│   │   ├── project.go
│   │   ├── task.go
│   │   └── comment.go
│   │
│   ├── service/                  # Business logic layer (one per entity)
│   │   ├── user.go
│   │   ├── project.go
│   │   ├── task.go
│   │   └── comment.go
│   │
│   ├── handler/                  # HTTP handlers (one per entity)
│   │   ├── user.go
│   │   ├── project.go
│   │   ├── task.go
│   │   ├── comment.go
│   │   └── responses.go          # Shared DTOs
│   │
│   ├── middleware/
│   │   └── auth.go               # JWT auth, error handling, logging
│   │
│   ├── errors/
│   │   └── errors.go             # Custom error types
│   │
│   └── utils/                    # Utility functions
│       ├── jwt.go                # Token generation/validation
│       ├── password.go           # Bcrypt hashing
│       └── validation.go         # Input validation
│
├── migrations/                    # Database migrations
│   ├── 000001_init_extensions.up/down.sql
│   ├── 000002_create_users_table.up/down.sql
│   ├── 000003_create_projects_table.up/down.sql
│   ├── 000004_create_tasks_table.up/down.sql
│   ├── 000005_create_task_assignments_table.up/down.sql
│   └── 000006_create_comments_table.up/down.sql
│
├── .env                           # Environment variables (configured)
├── .gitignore                     # Git ignore rules
├── go.mod                         # Go module file
├── go.sum                         # Go dependencies hash
├── Makefile                       # Build and run commands
└── README.md                      # This file
```

## 🗄️ Database Schema

### Tables
- **users** - User accounts with email and password hash
- **projects** - Projects belonging to users
- **tasks** - Tasks within projects with status and priority
- **task_assignments** - User assignments to tasks
- **comments** - Comments on tasks

### Indexes
All indexes are created in migrations for optimal performance on:
- User email lookups
- Project filtering by user
- Task filtering by project, status, and assignee
- Comment retrieval by task and date

## 🧪 Testing

Test the API using curl:

### Health Check
```bash
curl http://localhost:8080/health
```

### Sign Up
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

## 🔐 Security

- Password hashing with bcrypt
- JWT authentication
- Input validation
- SQL injection prevention via prepared statements
- CORS middleware (to be implemented)

## 📊 API Endpoints

### Public Endpoints
- `GET /health` - Health check

### Authentication (Public)
- `POST /api/auth/signup` - Register new user
- `POST /api/auth/login` - Login user

### Protected Endpoints (Require JWT)

#### User
- `GET /api/auth/me` - Get current user profile

#### Projects
- `GET /api/projects` - List user's projects
- `POST /api/projects` - Create project
- `PUT /api/projects/{project_id}` - Update project
- `DELETE /api/projects/{project_id}` - Delete project

#### Tasks
- `GET /api/projects/{project_id}/tasks` - List tasks (with status/priority filters)
- `POST /api/projects/{project_id}/tasks` - Create task
- `PUT /api/projects/{project_id}/tasks/{task_id}` - Update task
- `PATCH /api/projects/{project_id}/tasks/{task_id}/status` - Update task status
- `POST /api/projects/{project_id}/tasks/{task_id}/assign` - Assign task to user
- `DELETE /api/projects/{project_id}/tasks/{task_id}` - Delete task

#### Comments
- `GET /api/projects/{project_id}/tasks/{task_id}/comments` - List comments
- `POST /api/projects/{project_id}/tasks/{task_id}/comments` - Add comment
- `PUT /api/projects/{project_id}/tasks/{task_id}/comments/{comment_id}` - Update comment
- `DELETE /api/projects/{project_id}/tasks/{task_id}/comments/{comment_id}` - Delete comment

## 🛠️ Development

### Using Makefile
```bash
make build      # Build the application
make run        # Run the application  
make clean      # Clean build artifacts
```

### Environment Variables
See `.env` for all available configuration options.

## � Dependencies

- **chi/v5** - HTTP router
- **pgx/v5** - PostgreSQL driver with connection pooling
- **golang-migrate/migrate/v4** - Database migrations
- **golang-jwt/jwt/v5** - JWT token handling
- **crypto** - Password hashing with bcrypt
- **godotenv** - Environment variable loading

## 🚀 Deployment

For production deployment:
1. Set strong `JWT_SECRET`
2. Enable SSL/TLS
3. Use environment-based configuration
4. Run migrations before startup
5. Set up proper logging
6. Configure database backups

## 📄 License

MIT License
