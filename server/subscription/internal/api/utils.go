package api

import (
	"subscription/internal/api/models"
	"time"

	"github.com/idsulik/go-collections/deque"
)

type Queue interface {
	GetAllQueue() *deque.Deque[int]
	AddSubToEnd(int) error
	MoveToNotificationQueue(time.Time) error
	RemoveFromNotificationQueue(time.Time) error

	NotifiedQueuePeekBack() (models.Subscription, bool)
	NotifiedQueuePopBack() (models.Subscription, bool)
}
