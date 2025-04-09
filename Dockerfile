# Use a specific version of golang to avoid compatibility issues
FROM golang:1.21 as builder

# Set working directory
WORKDIR /app

# Copy go.mod first to leverage Docker cache
COPY go.mod ./

# Copy the rest of the code
COPY . .

# Download dependencies and build
RUN go mod tidy
RUN go build -o server .

# Use a smaller image for the final container
FROM debian:bullseye-slim

# Install CA certificates for HTTPS connections and curl for healthchecks
RUN apt-get update && \
    apt-get install -y ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Create logs directory
RUN mkdir -p /app/logs

# Copy the binary from the builder stage
COPY --from=builder /app/server /app/server

# Copy the .env file if it exists
COPY .env* /app/

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
