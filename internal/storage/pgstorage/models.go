package pgstorage

import "time"

type NotificationStatus struct {
	NotificationID string     `db:"notification_id"` // ID уведомления
	OldStatus      string     `db:"old_status"`      // предыдущий статус
	Status         string     `db:"status"`          // текущий статус
	ErrorMessage   string     `db:"error"`           // ошибка, если есть
	SentAt         *time.Time `db:"sent_at"`         // время отправки
	UpdatedAt      time.Time  `db:"updated_at"`      // время обновления статуса
}

const (
	tableName        = "notification_status"
	IDColumnName     = "notification_id"
	OldStatusColumnName = "old_status"
	StatusColumnName = "status"
	ErrorColumnName  = "error"
	SentAtColumnName = "sent_at"
	UpdatedAtColumnName = "updated_at"
) 
