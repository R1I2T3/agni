# Stage 1: Build
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev curl git

# Install air
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
ENV PATH="/root/.air:$PATH"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Optional: run `air` from here (development)
# ENTRYPOINT ["air"]

# Or still build binary for prod
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o agni ./cmd/server

# Stage 2: Minimal final image
FROM alpine:latest
RUN apk add --no-cache sqlite-libs

WORKDIR /app
COPY --from=builder /app/agni .
# Only needed if you're copying air too
# COPY --from=builder /root/.air/air /usr/local/bin/

RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./agni"]
