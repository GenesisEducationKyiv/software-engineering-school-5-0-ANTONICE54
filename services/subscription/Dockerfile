#Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY services/subscription/go.mod services/subscription/go.sum ./services/subscription/
COPY pkg/go.mod pkg/go.sum ./pkg/

COPY services/subscription/ ./services/subscription/
COPY pkg/ ./pkg/

WORKDIR /app/services/subscription


RUN go mod tidy
RUN go mod download

RUN go build -o main ./cmd/subscription-service/main.go

#Run stage
FROM alpine:3.19
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/services/subscription/main .
COPY services/subscription/.env .




EXPOSE 8082
