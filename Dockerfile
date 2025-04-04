FROM golang:1.24.0

WORKDIR /app
COPY . .

RUN go mod download

CMD ["go", "run", "cmd/main.go"]
