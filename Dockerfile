#Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/weather-forecast/main.go

#Run stage
FROM alpine:3.19
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/subscription.html .



EXPOSE 8080
