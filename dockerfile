# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install required build tools and SQLite dev
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy only go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies before copying full source
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application statically
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o agni ./cmd/server

# Stage 2: Minimal final image
FROM alpine:latest

# Install only SQLite runtime if needed (optional)
RUN apk add --no-cache sqlite-libs

WORKDIR /app

# Copy built binary
COPY --from=builder /app/agni .

# Create writable directory for SQLite data (optional)
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./agni"]
