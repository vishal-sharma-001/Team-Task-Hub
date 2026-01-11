# Team Task Hub - API Documentation

## Base URL
```
http://localhost:8080/api
```

## Authentication
All endpoints (except `/auth/signup` and `/auth/login`) require JWT token in the `Authorization` header:
```
Authorization: Bearer {token}
```

## Response Format

### Success Response
```json
{
  "status": "success",
  "message": "Operation successful",
  "data": { /* resource data */ }
}
```

### Error Response
```json
{
  "status": "error",
  "error": "ErrorCode",
  "message": "Human readable error message"
}
```

### Paginated Response
```json
{
  "status": "success",
  "data": [ /* items */ ],
  "total": 100,
  "page": 1,
  "pages": 10,
  "message": "Items retrieved successfully"
}
```

---

## Auth Endpoints

### POST /auth/signup
Register a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": ""
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**Status Codes:** 201 Created, 400 Bad Request, 500 Internal Server Error

---

### POST /auth/login
Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**Status Codes:** 200 OK, 401 Unauthorized, 500 Internal Server Error

---

### GET /auth/me
Get current authenticated user's profile.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2026-01-12T10:00:00Z",
    "updated_at": "2026-01-12T10:00:00Z"
  }
}
```

**Status Codes:** 200 OK, 401 Unauthorized

---

### PUT /auth/me
Update current user's profile.

**Request Body:**
```json
{
  "name": "Jane Doe"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "Jane Doe"
  }
}
```

**Status Codes:** 200 OK, 400 Bad Request, 401 Unauthorized

---

## User Endpoints

### GET /users
List all users in the system.

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user1@example.com",
    "name": "User One",
    "created_at": "2026-01-12T10:00:00Z"
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "email": "user2@example.com",
    "name": "User Two",
    "created_at": "2026-01-12T10:05:00Z"
  }
]
```

**Status Codes:** 200 OK, 401 Unauthorized

---

## Project Endpoints

### GET /projects
List all projects (paginated).

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 20, max: 100)

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Website Redesign",
      "description": "Redesign company website",
      "created_by_id": "550e8400-e29b-41d4-a716-446655440000",
      "created_by": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "email": "user@example.com",
        "name": "John Doe"
      },
      "created_at": "2026-01-12T10:00:00Z",
      "updated_at": "2026-01-12T10:00:00Z"
    }
  ],
  "total": 5,
  "page": 1,
  "pages": 1
}
```

**Status Codes:** 200 OK, 401 Unauthorized

---

### POST /projects
Create a new project.

**Request Body:**
```json
{
  "name": "Mobile App",
  "description": "Build iOS and Android app"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Mobile App",
    "description": "Build iOS and Android app",
    "created_by_id": "550e8400-e29b-41d4-a716-446655440000",
    "created_by": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe"
    },
    "created_at": "2026-01-12T10:10:00Z",
    "updated_at": "2026-01-12T10:10:00Z"
  },
  "message": "Project created successfully"
}
```

**Status Codes:** 201 Created, 400 Bad Request, 401 Unauthorized

---

### GET /projects/{id}
Get project details.

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Website Redesign",
    "description": "Redesign company website",
    "created_by_id": "550e8400-e29b-41d4-a716-446655440000",
    "created_by": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe"
    },
    "created_at": "2026-01-12T10:00:00Z",
    "updated_at": "2026-01-12T10:00:00Z"
  }
}
```

**Status Codes:** 200 OK, 404 Not Found, 401 Unauthorized

---

### PUT /projects/{id}
Update a project.

**Request Body:**
```json
{
  "name": "Website Redesign v2",
  "description": "Updated website redesign plan"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Website Redesign v2",
    "description": "Updated website redesign plan",
    "created_by_id": "550e8400-e29b-41d4-a716-446655440000",
    "created_by": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe"
    },
    "created_at": "2026-01-12T10:00:00Z",
    "updated_at": "2026-01-12T10:15:00Z"
  }
}
```

**Status Codes:** 200 OK, 400 Bad Request, 404 Not Found, 401 Unauthorized

---

### DELETE /projects/{id}
Delete a project (cascades to tasks).

**Response:**
```json
{
  "status": "success",
  "message": "Project deleted successfully"
}
```

**Status Codes:** 200 OK, 404 Not Found, 401 Unauthorized

---

## Task Endpoints

### GET /projects/{projectId}/tasks
List tasks for a project (paginated, with optional filters).

