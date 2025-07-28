# Use official Go image as build environment
FROM golang:alpine AS builder

# Build arguments
ARG BUILDTIME
ARG VERSION
ARG REVISION

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with version info
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags "-s -w -X main.buildTime=${BUILDTIME} -X main.version=${VERSION} -X main.revision=${REVISION}" \
    -o main cmd/server/main.go

# Use minimal alpine image for final stage
FROM alpine:latest

# Install security updates and create non-root user
RUN apk --no-cache add ca-certificates tzdata wget && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create app directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]