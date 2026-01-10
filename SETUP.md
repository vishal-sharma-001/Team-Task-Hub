# Team Task Hub - Full Stack Setup Guide

Complete setup and deployment guide for the Team Task Hub application (Backend + Frontend).

## Quick Start

### Backend Setup (Go)

```bash
cd "d:\LV Project\team-task-hub-backend"
go run cmd/main.go
```

Backend will be available at `http://localhost:8080`

### Frontend Setup (React)

```bash
cd "d:\LV Project\team-task-hub-ui"
npm install
npm run dev
```

Frontend will be available at `http://localhost:3000`

## Project Overview

**Team Task Hub** is a full-stack task management application with two main components:

### Backend (Go)
- REST API for project and task management
- PostgreSQL database with 6 migrations
- JWT authentication
- Clean architecture (Domain, Repository, Service, Handler layers)
- 27 implemented endpoints

### Frontend (React)
- Modern single-page application
- Responsive design with Tailwind CSS
- User authentication with JWT tokens
- Project and task management UI
- Task commenting and collaboration

## Directory Structure

```
d:\LV Project\
├── team-task-hub-backend/      # Go Backend
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── domain/             # Entity models
│   │   ├── repository/         # Data access
│   │   ├── service/            # Business logic
│   │   ├── handler/            # HTTP handlers
│   │   ├── middleware/         # Authentication, logging
│   │   ├── config/             # Configuration
│   │   ├── errors/             # Custom errors
│   │   ├── app/                # App initialization
│   │   └── utils/              # Utilities
│   ├── migrations/             # Database migrations
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile
│   └── README.md
│
└── team-task-hub-ui/           # React Frontend
    ├── src/
    │   ├── api/                # API client
    │   ├── components/         # React components
    │   ├── hooks/              # Custom hooks
    │   ├── pages/              # Page components
    │   ├── App.jsx
    │   ├── main.jsx
    │   └── index.css
    ├── index.html
    ├── package.json
    ├── vite.config.js
    ├── tailwind.config.js
    ├── postcss.config.js
    ├── .gitignore
    ├── README.md
    └── SETUP.md
```

## Detailed Backend Setup

### 1. Prerequisites

- Go 1.23 or higher
- PostgreSQL 15 (Docker container recommended)
- Git

### 2. Start PostgreSQL Database

Using Docker:

```bash
docker run -d \
  --name postgres-task-hub \
  -e POSTGRES_DB=taskdb \
  -e POSTGRES_USER=taskuser \
  -e POSTGRES_PASSWORD=taskpass \
  -p 5435:5432 \
  postgres:15
```

### 3. Environment Configuration

Backend reads from `.env` file (auto-created if missing):

```
DB_HOST=localhost
DB_PORT=5435
DB_USER=taskuser
DB_PASSWORD=taskpass
DB_NAME=taskdb
JWT_SECRET=your-secret-key-change-this
SERVER_PORT=8080
```

### 4. Run Backend

```bash
cd "d:\LV Project\team-task-hub-backend"

# Run database migrations
make migrate

# Start the server
make run

# Or directly
go run cmd/main.go
```

Backend endpoints will be available at `http://localhost:8080`

### 5. Available Make Commands

```bash
make help          # Show all available commands
make run          # Run the server
make build        # Build executable
make migrate      # Run database migrations
make test         # Run tests
make clean        # Clean build files
make fmt          # Format code
```

## Detailed Frontend Setup

### 1. Prerequisites

- Node.js 16.x or higher
- npm 7.x or higher
- Backend API running on `http://localhost:8080`

### 2. Installation

```bash
cd "d:\LV Project\team-task-hub-ui"
npm install
```

This installs all dependencies:
- React 18.2.0
- React Router 6.18.0
- Axios 1.6.0
- Tailwind CSS 3.3.0
- Vite 5.0.0

### 3. Start Development Server

```bash
npm run dev
```

Opens at `http://localhost:3000` with:
- Hot module replacement (HMR)
- Automatic browser reload
- API proxy to localhost:8080

### 4. Production Build

```bash
npm run build
npm run preview
```

Creates optimized bundle in `dist/` folder.

## API Endpoints Reference

### Authentication (Public)

```
POST   /api/auth/signup           Create new account
POST   /api/auth/login            Authenticate user
GET    /api/auth/me               Get current user profile
```

### Projects (Protected)

```
GET    /api/projects              List all projects
POST   /api/projects              Create project
GET    /api/projects/:id          Get project details
PUT    /api/projects/:id          Update project
DELETE /api/projects/:id          Delete project
```

### Tasks (Protected)

```
GET    /api/projects/:id/tasks           List project tasks
POST   /api/projects/:id/tasks           Create task
GET    /api/tasks/:id                    Get task details
PUT    /api/tasks/:id                    Update task
PATCH  /api/tasks/:id/status             Update task status
DELETE /api/tasks/:id                    Delete task
POST   /api/tasks/:id/assign             Assign task to user
```

### Comments (Protected)

```
GET    /api/tasks/:id/comments           List task comments
POST   /api/tasks/:id/comments           Create comment
PUT    /api/comments/:id                 Update comment
DELETE /api/comments/:id                 Delete comment
```

## Database Schema

### Users Table
```sql
id, email, password_hash, created_at, updated_at
```

