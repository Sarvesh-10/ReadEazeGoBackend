# Stage 1: Build
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o backend ./cmd

# Stage 2: Run
FROM alpine:latest

# Install libc (some Go binaries need it)
RUN apk add --no-cache libc6-compat

COPY --from=builder /app/backend /usr/local/bin/backend

EXPOSE 8080

CMD ["/usr/local/bin/backend"]
