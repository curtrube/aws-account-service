FROM golang:1.26-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
COPY cmd ./cmd/
COPY internal ./internal/
#COPY docs ./docs/

RUN CGO_ENABLED=0 GOOS=linux go build -o aws-account-service main.go