### Projects Table
```sql
id, name, description, user_id, created_at, updated_at
```

### Tasks Table
```sql
id, project_id, title, description, status, priority, 
assignee, due_date, created_at, updated_at
```

### Task_Assignments Table
```sql
id, task_id, user_id, assigned_at
```

### Comments Table
```sql
id, task_id, user_id, content, created_at, updated_at
```

### Schema Extensions Table
```sql
version, description, installed_on
```

## Authentication Flow

### Login Sequence

1. **User enters credentials** on Login page
2. **Frontend sends** POST to `/api/auth/login` with email and password
3. **Backend validates** credentials against PostgreSQL
4. **Backend returns** JWT token and user data
5. **Frontend stores** token in `localStorage` as `authToken`
6. **Subsequent requests** include token in `Authorization: Bearer <token>` header
7. **Backend validates** token in JWT middleware
8. **User redirected** to Dashboard on success

### Token Storage & Security

- Token stored in `localStorage` (browser-based)
- Token included in all protected API requests via Axios interceptor
- 401 responses trigger logout and redirect to login
- Token removed on logout

## Development Workflow

### Daily Development

1. **Start PostgreSQL**
   ```bash
   docker ps  # Check if running
   ```

2. **Start Backend**
   ```bash
   cd team-task-hub-backend
   go run cmd/main.go
   ```

3. **Start Frontend**
   ```bash
   cd team-task-hub-ui
   npm run dev
   ```

4. **Open in Browser**
   - Frontend: `http://localhost:3000`
   - Backend API Docs: Check handler files for endpoint details

### Testing Features

**Create Test Account:**
- Email: `test@example.com`
- Password: `password123`

**Test Workflows:**
1. Sign up with new account
2. Create a project (e.g., "My Project")
3. Create tasks in the project
4. Update task status
5. Add comments to tasks
6. Edit/delete tasks and projects

## Troubleshooting

### Backend Issues

**Port already in use**
```bash
# Change port in backend
# Edit: internal/config/config.go or set SERVER_PORT env var
```

**Database connection failed**
```bash
# Verify PostgreSQL is running
docker ps
# Check credentials in .env file
```

**Migrations not applied**
```bash
make migrate
# Or manually: go run cmd/main.go (auto-runs migrations)
```

### Frontend Issues

**Port 3000 already in use**
```bash
# Change port in vite.config.js
export default {
  server: {
    port: 3001  // Change to different port
  }
}
```

**API requests failing**
- Ensure backend is running on port 8080
- Check Network tab in DevTools
- Clear localStorage and reload

**Build errors**
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
npm run build
```

### Common Errors

| Error | Solution |
|-------|----------|
| Cannot connect to backend | Verify backend is running on 8080 |
| 401 Unauthorized | Clear localStorage, login again |
| CORS errors | Check proxy settings in vite.config.js |
| Password too short | Minimum 8 characters required |
| Email already exists | Use different email for signup |

## Performance Optimization

### Backend
- Database connection pooling (pgxpool)
- JWT token validation
- Efficient queries with proper indexing

### Frontend
- Code splitting with React Router
- Lazy loading components
- Optimized builds with Vite
- Tailwind CSS purging in production

## Security Best Practices

1. **Environment Variables**
   - Never commit `.env` to git
   - Use `.gitignore` to exclude secrets

2. **JWT Tokens**
   - Stored securely in localStorage
   - Sent only to trusted API endpoints
   - Validated on each request

3. **Password Security**
   - Minimum 8 characters
   - Bcrypt hashing on backend
   - Never logged or exposed

4. **API Security**
   - All non-auth endpoints require JWT
   - CORS configured for localhost only
   - Input validation on all requests

## Deployment

### Frontend Deployment (Vercel/Netlify)

1. Build locally
   ```bash
   npm run build
   ```

2. Deploy `dist` folder to:
   - Vercel: Connect GitHub repo
   - Netlify: Drag and drop `dist` folder

3. Set API URL environment variable:
   ```
   VITE_API_URL=https://api.example.com
   ```

### Backend Deployment (Railway/Render)

1. Set environment variables:
   ```
   DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
   JWT_SECRET, SERVER_PORT
   ```

2. Run migrations on deployment
3. Connect PostgreSQL database
4. Deploy Go binary

## Monitoring & Logging

### Backend Logging

Logs include:
- Request/response times
- Database queries
- Authentication events
- Error details

### Frontend Monitoring

Check browser console for:
- API call errors
- Form validation messages
- Authentication status
- Network issues

## Next Steps

1. **Development**
   - Add more task features
   - Implement real-time updates
   - Add user notifications
   - Create admin dashboard

2. **Deployment**
   - Set up CI/CD pipeline
   - Configure SSL certificates
   - Set up monitoring and alerts
   - Enable auto-scaling

3. **Enhancements**
   - Add email notifications
   - Implement file attachments
   - Add task templates
   - Create team collaboration features

## Support & Documentation

- **Backend**: See `team-task-hub-backend/README.md`
- **Frontend**: See `team-task-hub-ui/README.md`
- **API Documentation**: Check handler implementations
- **Database**: See migrations in `internal/repository/migration.go`

## License

Team Task Hub - Full Stack Project

---

**Setup Last Updated**: 2024
**Go Version**: 1.23
**Node Version**: 16+
**PostgreSQL Version**: 15
