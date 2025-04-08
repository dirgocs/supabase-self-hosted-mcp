# Build stage
FROM golang:1.21-alpine AS build

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/server .

# Copy necessary files
COPY --from=build /app/.env.example ./.env

# Expose the port
EXPOSE 3000

# Run the application
CMD ["./server"]
