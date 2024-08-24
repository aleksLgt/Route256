FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.mod

RUN go mod download

COPY . .

ARG DB_CONN_READ
ARG DB_CONN_WRITE
ARG DB_CONN_TEST
ARG GRPC_PORT
ARG HTTP_PORT
ARG JAEGER_HOST

RUN echo "DB_CONN_READ=$DB_CONN_READ" > ./.env
RUN echo "DB_CONN_WRITE=$DB_CONN_WRITE" >> ./.env
RUN echo "DB_CONN_TEST=$DB_CONN_TEST" >> ./.env
RUN echo "GRPC_PORT=$GRPC_PORT" >> ./.env
RUN echo "HTTP_PORT=$HTTP_PORT" >> ./.env
RUN echo "JAEGER_HOST=$JAEGER_HOST" >> ./.env

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/app/main.go

ENTRYPOINT ["./app"]
