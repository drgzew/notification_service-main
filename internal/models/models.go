package models

import "time"

type Notification struct {
	ID        string     `json:"id"`
	Recipient string     `json:"recipient"`
	Message   string     `json:"message"`
	CreatedAt *time.Time `json:"created_at"`
}

type NotificationStatus struct {
	NotificationID string     `db:"notification_id"` 
	OldStatus      string     `db:"old_status"`      
	Status         string     `db:"status"`          
	ErrorMessage   string     `db:"error"` 
	CreatedAt      *time.Time `json:"created_at"`          
	SentAt         *time.Time `db:"sent_at"`      
	UpdatedAt      time.Time  `db:"updated_at"`     
}