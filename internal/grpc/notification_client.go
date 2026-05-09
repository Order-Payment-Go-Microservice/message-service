package grpc

import (
	"context"
	"log"
	"message-service/internal/service"
	pb "message-service/proto"
	"time"

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
	conn, err := grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("[gRPC Client] Failed to connect: %v", err)
		return err
	}
	defer conn.Close()

	client := pb.NewNotificationServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.SendNotification(ctx, &pb.NotificationRequest{
		UserId:  recipientID,
		Title:   "Новое сообщение",
		Message: content,
	})

	if err != nil {
		log.Printf("[gRPC Client] Error sending notification: %v", err)
		return err
	}

	log.Printf("[gRPC Client] Notification sent successfully to %s", recipientID)
	return nil
}
