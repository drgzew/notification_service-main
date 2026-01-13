package pgstorage

import (
	"context"

	"notification_service/internal/models"
)

func (s *PGStorage) UpdateNotificationStatus(ctx context.Context, statuses []*models.NotificationStatus) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt := `
	INSERT INTO notification_status(notification_id, status, error, sent_at)
	VALUES($1, $2, $3, $4)
	ON CONFLICT(notification_id)
	DO UPDATE SET status=EXCLUDED.status, error=EXCLUDED.error, sent_at=EXCLUDED.sent_at
	`

	for _, ns := range statuses {
		_, err := tx.Exec(ctx, stmt,
			ns.NotificationID,
			ns.Status,
			ns.ErrorMessage,
			ns.SentAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

