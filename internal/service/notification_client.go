package service

type NotificationClient interface {
	SendPushNotification(recipientID, content string) error
}
