# Stage 1: Build the Go binary statically with CGO disabled
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o backend ./cmd

# Stage 2: Use Debian slim base image (includes glibc)
FROM debian:bullseye-slim

COPY --from=builder /app/backend /usr/local/bin/backend

EXPOSE 8080

CMD ["/usr/local/bin/backend"]
