# Stage 1: Build
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN go build -o backend .

# Stage 2: Run
FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/backend /usr/local/bin/backend

EXPOSE 8080

CMD ["backend"]
