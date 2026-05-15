package grpc

import (
	"context"
	"log"
	"time"

	"github.com/Order-Payment-Go-Microservice/message-service/internal/service"
	notificationv1 "github.com/Order-Payment-Go-Microservice/proto-generation/gen/notification/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type notificationClient struct {
	addr string
}

func NewNotificationClient(addr string) service.NotificationClient {
	return &notificationClient{addr: addr}
}

func (c *notificationClient) SendPushNotification(recipientID, content string) error {
	conn, err := grpc.NewClient(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("[gRPC Client] Failed to connect: %v", err)
		return err
	}
	defer conn.Close()

	client := notificationv1.NewNotificationServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.SendNotification(ctx, &notificationv1.NotificationRequest{
		UserId:  recipientID,
		Title:   "Новое сообщение",
		Message: content,
		Type:    "push",
	})

	if err != nil {
		log.Printf("[gRPC Client] Error sending notification: %v", err)
		return err
	}

	log.Printf("[gRPC Client] Notification sent successfully to %s", recipientID)
	return nil
}
