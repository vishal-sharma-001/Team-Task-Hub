# Team Task Hub

A full-stack task management application built with Go backend and React frontend, enabling teams to collaborate on projects and tasks.

## Features

- **Project Management**: Create and manage projects visible to all team members
- **Task Management**: Create, assign, and track tasks within projects
- **User Assignment**: Assign tasks to team members with visibility of assignees
- **Task Status & Priority**: Track task progress with status (OPEN, IN_PROGRESS, DONE) and priority levels (LOW, MEDIUM, HIGH)
- **User Authentication**: Secure JWT-based authentication
- **Team Collaboration**: All users can access all projects and tasks created by team members

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Chi (HTTP router)
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Migrations**: golang-migrate

### Frontend
- **Framework**: React 18+
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **HTTP Client**: Axios

## Project Structure

```
team-task-hub/
├── team-task-hub-backend/          # Go backend
│   ├── cmd/team-task-hub/          # Entry point
│   ├── internal/
│   │   ├── handler/                # HTTP handlers
│   │   ├── service/                # Business logic
│   │   ├── repository/             # Data access layer
│   │   ├── domain/                 # Domain models
│   │   ├── middleware/             # HTTP middleware
│   │   ├── config/                 # Configuration
│   │   ├── errors/                 # Error handling
│   │   └── utils/                  # Utilities
│   ├── migrations/                 # Database migrations
│   ├── go.mod & go.sum             # Go dependencies
│   └── .env                        # Environment configuration
│
└── team-task-hub-ui/               # React frontend
    ├── src/
    │   ├── components/             # Reusable React components
    │   ├── pages/                  # Page components
    │   ├── api/                    # API client
    │   ├── hooks/                  # Custom React hooks
    │   ├── App.jsx                 # Main app component
    │   └── index.css               # Global styles
    ├── index.html
    ├── vite.config.js
    ├── tailwind.config.js
    └── package.json
```

## How to Run

### Prerequisites
- Go 1.21 or higher
- Node.js 16+ and npm
- PostgreSQL 12+

### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd team-task-hub-backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   Create/update `.env` file:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=task_hub
   JWT_SECRET=your_jwt_secret_key
   PORT=8080
   ```

4. **Run database migrations**
   ```bash
   go run ./cmd/team-task-hub/
   ```
   The server will automatically run migrations on startup.

5. **Start the server**
   ```bash
   go run ./cmd/team-task-hub/
   ```
   Server will be available at `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd team-task-hub-ui
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Start development server**
   ```bash
   npm run dev
   ```
   Application will be available at `http://localhost:3001` (or next available port)

4. **Build for production**
   ```bash
   npm run build
   ```

## Architecture Overview

### Backend Architecture

**Layered Architecture Pattern:**

```
┌─────────────────────────────────┐
│      HTTP Handlers (REST API)   │
│     (handler/task.go, etc.)     │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│     Service Layer (Business)    │
│    (service/task.go, etc.)      │
│  - Validation & Business Logic  │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│   Repository Layer (Data Access)│
│   (repository/task.go, etc.)    │
│  - Database queries & mapping   │
└──────────────┬──────────────────┘
               │
┌──────────────▼──────────────────┐
│      PostgreSQL Database        │
│  (users, projects, tasks, etc.) │
└─────────────────────────────────┘
```

**Key Components:**
- **Domain Models**: Data structures representing core entities (User, Project, Task)
- **Handlers**: HTTP request/response processing
- **Services**: Business logic and validation
- **Repositories**: Database operations
- **Middleware**: Authentication (JWT), CORS, request logging

### Frontend Architecture

**Component-Based Architecture:**

```
App.jsx (Router)
├── Pages/
│   ├── Login/Signup
│   ├── Projects (List all projects)
│   ├── TaskBoard (Project tasks)
│   └── TaskDetail (Single task)
└── Components/
    ├── TaskForm (Create/Edit tasks)
    ├── ProjectForm (Create/Edit projects)
    ├── Modal
    └── ErrorMessage
```

**State Management:**
- React Hooks (useState, useEffect, useContext)
- Custom hooks for async operations (useAsync)
- API client (Axios) for backend communication

## API Documentation

### Authentication Endpoints

**POST /api/auth/signup**
- Register a new user
- Body: `{ email, password }`
- Response: `{ user, token }`

**POST /api/auth/login**
- Login user
- Body: `{ email, password }`
- Response: `{ user, token }`

### Project Endpoints

**GET /api/projects**
- List all projects (paginated)
- Query: `page`, `page_size`
- Response: `{ data: [projects], total, page, pages }`

**POST /api/projects**
- Create new project
- Body: `{ name, description }`
- Response: `{ data: project }`

**GET /api/projects/{id}**
- Get project details
- Response: `{ data: project }`

**PUT /api/projects/{id}**
- Update project
- Body: `{ name, description }`
- Response: `{ data: project }`

**DELETE /api/projects/{id}**
- Delete project

### Task Endpoints

**GET /api/projects/{projectId}/tasks**
- List project tasks
- Query: `status`, `priority`, `page`, `page_size`
- Response: `{ data: [tasks], total, page, pages }`

**POST /api/projects/{projectId}/tasks**
- Create task
- Body: `{ title, description, priority, assignee_id, due_date }`
- Response: `{ data: task }`

**GET /api/tasks/{id}**
- Get task details
- Response: `{ data: task }`

**PUT /api/tasks/{id}**
- Update task
- Body: `{ title, description, status, priority, assignee_id, due_date }`
- Response: `{ data: task }`

**PATCH /api/tasks/{id}/status**
- Update task status
- Body: `{ status }`
- Response: `{ data: task }`

**PATCH /api/tasks/{id}/priority**
- Update task priority
- Body: `{ priority }`
- Response: `{ data: task }`

**PATCH /api/tasks/{id}/assignee**
- Update task assignee
- Body: `{ assignee_id }`
- Response: `{ data: task }`

**DELETE /api/tasks/{id}**
- Delete task

### User Endpoints

**GET /api/users**
- List all users
- Response: `[users]`

**GET /api/auth/me**
- Get current user profile
- Response: `{ data: user }`

## Database Design

### Schema Overview

**Users Table**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Projects Table**
```sql
CREATE TABLE projects (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Tasks Table**
```sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'OPEN',
    priority VARCHAR(50) DEFAULT 'MEDIUM',
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_by_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Comments Table**
```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Key Relationships

- **Users → Projects**: One user can create many projects (user_id in projects)
- **Users → Tasks**: Users can be assigned to multiple tasks (assignee_id in tasks)
- **Projects → Tasks**: One project contains many tasks (project_id in tasks)
- **Tasks → Comments**: One task can have multiple comments (task_id in comments)
- **Users → Comments**: One user can write multiple comments (user_id in comments)

### Data Integrity
- Foreign keys ensure referential integrity
- CASCADE delete for related records
- SET NULL for optional references (e.g., deleted assignees)
- Timestamps automatically managed (created_at, updated_at)

## Authentication Flow

1. User registers/logs in
2. Backend validates credentials and returns JWT token
3. Frontend stores token in localStorage
4. Subsequent requests include token in Authorization header: `Bearer {token}`
5. Backend middleware validates token and extracts user ID
6. User ID is available in request context for authorization

## Team Collaboration Features

- **Shared Projects**: All projects are visible to all authenticated users
- **Task Assignment**: Tasks can be assigned to any team member
- **Creator Tracking**: Original creator of projects/tasks is recorded
- **Assignee Tracking**: Current assignee is visible in task details
- **Audit Trail**: created_by_id and assigned_by_id fields track who performed actions
