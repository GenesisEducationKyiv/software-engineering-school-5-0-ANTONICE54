#Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY services/email/go.mod services/email/go.sum ./services/email/
COPY pkg/go.mod pkg/go.sum ./pkg/

COPY services/email/ ./services/email/
COPY pkg/ ./pkg/

WORKDIR /app/services/email


RUN go mod tidy
RUN go mod download

RUN go build -o main ./cmd/email-service/main.go

#Run stage
FROM alpine:3.19
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/services/email/main .
COPY services/email/.env .




