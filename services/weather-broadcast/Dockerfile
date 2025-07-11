#Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY services/weather-broadcast/go.mod services/weather-broadcast/go.sum ./services/weather-broadcast/
COPY pkg/go.mod pkg/go.sum ./pkg/

COPY services/weather-broadcast/ ./services/weather-broadcast/
COPY pkg/ ./pkg/

WORKDIR /app/services/weather-broadcast


RUN go mod tidy
RUN go mod download

RUN go build -o main ./cmd/weather-broadcast-service/main.go

#Run stage
FROM alpine:3.19
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/services/weather-broadcast/main .
COPY services/weather-broadcast/.env .




