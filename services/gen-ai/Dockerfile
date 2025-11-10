# Build stage
FROM golang:1.23.1-alpine AS builder

# Set working directory
WORKDIR /workspace

# Install git (needed for go mod download)
RUN apk add --no-cache git

# Copy the entire workspace (this Dockerfile should be run from project root)
COPY . .

# Set working directory to the service
WORKDIR /workspace/services/gen-ai

# Download dependencies using go workspace
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /workspace/services/gen-ai/main .

# Change ownership to non-root user
RUN chown appuser:appgroup main

# Switch to non-root user
USER appuser

# Expose port (will be set via environment variable)
EXPOSE 8081

# Run the binary
CMD ["./main"]
