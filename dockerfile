FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

# basically cd to /app
# and run all commands from there
WORKDIR /app 

COPY go.mod go.sum ./


# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o agni cmd/server/main.go

# Final stage: create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/agni .

# Create data directory for SQLite
RUN  mkdir -p /app/data

EXPOSE 8080

CMD ["./agni"]

