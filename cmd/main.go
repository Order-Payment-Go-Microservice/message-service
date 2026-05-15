package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/Order-Payment-Go-Microservice/message-service/internal/config"
	"github.com/Order-Payment-Go-Microservice/message-service/internal/database"
	internalGrpc "github.com/Order-Payment-Go-Microservice/message-service/internal/grpc"
	"github.com/Order-Payment-Go-Microservice/message-service/internal/handler"
	"github.com/Order-Payment-Go-Microservice/message-service/internal/repository"
	"github.com/Order-Payment-Go-Microservice/message-service/internal/service"
	messagev1 "github.com/Order-Payment-Go-Microservice/proto-generation/gen/message/v1"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	dbAddr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	var nc *nats.Conn
	nc, err = nats.Connect(cfg.NatsURL)
	if err != nil {
		log.Printf("NATS connection failed: %v", err)
	} else {
		defer nc.Close()
	}

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("Redis connection failed: %v", err)
	}

	repo := repository.NewPostgresMessageRepository(db)
	hub := service.NewHub()
	notificationClient := internalGrpc.NewNotificationClient(cfg.NotificationServiceAddr)
	messageService := service.NewMessageService(repo, notificationClient, hub, nc, rdb)

	messageHandler := handler.NewMessageHandler(messageService)
	wsHandler := handler.NewWebSocketHandler(hub)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		messagev1.RegisterMessageServiceServer(s, internalGrpc.NewMessageServer(messageService))
		log.Println("gRPC Message Server starting on port 50051...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	router := gin.Default()
	router.POST("/messages", messageHandler.SendMessage)
	router.GET("/messages/:chatId", messageHandler.GetHistory)
	router.PATCH("/messages/:id/read", messageHandler.MarkRead)
	router.GET("/ws", wsHandler.HandleConnections)

	log.Printf("Message Service starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
