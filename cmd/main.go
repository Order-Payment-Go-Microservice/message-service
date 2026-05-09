package main

import (
	"database/sql"
	"fmt"
	"log"
	"message-service/internal/config"
	"message-service/internal/grpc"
	"message-service/internal/handler"
	"message-service/internal/repository"
	"message-service/internal/service"
	"time"


	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	dbAddr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	
	var db *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", dbAddr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Waiting for database... attempt %d/5", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after retries: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id UUID PRIMARY KEY,
		chat_id UUID NOT NULL,
		sender_id UUID NOT NULL,
		receiver_id UUID NOT NULL,
		content TEXT NOT NULL,
		message_type VARCHAR(20) DEFAULT 'text',
		is_read BOOLEAN DEFAULT FALSE,
		is_delivered BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	repo := repository.NewPostgresMessageRepository(db)
	hub := service.NewHub()
	notificationClient := grpc.NewNotificationClient(cfg.NotificationServiceAddr)
	messageService := service.NewMessageService(repo, notificationClient, hub)

	messageHandler := handler.NewMessageHandler(messageService)
	wsHandler := handler.NewWebSocketHandler(hub)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.POST("/messages", messageHandler.SendMessage)
	router.GET("/messages/:chatId", messageHandler.GetHistory)
	router.POST("/messages/read/:id", messageHandler.MarkRead)
	router.GET("/ws", wsHandler.HandleConnections)

	router.GET("/messages/message/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "GET single message placeholder"}) })
	router.PUT("/messages/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "PUT message placeholder"}) })
	router.DELETE("/messages/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "DELETE message placeholder"}) })
	router.POST("/messages/delivered/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "POST delivered placeholder"}) })
	router.GET("/messages/search", func(c *gin.Context) { c.JSON(200, gin.H{"message": "GET search placeholder"}) })
	router.POST("/messages/reply", func(c *gin.Context) { c.JSON(200, gin.H{"message": "POST reply placeholder"}) })
	router.POST("/messages/forward", func(c *gin.Context) { c.JSON(200, gin.H{"message": "POST forward placeholder"}) })
	router.POST("/messages/react", func(c *gin.Context) { c.JSON(200, gin.H{"message": "POST react placeholder"}) })
	router.GET("/messages/pinned", func(c *gin.Context) { c.JSON(200, gin.H{"message": "GET pinned placeholder"}) })

	log.Printf("Message Service starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
