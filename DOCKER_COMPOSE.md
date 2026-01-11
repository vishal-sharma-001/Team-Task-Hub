# Team Task Hub - Docker Compose Setup

This guide explains how to run the entire application stack using Docker Compose.

## Prerequisites

- Docker (20.10+)
- Docker Compose (2.0+)

## Quick Start

1. **Create `docker-compose.yml` in project root:**

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

2. **Create `Dockerfile` for backend** (`team-task-hub-backend/Dockerfile`):

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

3. **Create `Dockerfile` for frontend** (`team-task-hub-ui/Dockerfile`):

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

# Rebuild images
docker-compose up -d --build
```

## Accessing the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/api
- **PostgreSQL**: localhost:5432

## Environment Variables

### Backend (.env or docker-compose.yml)

- `DB_HOST`: Database host (postgres)
- `DB_PORT`: Database port (5432)
- `DB_USER`: Database user (postgres)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (task_hub)
- `JWT_SECRET`: Secret key for JWT signing (change in production)
- `PORT`: API port (8080)

### Frontend

- `VITE_API_URL`: Backend API URL (http://localhost:8080/api)

## Production Considerations

For production deployment:

1. **Security:**
   - Change default database password
   - Use strong JWT secret
   - Enable HTTPS/SSL
   - Set appropriate CORS headers

2. **Database:**
   - Use managed PostgreSQL service (AWS RDS, Cloud SQL, etc.)
   - Enable automated backups
   - Configure proper resource limits

3. **Environment:**
   - Use `.env` files (not committed to repo)
   - Use secrets management (Docker Secrets, Kubernetes Secrets)
   - Set NODE_ENV=production
   - Enable logging and monitoring

4. **Scaling:**
   - Use load balancer for multiple backend instances
   - Use CDN for frontend assets
   - Consider container orchestration (Kubernetes)

## Example Production docker-compose.yml

```yaml
version: '3.8'

services:
  backend:
    image: your-registry/task-hub-backend:latest
    environment:
      DB_HOST: ${DB_HOST}
      DB_PASSWORD: ${DB_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
    restart: always
    networks:
      - task-hub-network
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M

  frontend:
    image: your-registry/task-hub-frontend:latest
    environment:
      VITE_API_URL: ${API_URL}
    restart: always
    networks:
      - task-hub-network
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M

networks:
  task-hub-network:
```

## Troubleshooting

**Database connection fails:**
```bash
# Check database is healthy
docker-compose ps postgres

# Check logs
docker-compose logs postgres
```

**Backend can't connect to database:**
```bash
# Ensure postgres service is started first
docker-compose up -d postgres
docker-compose up -d backend
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
