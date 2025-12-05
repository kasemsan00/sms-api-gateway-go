# -----------------------------------------------------------------------------
# Multi-stage build for production optimization
# -----------------------------------------------------------------------------

# Stage 1: Build the Go application
FROM golang:1.25.4-alpine AS builder

# Set working directory
WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary, GOOS=linux for Linux binary
# Build from cmd/api directory (new project structure)
RUN CGO_ENABLED=0 GOOS=linux go build -o zteco-api-go ./cmd/api

# Stage 2: Create the production image
FROM alpine:latest

# Install ca-certificates for HTTPS requests and tzdata for timezone support
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user for security
RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/zteco-api-go /app/zteco-api-go

# Copy migrations directory (optional: for manual migration runs)
# COPY --from=builder /app/migrations /app/migrations

# Create logs directory
RUN mkdir -p /app/logs && \
    chown -R appuser:appgroup /app/logs

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8090

# Set environment variables for production
ENV ENVIRONMENT=production
ENV PORT=8090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8090/health || exit 1

# Run the application
CMD ["./zteco-api-go"]
