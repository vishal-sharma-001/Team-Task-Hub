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
