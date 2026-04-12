# syntax=docker/dockerfile:1

# --- Build Stage ---
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install git for go mod download if needed
RUN apk add --no-cache git

# Copy go mod, sum, and vendor folder
COPY go.mod go.sum ./
COPY vendor/ ./vendor/
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app using vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o travels_backend_api main.go

# --- Run Stage ---
FROM alpine:latest
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/travels_backend_api .

# Copy any static/config files if needed (uncomment if required)
COPY api-ngnix.conf ./

# Environment variables will be provided at runtime using --env-file .env

EXPOSE 8080

CMD ["./travels_backend_api"]
