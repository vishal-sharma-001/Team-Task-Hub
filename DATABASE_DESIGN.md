# Team Task Hub - Database Design

## Overview

The Team Task Hub uses PostgreSQL as its relational database. The schema is designed to support project management with user collaboration, task tracking, and audit capabilities.

## Entity Relationship Diagram (ERD)

```
┌──────────────┐
│    users     │
├──────────────┤
│ id (PK)      │──┐
│ email        │  │
│ name         │  │
│ password_hash│  │
│ created_at   │  │
│ updated_at   │  │
└──────────────┘  │
       ▲          │
       │          │
       │  ┌───────┴────────────────────────┐
       │  │                                │
       │  │                                │
┌──────┴──┴────────┐          ┌───────────┴─────────┐
│   projects       │          │      tasks          │
├──────────────────┤          ├─────────────────────┤
│ id (PK)          │          │ id (PK)             │
│ user_id (FK)     │──────────│ project_id (FK)     │
│ created_by_id(FK)│          │ assignee_id (FK)    │
│ name             │          │ assigned_by_id(FK)  │
│ description      │          │ created_by_id(FK)   │
│ created_at       │          │ title               │
│ updated_at       │          │ description         │
└──────────────────┘          │ status              │
       │                      │ priority            │
       │ 1:N                  │ due_date            │
       │                      │ created_at          │
       │                      │ updated_at          │
       │              ┌───────┴─────────────┐
       │              │                     │
       │              │                     │
       │         ┌────▼────────────┐  ┌─────▼───────────┐
       │         │    comments     │  │ task_assignments│
       │         ├─────────────────┤  ├──────────────────┤
       │         │ id (PK)         │  │ id (PK)          │
       │         │ task_id (FK)    │  │ task_id (FK)     │
       │         │ user_id (FK)    │  │ user_id (FK)     │
       │         │ content         │  │ assigned_by(FK)  │
       │         │ created_at      │  │ created_at       │
       │         │ updated_at      │  └──────────────────┘
       │         └─────────────────┘
       │
```

## Tables

### users

Stores user information for authentication and profile management.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

**Fields:**
- `id`: Unique identifier (UUID)
- `email`: Email address (unique, for login)
- `name`: User's display name
- `password_hash`: Bcrypt hash of password
- `created_at`: Account creation timestamp
- `updated_at`: Last profile update timestamp

**Indexes:**
- Primary key on `id`
- Unique index on `email` (for fast login)

---

### projects

Stores project information accessible to all team members.

```sql
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    created_by_id UUID,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_created_by ON projects(created_by_id);
CREATE INDEX idx_projects_created_at ON projects(created_at DESC);
```

**Fields:**
- `id`: Unique identifier (UUID)
- `user_id`: Project owner (references users)
- `created_by_id`: Who created this project (audit trail)
- `name`: Project name
- `description`: Project description
- `created_at`: Project creation timestamp
- `updated_at`: Last update timestamp

**Relationships:**
- `user_id` → `users.id`: One-to-Many (User can own multiple projects)
- `created_by_id` → `users.id`: One-to-Many (User can create multiple projects)

**Cascade Rules:**
- When a user is deleted, their projects are deleted (CASCADE)
- When a creator is deleted, `created_by_id` is set to NULL (SET NULL)

---

### tasks

Stores task information with assignment and tracking details.

```sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'OPEN',
    priority VARCHAR(50) DEFAULT 'MEDIUM',
    assignee_id UUID,
    assigned_by_id UUID,
    created_by_id UUID NOT NULL,
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (assigned_by_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_created_by ON tasks(created_by_id);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
```

**Fields:**
- `id`: Unique identifier (UUID)
- `project_id`: Parent project (references projects)
- `title`: Task title
- `description`: Task description
- `status`: Task status (OPEN, IN_PROGRESS, DONE)
- `priority`: Task priority (LOW, MEDIUM, HIGH)
- `assignee_id`: User task is assigned to (nullable, unassigned if NULL)
- `assigned_by_id`: Who assigned this task (audit trail)
- `created_by_id`: Who created this task (audit trail)
- `due_date`: Task deadline (nullable)
- `created_at`: Task creation timestamp
- `updated_at`: Last update timestamp

**Relationships:**
- `project_id` → `projects.id`: Many-to-One (Task belongs to one project)
- `assignee_id` → `users.id`: Many-to-One (Multiple tasks to one user)
- `assigned_by_id` → `users.id`: Many-to-One (User can assign multiple tasks)
- `created_by_id` → `users.id`: Many-to-One (User can create multiple tasks)

