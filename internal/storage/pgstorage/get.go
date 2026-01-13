package pgstorage

import (
	"context"

	"notification_service/internal/models"
)


func (s *PGStorage) GetNotificationStatusByIDs(ctx context.Context, IDs []string) ([]*models.NotificationStatus, error) {
	query := `
	SELECT notification_id, status, error, sent_at
	FROM notification_status
	WHERE notification_id = ANY($1)
	`
	rows, err := s.db.Query(ctx, query, IDs) 
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.NotificationStatus
	for rows.Next() {
		var ns models.NotificationStatus
		if err := rows.Scan(&ns.NotificationID, &ns.Status, &ns.ErrorMessage, &ns.SentAt); err != nil {
			return nil, err
		}
		results = append(results, &ns)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}