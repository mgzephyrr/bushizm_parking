package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"subscription/internal/models"
	"time"
)

func (s *PostgresStorage) GetManInQueueByID(ctx context.Context, userID int) (*models.ManInQueue, error) {
	rows, err := s.GetByID(ctx, "parking_queue", userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			slog.Error("Error while trying to close rows", slog.String("error", err.Error()))
		}
	}()

	for rows.Next() {
		return scanIntoManInQueue(rows)
	}

	return nil, &ErrorDoesNotExist{
		What:  fmt.Sprintf("man in queue with ID = %d", userID),
		Inner: nil,
	}
}

func (s *PostgresStorage) CreateManInQueue(ctx context.Context, id int) error {
	query := `INSERT INTO parking_queue
	(user_id)
	VALUES
	($1)
	RETURNING id`

	_, err := s.Create(ctx, query, "parking_queue", id)
	if err == nil {
		return nil
	}

	var notFoundErr *ErrorDoesNotExist
	if errors.As(err, &notFoundErr) {
		return notFoundErr
	}

	return fmt.Errorf("unexpected error: %w", err)
}

func (s *PostgresStorage) DeleteEarliest(ctx context.Context) (*models.ManInQueue, error) {
	query := `WITH cte AS (
		SELECT id FROM parking_queue
		ORDER BY created_at
		LIMIT 1
	)
	DELETE FROM parking_queue
	WHERE id IN (SELECT id FROM cte)
	RETURNING id, user_id, created_at;`

	var (
		id        int
		userID    int
		createdAt time.Time
	)

	err := s.tx.GetQueryEngine(ctx).QueryRow(query).Scan(&id, &userID, &createdAt)
	if err != nil {
		slog.Error("Error while deleting rows", slog.String("table", "parking_queue"), slog.String("error", err.Error()))
		return nil, fmt.Errorf("unexpected error while deleting earliest man in queue")
	}

	if id == 0 {
		return nil, &ErrorDoesNotExist{
			What:  "earliest man in queue",
			Inner: nil,
		}
	}

	return &models.ManInQueue{ID: id, UserID: userID, CreatedAt: createdAt}, nil
}

func (s *PostgresStorage) GetAllQueue(ctx context.Context) ([]int, error) {
	rows, err := s.GetAll(ctx, "parking_queue")
	if err != nil {
		return nil, err
	}

	miqs, err := getManInQueueSlice(rows)
	if err != nil {
		return nil, err
	}

	result := make([]int, len(miqs))
	for i, miq := range miqs {
		result[i] = miq.UserID
	}

	return result, nil
}

func (s *PostgresStorage) AddSubToEnd(ctx context.Context, id int) error {
	// проверить что еще нет в очереди
	_, err := s.GetManInQueueByID(ctx, id)
	if err != nil {
		return err
	}

	// добавить в очередь
	err = s.CreateManInQueue(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) MoveToNotificationQueue(ctx context.Context, now time.Time) error {
	rows, err := s.GetAll(ctx, "parking_queue")
	if err != nil {
		return err
	}

	miqs, err := getManInQueueSlice(rows)
	if err != nil {
		return err
	}

	if len(miqs) == 0 {
		return nil
	}

	miq, err := s.DeleteEarliest(ctx)
	if err != nil {
		return err
	}

	err = s.AddSubToNotificationQueue(ctx, miq)
	if err != nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastDequeueTimes = append(s.lastDequeueTimes, now)
	if len(s.lastDequeueTimes) > maxDequeueTimes {
		s.lastDequeueTimes = s.lastDequeueTimes[1:] // храним последние 20
	}

	return nil
}
