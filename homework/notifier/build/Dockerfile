# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY ./ ./

RUN pwd && ls -alg && CGO_ENABLED=0 GOOS=linux go build ./cmd/notifier

# Deploy the application binary into a lean image
FROM alpine:3 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/notifier /notifier

ENTRYPOINT ["/notifier", "-bootstrap-server=kafka0:29092"]
