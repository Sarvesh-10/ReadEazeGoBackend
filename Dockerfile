# Stage 1: Build the Go binary statically
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO disabled for a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend ./cmd

# Stage 2: Minimal runtime image
FROM alpine:latest

# Install libc6-compat only if you need it (for dynamic binaries)
# For static binary, you can skip this line
# RUN apk add --no-cache libc6-compat

COPY --from=builder /app/backend /usr/local/bin/backend

EXPOSE 8080

CMD ["/usr/local/bin/backend"]
