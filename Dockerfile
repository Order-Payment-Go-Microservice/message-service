FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN go build -o message-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/message-service .
COPY --from=builder /app/.env .

EXPOSE 9004
EXPOSE 50051
EXPOSE 9005

CMD ["./message-service"]
