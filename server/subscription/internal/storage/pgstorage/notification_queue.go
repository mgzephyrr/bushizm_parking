package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"subscription/internal/models"
	"time"
)

func (s *PostgresStorage) AddSubToNotificationQueue(ctx context.Context, miq *models.ManInQueue) error {
	now := time.Now().Truncate(time.Second)
	sub := models.NewSubscription(0, miq.UserID, now)

	query := `INSERT INTO notification_queue
	(user_id, created_at, expires_at, notify_attempts)
	VALUES
	($1, $2, $3, $4)`

	_, err := s.Create(ctx, query, "parking_queue", sub.UserID, sub.CreatedAt, sub.ExpiresAt, sub.NotifyAttempts)
	if err == nil {
		return nil
	}

	var notFoundErr *ErrorDoesNotExist
	if errors.As(err, &notFoundErr) {
		return notFoundErr
	}

	return fmt.Errorf("unexpected error: %w", err)
}

func (s *PostgresStorage) NotifiedQueuePeekBack(ctx context.Context) (models.Subscription, bool) {
	query := `SELECT id, user_id, created_at, expires_at, notify_attempts
	FROM notification_queue
	ORDER BY created_at
	LIMIT 1`

	var (
		id             int
		userID         int
		createdAt      time.Time
		expiresAt      time.Time
		notifyAttempts int
	)

	err := s.tx.GetQueryEngine(ctx).QueryRow(query).Scan(&id, &userID, &createdAt, &expiresAt, &notifyAttempts)
	if err != nil {
		slog.Error("Error while selecting rows", slog.String("table", "parking_queue"), slog.String("error", err.Error()))
		return models.Subscription{}, false
	}

	return models.Subscription{
		ID:             id,
		UserID:         userID,
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAt,
		NotifyAttempts: notifyAttempts,
	}, true
}

func (s *PostgresStorage) NotifiedQueuePopBack(ctx context.Context) (models.Subscription, bool) {
	query := `WITH cte AS (
		SELECT id
		FROM notification_queue
		ORDER BY created_at
		LIMIT 1
	)
	DELETE FROM notification_queue
	WHERE id IN (SELECT id FROM cte)
	RETURNING id, user_id, created_at, expires_at, notify_attempts;`

	var (
		id             int
		userID         int
		createdAt      time.Time
		expiresAt      time.Time
		notifyAttempts int
	)

	err := s.tx.GetQueryEngine(ctx).QueryRow(query).Scan(&id, &userID, &createdAt, &expiresAt, &notifyAttempts)
	if err != nil {
		slog.Error("Error while selecting rows", slog.String("table", "parking_queue"), slog.String("error", err.Error()))
		return models.Subscription{}, false
	}

	return models.Subscription{
		ID:             id,
		UserID:         userID,
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAt,
		NotifyAttempts: notifyAttempts,
	}, true
}
