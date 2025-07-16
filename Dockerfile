FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ✅ Enable CGO and install C compiler + mupdf libs
ENV CGO_ENABLED=1
RUN apt-get update && apt-get install -y gcc libc6-dev libmupdf-dev

RUN go build -o backend ./cmd/main.go

FROM debian:bullseye-slim

# ✅ Install runtime dependencies
RUN apt-get update && apt-get install -y libmupdf-dev && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/backend /usr/local/bin/backend

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/backend"]
