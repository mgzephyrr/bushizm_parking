package inmem

import (
	"subscription/internal/api/models"
	"time"

	"github.com/idsulik/go-collections/deque"
)

type InMemStorage struct {
	queue *deque.Deque[models.Subscription]
}

func NewInMemStorage() *InMemStorage {
	return &InMemStorage{
		queue: deque.New[models.Subscription](10),
	}
}

func (s *InMemStorage) AddSubToEnd(id int, now time.Time) (models.Subscription, error) {
	sub := models.NewSubscription(id, now)
	s.queue.PushFront(sub)

	return sub, nil
}
