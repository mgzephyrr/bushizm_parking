package models

import (
	"fmt"
	"log/slog"
	"time"
)

const (
	EXPIRATION_TIME   = 15 * time.Minute
	maxNotifyAttempts = 3
)

type Subscription struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	NotifyAttempts int       `json:"notify_attempts"`
}

func NewSubscription(id int, userID int, now time.Time) Subscription {
	return Subscription{
		ID:        id,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(EXPIRATION_TIME),
	}
}

func (sub *Subscription) Notify() error {
	if sub.NotifyAttempts >= maxNotifyAttempts {
		return fmt.Errorf("max notification attempts reached")
	}

	sub.NotifyAttempts++

	slog.Info("Notifying user with", slog.Int("ID", sub.UserID))

	return nil
}
