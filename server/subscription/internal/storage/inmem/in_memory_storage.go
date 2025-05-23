package inmem

import (
	"context"
	"subscription/internal/api/models"
	"sync"
	"time"

	"github.com/idsulik/go-collections/deque"
)

const WORKERS_COUNT = 1

type InMemStorage struct {
	wantToPark    *deque.Deque[int]
	notifiedQueue *deque.Deque[models.Subscription]
	mu            sync.Mutex
}

func NewInMemStorage(ctx context.Context) *InMemStorage {
	storage := &InMemStorage{
		wantToPark:    deque.New[int](10),
		notifiedQueue: deque.New[models.Subscription](10),
	}

	for i := range WORKERS_COUNT {
		go NewQueueWorker(i + 1, storage).Process(ctx)
	}

	return storage
}

func (s *InMemStorage) GetAllQueue() *deque.Deque[int] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.wantToPark
}

func (s *InMemStorage) AddSubToEnd(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wantToPark.PushFront(id)
	return nil
}

func (s *InMemStorage) MoveToNotificationQueue(now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID, ok := s.wantToPark.PopBack()
	if !ok {
		return nil
	}

	s.notifiedQueue.PushFront(models.NewSubscription(userID, now))
	return nil
}

func (s *InMemStorage) NotifiedQueuePeekBack() (models.Subscription, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.notifiedQueue.PeekBack()
}

func (s *InMemStorage) NotifiedQueuePopBack() (models.Subscription, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.notifiedQueue.PopBack()
}
