FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.mod

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/app/main.go

ENTRYPOINT ["./app"]
