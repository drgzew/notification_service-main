package models

import "time"

type KafkaNotification struct {
	BatchID      string
	Notification Notification
}

type Notification struct {
	ID        string     `json:"id"`
	Recipient string     `json:"recipient"`
	Channel   string     `json:"channel"`
	Message   string     `json:"message"`
	CreatedAt *time.Time `json:"created_at"`
}

type NotificationStatus struct {
	NotificationID string     `json:"notification_id"` 
	Status         string     `json:"status"`          
	ErrorMessage   string     `json:"error"`          
	SentAt         *time.Time `json:"sent_at"`         
}
