#Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/
COPY pkg/go.mod pkg/go.sum ./pkg/

COPY services/gateway/ ./services/gateway/
COPY pkg/ ./pkg/

WORKDIR /app/services/gateway


RUN go mod tidy
RUN go mod download

RUN go build -o main ./cmd/gateway/main.go

#Run stage
FROM alpine:3.19
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/services/gateway/main .
COPY /services/gateway/static ./static
COPY services/gateway/.env .




