# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install curl for healthcheck
RUN apk --no-cache add curl

# Copy the binary from builder
COPY --from=builder /app/server /app/
COPY --from=builder /app/.env* /app/

# Create logs directory
RUN mkdir -p /app/logs

# Expose the port
EXPOSE 3000

# Add metadata labels
LABEL org.opencontainers.image.title="Supabase Self-Hosted MCP Server"
LABEL org.opencontainers.image.description="Go implementation of Supabase Self-Hosted MCP Server"
LABEL org.opencontainers.image.version="1.0.0"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1

# Run the application
CMD ["/app/server"]
