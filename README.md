# Team Task Hub - Full Stack Application

A modern full-stack task management application built with Go and React.

## Quick Start

### Prerequisites
- Go 1.23+
- Node.js 18+
- PostgreSQL 15+ (or Docker)

### Backend Setup
```bash
cd team-task-hub-backend
go run ./cmd/team-task-hub/main.go
```
Server runs on `http://localhost:8080`

### Frontend Setup
```bash
cd team-task-hub-ui
npm install
npm run dev
```
App runs on `http://localhost:3000`

## Features

- **Authentication**: Secure JWT-based auth with password hashing
- **Projects**: Create and manage projects
- **Tasks**: Create, assign, and track tasks with status and priority
- **Comments**: Collaborate on tasks with comments
- **Dashboard**: View assigned tasks, grouped by project with progress tracking
- **Responsive UI**: Works on desktop, tablet, and mobile

## Documentation

See `SETUP.md` for complete setup, deployment, and architecture details.

## Project Structure

```
team-task-hub-backend/   # Go API server
├── cmd/                 # Application entry point
├── internal/            # Core business logic
│   ├── handler/         # HTTP request handlers
│   ├── service/         # Business logic
│   ├── repository/      # Data access layer
│   ├── domain/          # Data models
│   └── middleware/      # Auth & logging
└── migrations/          # Database migrations

team-task-hub-ui/        # React frontend
├── src/
│   ├── api/            # API client
│   ├── components/     # Reusable components
│   ├── pages/          # Page components
│   ├── hooks/          # Custom hooks
│   └── index.css       # Tailwind CSS
└── public/             # Static assets
```

## API Endpoints

### Authentication
- `POST /api/auth/signup` - Create account
- `POST /api/auth/login` - Login
- `GET /api/auth/me` - Current user

### Projects
- `GET /api/projects` - List projects
- `POST /api/projects` - Create project
- `PUT /api/projects/{id}` - Update project
- `DELETE /api/projects/{id}` - Delete project

### Tasks
- `GET /api/projects/{id}/tasks` - List project tasks
- `GET /api/tasks/assigned` - Get your assigned tasks
- `POST /api/projects/{id}/tasks` - Create task
- `PUT /api/tasks/{id}` - Update task
- `PATCH /api/tasks/{id}/status` - Update task status
- `POST /api/tasks/{id}/assign` - Assign task to user
- `DELETE /api/tasks/{id}` - Delete task

### Comments
- `GET /api/tasks/{id}/comments` - List task comments
- `GET /api/comments/recent` - Recent comments feed
- `POST /api/tasks/{id}/comments` - Add comment
- `PUT /api/comments/{id}` - Update comment
- `DELETE /api/comments/{id}` - Delete comment

## Tech Stack

**Backend:**
- Go 1.23
- PostgreSQL 15
- Chi router
- JWT authentication

**Frontend:**
- React 18
- React Router 6
- Axios
- Tailwind CSS
- Vite

## Development

### Running Tests
```bash
cd team-task-hub-backend
go test ./...
```

### Building for Production
```bash
# Backend
cd team-task-hub-backend
go build -o dist/app ./cmd/team-task-hub

# Frontend
cd team-task-hub-ui
npm run build
```

## Database

The application uses PostgreSQL with automated migrations. Database schema includes:
- users
- projects
- tasks
- comments
- task_assignments

## Security Features

- JWT-based authentication
- Password hashing with bcrypt
- CORS middleware
- SQL injection prevention via parameterized queries
- Protected routes requiring authentication
- User isolation (can only see own data)

## License

MIT
