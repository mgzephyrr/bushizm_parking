package api

import (
	"errors"
	"subscription/internal/api/models"
	"time"

	"github.com/idsulik/go-collections/deque"
)

var (
	ErrQueueFull      = errors.New("queue is full")
	ErrAlreadyInQueue = errors.New("user already in queue")
)

type Queue interface {
	GetAllQueue() *deque.Deque[int]
	AddSubToEnd(int) error
	MoveToNotificationQueue(time.Time) error

	NotifiedQueuePeekBack() (models.Subscription, bool)
	NotifiedQueuePopBack() (models.Subscription, bool)
}
