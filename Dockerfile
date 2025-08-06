FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main .

FROM alpine:latest

RUN apk add --no-cache postgresql-client

COPY --from=builder /app/main /main

COPY wait.sh /wait.sh
RUN chmod +x /wait.sh

EXPOSE 8080