**Cascade Rules:**
- When a project is deleted, its tasks are deleted (CASCADE)
- When an assignee is deleted, tasks remain unassigned (SET NULL)
- When creator is deleted, tasks are deleted (CASCADE)

**Indexes:**
- `idx_tasks_project_id`: Fast filtering by project
- `idx_tasks_assignee_id`: Fast filtering by assignee
- `idx_tasks_status`: Fast filtering by status
- `idx_tasks_priority`: Fast filtering by priority
- `idx_tasks_created_by`: Fast audit trail queries
- `idx_tasks_due_date`: Fast date-based queries

---

### comments

Stores task comments for collaboration.

```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_comments_task_id ON comments(task_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
```

**Fields:**
- `id`: Unique identifier (UUID)
- `task_id`: Parent task (references tasks)
- `user_id`: Comment author (references users)
- `content`: Comment text
- `created_at`: Comment creation timestamp
- `updated_at`: Last edit timestamp

**Relationships:**
- `task_id` → `tasks.id`: Many-to-One (Task can have multiple comments)
- `user_id` → `users.id`: Many-to-One (User can write multiple comments)

**Cascade Rules:**
- When a task is deleted, its comments are deleted (CASCADE)
- When a user is deleted, their comments are deleted (CASCADE)

---

## Data Integrity & Constraints

### Primary Keys
- All tables use UUID as primary key for distributed uniqueness
- UUIDs are generated server-side using `uuid_generate_v4()` PostgreSQL extension

### Foreign Keys
- Enforce referential integrity
- Prevent orphaned records
- Cascade deletes where appropriate (task deletion with project)
- Set to NULL where soft-delete is preferred (creator deletion)

### Unique Constraints
- `users.email`: Ensures no duplicate email registrations

### Not Null Constraints
- User core fields: email, password_hash
- Project core fields: name, user_id
- Task core fields: title, project_id, created_by_id
- Comment core fields: content, task_id, user_id

## Audit Trail

The schema supports audit trails through:

1. **Created By Fields** (`created_by_id`):
   - Projects: Who created the project
   - Tasks: Who created the task
   - Comments: Implicit (user_id)

2. **Assigned By Fields** (`assigned_by_id` in tasks):
   - Records who assigned a task (useful for responsibility tracking)

3. **Timestamps** (`created_at`, `updated_at`):
   - All tables track creation and modification times
   - Automatically managed by database triggers (SET DEFAULT CURRENT_TIMESTAMP)

## Performance Considerations

### Indexes
- `projects.user_id`: Fast project list queries
- `tasks.project_id`: Fast task list queries
- `tasks.assignee_id`: Fast filtering by assignee
- `tasks.status`: Fast filtering by status for dashboards
- `tasks.priority`: Fast sorting/filtering by priority
- `tasks.due_date`: Fast queries for overdue tasks
- `comments.task_id`: Fast comment retrieval for task details

### Query Patterns

**List projects for user:**
```sql
SELECT p.* FROM projects p
WHERE p.user_id = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;
```

**List tasks for project with filters:**
```sql
SELECT t.*, u.*, ab.*, cb.* FROM tasks t
LEFT JOIN users u ON t.assignee_id = u.id
LEFT JOIN users ab ON t.assigned_by_id = ab.id
LEFT JOIN users cb ON t.created_by_id = cb.id
WHERE t.project_id = $1
  AND ($2 = '' OR t.status = $2)
  AND ($3 = '' OR t.priority = $3)
ORDER BY t.created_at DESC
LIMIT $4 OFFSET $5;
```

**Find user's assigned tasks:**
```sql
SELECT t.*, p.* FROM tasks t
JOIN projects p ON t.project_id = p.id
WHERE t.assignee_id = $1
ORDER BY t.due_date ASC, t.priority DESC;
```

## Migrations

Database schema is managed through migrations (located in `migrations/` directory):

1. `000001_create_users_table.up.sql` - Create users table
2. `000002_create_projects_table.up.sql` - Create projects table
3. `000003_create_tasks_table.up.sql` - Create tasks table
4. `000004_create_comments_table.up.sql` - Create comments table

Migrations are automatically applied on server startup using `golang-migrate`.

## Future Optimization Opportunities

1. **Materialized Views**: For frequently accessed aggregations (task count per project)
2. **Partitioning**: For large tasks/comments tables (by project_id or date range)
3. **Caching**: Redis cache layer for frequently accessed projects/users
4. **Read Replicas**: For scaling read-heavy operations (dashboards)
5. **Archive Tables**: Move old tasks/comments to archive for performance
