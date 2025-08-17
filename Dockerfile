# Stage 1: Build with CGO enabled
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    libmupdf-dev \
    mupdf-tools \
    libharfbuzz-dev \
    libjpeg-dev \
    libopenjp2-7-dev

RUN go build -o backend ./cmd/main.go

# Stage 2: Runtime — use SAME image to avoid glibc issues
FROM golang:1.24 AS backend

RUN apt-get update && apt-get install -y \
    libmupdf-dev \
    libharfbuzz0b \
    libjpeg62-turbo \
    libopenjp2-7 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/backend /usr/local/bin/backend

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/backend"]

## Stage 3: Migrate db image
# Migrations image
FROM alpine:3.20 AS migrate
RUN apk add --no-cache bash curl

# install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz \
    | tar xvz -C /usr/local/bin

# copy migrations
WORKDIR /migrations
COPY migrations .

# entrypoint does nothing, Job will override CMD
ENTRYPOINT [ "migrate" ]


