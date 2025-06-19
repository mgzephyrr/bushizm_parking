package api

import (
	"errors"
	"subscription/internal/models"
	"time"
)

var (
	ErrQueueFull      = errors.New("queue is full")
	ErrAlreadyInQueue = errors.New("user already in queue")
)

type Queue interface {
	GetAllQueue() []int
	AddSubToEnd(int) error
	MoveToNotificationQueue(time.Time) error
	GetUserPosition(int) (int, error)

	NotifiedQueuePeekBack() (models.Subscription, bool)
	NotifiedQueuePopBack() (models.Subscription, bool)
	EstimateWaitTime(int) time.Duration
}

type ParkingService interface {
	CheckAvailableSpots() (int, error)
}

type NotificationService interface {

}
