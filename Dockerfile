FROM golang:1.25-alpine AS builder

WORKDIR /app

# proto-generation нужен из-за replace в go.mod
COPY proto-generation /app/proto-generation

WORKDIR /app/message-service
COPY message-service/go.mod message-service/go.sum ./
RUN go mod download

COPY message-service/ .
RUN go build -o message-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/message-service/message-service .

EXPOSE 9004
EXPOSE 50051

CMD ["./message-service"]