**Query Parameters:**
- `status` (optional): Filter by status (OPEN, IN_PROGRESS, DONE)
- `priority` (optional): Filter by priority (LOW, MEDIUM, HIGH)
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 20, max: 100)

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440000",
      "project_id": "660e8400-e29b-41d4-a716-446655440000",
      "title": "Design homepage",
      "description": "Create mockups for homepage",
      "status": "IN_PROGRESS",
      "priority": "HIGH",
      "assignee_id": "550e8400-e29b-41d4-a716-446655440001",
      "assignee": {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "email": "designer@example.com",
        "name": "Jane Designer"
      },
      "assigned_by_id": "550e8400-e29b-41d4-a716-446655440000",
      "assigned_by": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "email": "manager@example.com",
        "name": "John Manager"
      },
      "created_by_id": "550e8400-e29b-41d4-a716-446655440000",
      "created_by": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "email": "manager@example.com",
        "name": "John Manager"
      },
      "due_date": "2026-01-20T23:59:59Z",
      "created_at": "2026-01-12T10:00:00Z",
      "updated_at": "2026-01-12T10:05:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "pages": 1
}
```

**Status Codes:** 200 OK, 401 Unauthorized

---

### POST /projects/{projectId}/tasks
Create a new task.

**Request Body:**
```json
{
  "title": "Design homepage",
  "description": "Create mockups for homepage",
  "priority": "HIGH",
  "assignee_id": "550e8400-e29b-41d4-a716-446655440001",
  "due_date": "2026-01-20T23:59:59Z"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440001",
    "project_id": "660e8400-e29b-41d4-a716-446655440000",
    "title": "Design homepage",
    "description": "Create mockups for homepage",
    "status": "OPEN",
    "priority": "HIGH",
    "assignee_id": "550e8400-e29b-41d4-a716-446655440001",
    "assignee": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "email": "designer@example.com",
      "name": "Jane Designer"
    },
    "assigned_by_id": "550e8400-e29b-41d4-a716-446655440000",
    "assigned_by": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "manager@example.com",
      "name": "John Manager"
    },
    "created_by_id": "550e8400-e29b-41d4-a716-446655440000",
    "created_by": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "manager@example.com",
      "name": "John Manager"
    },
    "due_date": "2026-01-20T23:59:59Z",
    "created_at": "2026-01-12T10:10:00Z",
    "updated_at": "2026-01-12T10:10:00Z"
  }
}
```

**Status Codes:** 201 Created, 400 Bad Request, 401 Unauthorized

---

### GET /tasks/{id}
Get task details.

**Response:** Same as task object in POST response

**Status Codes:** 200 OK, 404 Not Found, 401 Unauthorized

---

### PUT /tasks/{id}
Update a task.

**Request Body:**
```json
{
  "title": "Design homepage and footer",
  "description": "Create mockups for homepage and footer",
  "status": "IN_PROGRESS",
  "priority": "HIGH",
  "assignee_id": "550e8400-e29b-41d4-a716-446655440001",
  "due_date": "2026-01-25T23:59:59Z"
}
```

**Response:** Updated task object

**Status Codes:** 200 OK, 400 Bad Request, 404 Not Found, 401 Unauthorized

---

### PATCH /tasks/{id}/status
Update only task status.

**Request Body:**
```json
{
  "status": "DONE"
}
```

**Response:** Updated task object

**Status Codes:** 200 OK, 400 Bad Request, 404 Not Found, 401 Unauthorized

---

### PATCH /tasks/{id}/priority
Update only task priority.

**Request Body:**
```json
{
  "priority": "LOW"
}
```

**Response:** Updated task object

**Status Codes:** 200 OK, 400 Bad Request, 404 Not Found, 401 Unauthorized

---

### PATCH /tasks/{id}/assignee
Update task assignee.

**Request Body:**
```json
{
  "assignee_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

or unassign:

```json
{
  "assignee_id": null
}
```

**Response:** Updated task object with current assignee

**Status Codes:** 200 OK, 400 Bad Request, 404 Not Found, 401 Unauthorized

---

### DELETE /tasks/{id}
Delete a task.

**Response:**
```json
{
  "status": "success",
  "message": "Task deleted successfully"
}
```

**Status Codes:** 200 OK, 404 Not Found, 401 Unauthorized

---

## Error Codes

| Error Code | Status | Description |
|-----------|--------|-------------|
| InvalidInput | 400 | Invalid request data |
| Unauthorized | 401 | Missing or invalid authentication token |
| NotFound | 404 | Resource not found |
| Conflict | 409 | Resource already exists (e.g., duplicate email) |
| InternalServerError | 500 | Server error |

---

## Validation Rules

### Project
- Name: Required, 3-100 characters
- Description: Optional, max 1000 characters

### Task
- Title: Required, 3-200 characters
- Description: Optional, max 2000 characters
- Status: One of: OPEN, IN_PROGRESS, DONE
- Priority: One of: LOW, MEDIUM, HIGH
- Due Date: Optional, ISO 8601 format

### User
- Email: Required, valid email format, unique
- Password: Required, minimum 8 characters
- Name: Optional, max 255 characters
