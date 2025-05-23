package models

import (
	"fmt"
	"log/slog"
	"time"
)

const (
	EXPIRATION_TIME = 15 * time.Minute
)

type Subscription struct {
	UserID       int       `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func NewSubscription(id int, now time.Time) Subscription {
	return Subscription{
		UserID:       id,
		CreatedAt:    now,
		ExpiresAt:    now.Add(EXPIRATION_TIME),
	}
}

func (sub *Subscription) Notify() error {
	slog.Info("User notified", slog.Int("USER_ID", sub.UserID))
	return fmt.Errorf("Mock")
}
