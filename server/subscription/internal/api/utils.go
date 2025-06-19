package api

import (
	"context"
	"errors"
	"subscription/internal/models"
	"time"
)

var (
	ErrQueueFull      = errors.New("queue is full")
	ErrAlreadyInQueue = errors.New("user already in queue")
)

type Queue interface {
	GetAllQueue(context.Context) ([]int, error)
	AddSubToEnd(context.Context, int) error
	MoveToNotificationQueue(context.Context, time.Time) error
	GetUserPosition(context.Context, int) (int, error)

	NotifiedQueuePeekBack(context.Context) (models.Subscription, bool)
	NotifiedQueuePopBack(context.Context) (models.Subscription, bool)
	EstimateWaitTime(int) time.Duration
}

type ParkingService interface {
	CheckAvailableSpots() (int, error)
}

type NotificationService interface {
	Notify(models.Subscription) error
}
