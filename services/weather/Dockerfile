#Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY services/weather/go.mod services/weather/go.sum ./services/weather/
COPY pkg/go.mod pkg/go.sum ./pkg/

COPY services/weather/ ./services/weather/
COPY pkg/ ./pkg/

WORKDIR /app/services/weather


RUN go mod tidy
RUN go mod download

RUN go build -o main ./cmd/weather-service/main.go

#Run stage
FROM alpine:3.19
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/services/weather/main .
COPY services/weather/.env .




EXPOSE 8081
