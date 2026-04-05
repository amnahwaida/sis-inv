# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies and ca-certificates (needed for GOTOOLCHAIN to download)
RUN apk add --no-cache gcc musl-dev ca-certificates

# Enable automatic toolchain downloads if dependencies require a newer version
ENV GOTOOLCHAIN=auto

# Copy go.mod and go.sum (if it exists)
COPY go.mod go.sum* ./
RUN go mod download

# Copy the source code
COPY . .

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o reset-admin ./cmd/reset-admin/main.go

# Stage 2: Run the binary
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for secure connections (Cloudflare, etc)
RUN apk --no-cache add ca-certificates tzdata

# Copy binaries and migrations from builder
COPY --from=builder /app/main .
COPY --from=builder /app/reset-admin .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Create uploads directory (needed for file storage)
RUN mkdir -p uploads

EXPOSE 8080

CMD ["./main"]
