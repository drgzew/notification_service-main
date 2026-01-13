package pgstorage

import "time"

type NotificationStatus struct {
	NotificationID string     `db:"notification_id"` 
	Status         string     `db:"status"`          
	ErrorMessage   string     `db:"error"`           
	SentAt         *time.Time `db:"sent_at"`        
}

const (
	tableName        = "notification_status"
	IDColumnName     = "notification_id"
	StatusColumnName = "status"
	ErrorColumnName  = "error"
	SentAtColumnName = "sent_at"
) 
