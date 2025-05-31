# Stage 1: Build
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

# Stage 2: Run
FROM alpine:latest


COPY --from=builder /app/backend /usr/local/bin/backend

ENV PORT=8080
EXPOSE $PORT

CMD ["backend"]
