FROM golang:1.24-alpine AS builder

# Install MySQL client & development headers for CGO builds
RUN apk add --no-cache gcc musl-dev mariadb-connector-c-dev curl git

# Install air for hot reload (development)
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
ENV PATH="/root/.air:$PATH"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Create needed dirs
RUN mkdir -p tmp data
RUN chmod 755 tmp

EXPOSE 8080

ENTRYPOINT ["air", "-c", ".air.toml", "-d"]