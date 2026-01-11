# Team Task Hub - Docker Compose Setup

This guide explains how to run the entire application stack using Docker Compose.

## Prerequisites

- Docker (20.10+)
- Docker Compose (2.0+)

## Quick Start

Run this command from the project root:

```bash
make up
```

This builds and starts all services. Access the app at:
- **Frontend:** http://localhost:3000
- **Backend:** http://localhost:8080/api  
- **Database:** localhost:5432 (postgres/postgres_password)

---

## Docker Files Included

### 1. `docker-compose.yml` (project root)

```yaml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: task-hub-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres_password
      POSTGRES_DB: task_hub
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - task-hub-network

  # Go Backend
  backend:
    build:
      context: ./team-task-hub-backend
      dockerfile: Dockerfile
    container_name: task-hub-backend
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres_password
      DB_NAME: task_hub
      JWT_SECRET: your_jwt_secret_key_change_in_production
      PORT: 8080
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - task-hub-network
    restart: unless-stopped

  # React Frontend
  frontend:
    build:
      context: ./team-task-hub-ui
      dockerfile: Dockerfile
    container_name: task-hub-frontend
    ports:
      - "3000:80"
    environment:
      VITE_API_URL: http://localhost:8080/api
    depends_on:
      - backend
    networks:
      - task-hub-network
    restart: unless-stopped

volumes:
  postgres_data:
    driver: local

networks:
  task-hub-network:
    driver: bridge
```

### 2. Backend Dockerfile (`team-task-hub-backend/Dockerfile`)

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o team-task-hub ./cmd/team-task-hub

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/team-task-hub .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./team-task-hub"]
```

### 3. Frontend Dockerfile (`team-task-hub-ui/Dockerfile`)

```dockerfile
# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build application
RUN npm run build

# Runtime stage
FROM nginx:alpine

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built application
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

4. **Create `nginx.conf` for frontend** (`team-task-hub-ui/nginx.conf`):

```nginx
server {
    listen 80;
    location / {
        root /usr/share/nginx/html;
        index index.html index.htm;
        try_files $uri $uri/ /index.html;
    }
    location /api {
        proxy_pass http://backend:8080;
    }
}
```

## Running with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Remove volumes (data will be deleted)
docker-compose down -v

## Common Commands

```bash
make up              # Build and start all services
make down            # Stop all services  
make logs            # View all logs (follow mode)
make logs-backend    # View backend logs only
make logs-frontend   # View frontend logs only
make logs-db         # View database logs only
make ps              # Show running containers
make restart         # Restart all services
make restart-backend # Restart backend only
make restart-frontend # Restart frontend only
make clean           # Remove all containers, images, volumes
```

See [Makefile](Makefile) for all available commands.

---

## Production Considerations

### Security
- Change default database password
- Use strong JWT secret
- Enable HTTPS/SSL in Nginx
- Set appropriate CORS headers
- Use environment files for secrets (not committed to repo)

### Database
- Use managed PostgreSQL service (AWS RDS, Google Cloud SQL, Azure)
- Enable automated backups
- Configure resource limits
- Set up monitoring and alerts

### Scaling
- Use load balancer for multiple backend instances
- Use CDN for frontend assets
- Consider container orchestration (Kubernetes)
- Implement rate limiting

### Example Production docker-compose.yml

```yaml
version: '3.8'

services:
  backend:
    image: your-registry/task-hub-backend:v1.0
    environment:
      DB_HOST: ${DB_HOST}
      DB_PASSWORD: ${DB_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
    restart: always
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M

  frontend:
    image: your-registry/task-hub-frontend:v1.0
    environment:
      VITE_API_URL: ${API_URL}
    restart: always
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
```

---

## Troubleshooting

### Database connection fails

```bash
# Check if postgres container is running
docker-compose ps postgres

# Check postgres logs
make logs-db

# Restart database
docker-compose restart postgres
```

### Backend can't reach database

```bash
# Start only postgres first
docker-compose up -d postgres

# Wait a few seconds, then start backend
docker-compose up -d backend

# Check backend logs
make logs-backend
```

### Frontend can't reach backend API

```bash
# Check if backend is running
docker-compose ps backend

# Check backend logs
make logs-backend

# Verify VITE_API_URL environment variable is set correctly
docker-compose config | grep VITE_API_URL
```

### Build failures

```bash
# Clean up and rebuild
make clean
make up

# Or with more verbosity
docker-compose up --build --verbose
```

### Port already in use

If port 3000, 8080, or 5432 is already in use, edit `docker-compose.yml` and change the host port mapping:

```yaml
frontend:
  ports:
    - "3001:80"  # Changed from 3000:80
```

---

## Volume Management

The database data is persisted in the `postgres_data` volume. To remove all data:

```bash
make clean  # Removes containers, images, and volumes
```

To backup the database:

```bash
docker exec task-hub-postgres pg_dump -U postgres task_hub > backup.sql
```

To restore the database:

```bash
docker exec -i task-hub-postgres psql -U postgres task_hub < backup.sql
```

**Frontend can't reach API:**
```bash
# Check backend is running
docker-compose logs backend

# Verify API URL in frontend environment
docker-compose logs frontend
```

**Permission issues on volumes:**
```bash
# Fix volume permissions
sudo chown -R 999:999 postgres_data
```
