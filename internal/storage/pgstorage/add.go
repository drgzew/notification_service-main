package pgstorage

import (
	"context"
	"log"
	"time"

	"notification_service/internal/models"
)

func (s *PGStorage) UpdateNotificationStatus(ctx context.Context, statuses []*models.NotificationStatus) error {
	if len(statuses) == 0 {
		return nil
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()

	stmt := `
	INSERT INTO notification_status(notification_id, status, old_status, error, sent_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (notification_id)
	DO UPDATE SET
		old_status = notification_status.status,
		status = EXCLUDED.status,
		error = EXCLUDED.error,
		sent_at = EXCLUDED.sent_at,
		updated_at = EXCLUDED.updated_at
	`

	for _, ns := range statuses {
		var oldStatus string
		err := tx.QueryRow(ctx, `SELECT status FROM notification_status WHERE notification_id=$1`, ns.NotificationID).Scan(&oldStatus)
		if err != nil && err.Error() != "no rows in result set" {
			tx.Rollback(ctx)
			log.Printf("Failed to get old status for %s: %v", ns.NotificationID, err)
			return err
		}
		ns.OldStatus = oldStatus
		ns.UpdatedAt = time.Now()

		log.Printf("Updating status for notification ID=%s: %s -> %s", ns.NotificationID, ns.OldStatus, ns.Status)

		_, err = tx.Exec(ctx, stmt,
			ns.NotificationID,
			ns.Status,
			ns.OldStatus,
			ns.ErrorMessage,
			ns.SentAt,
			ns.UpdatedAt,
		)
		if err != nil {
			tx.Rollback(ctx)
			log.Printf("Failed to update status for %s: %v", ns.NotificationID, err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction for notification statuses: %v", err)
		return err
	}

	return nil
}

func (s *PGStorage) AddNotification(ctx context.Context, n *models.Notification) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO notification(id, recipient, message, created_at) VALUES($1, $2, $3, $4)`,
		n.ID, n.Recipient, n.Message, n.CreatedAt,
	)
	return err
}


