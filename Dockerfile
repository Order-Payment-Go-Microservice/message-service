FROM golang:1.25-alpine AS builder

RUN apk add --no-cache protobuf protobuf-dev

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

WORKDIR /app

# Generate proto stubs for message and notification packages
COPY proto-generation/go.mod /app/proto-generation/go.mod
COPY proto-generation/gen/message /app/proto-generation/gen/message
COPY proto-generation/gen/notification /app/proto-generation/gen/notification

RUN cd /app/proto-generation && \
    protoc -I . -I /usr/include \
        --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        gen/message/v1/message.proto && \
    protoc -I . -I /usr/include \
        --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        gen/notification/v1/notification.proto

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
