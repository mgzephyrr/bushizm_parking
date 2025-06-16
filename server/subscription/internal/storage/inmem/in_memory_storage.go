package inmem

import (
	"context"
	"fmt"
	"subscription/internal/api"
	"subscription/internal/api/models"
	"sync"
	"time"

	"github.com/gammazero/deque"
)

const WORKERS_COUNT = 1

type InMemStorage struct {
	wantToPark       deque.Deque[int]
	notifiedQueue    deque.Deque[models.Subscription]
	mu               sync.Mutex
	maxQueueSize     int
	inQueue          map[int]struct{}
	lastDequeueTimes []time.Time
}

func NewInMemStorage(ctx context.Context, maxSize int) *InMemStorage {
	storage := &InMemStorage{
		wantToPark:    deque.Deque[int]{},
		notifiedQueue: deque.Deque[models.Subscription]{},
		maxQueueSize:  maxSize,
		inQueue:       make(map[int]struct{}),
	}

	for i := range WORKERS_COUNT {
		go NewQueueWorker(i+1, storage).Process(ctx)
	}

	return storage
}

func (s *InMemStorage) GetAllQueue() []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]int, s.wantToPark.Len())
	for i := 0; i < s.wantToPark.Len(); i++ {
		result[i] = s.wantToPark.At(i)
	}
	return result
}

func (s *InMemStorage) AddSubToEnd(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.wantToPark.Len() >= s.maxQueueSize {
		return api.ErrQueueFull
	}

	if _, exists := s.inQueue[id]; exists {
		return api.ErrAlreadyInQueue
	}

	s.wantToPark.PushBack(id)
	s.inQueue[id] = struct{}{}
	return nil
}

func (s *InMemStorage) MoveToNotificationQueue(now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.wantToPark.Len() == 0 {
		return nil
	}

	userID := s.wantToPark.PopFront()
	delete(s.inQueue, userID)
	s.notifiedQueue.PushBack(models.NewSubscription(userID, now))

	s.lastDequeueTimes = append(s.lastDequeueTimes, now)
	if len(s.lastDequeueTimes) > s.maxQueueSize {
		s.lastDequeueTimes = s.lastDequeueTimes[1:] // храним последние 10
	}

	return nil
}

func (s *InMemStorage) EstimateWaitTime(position int) time.Duration {
    s.mu.Lock()
    defer s.mu.Unlock()

    times := s.lastDequeueTimes
    if len(times) < 2 {
        return 0
    }

    var total time.Duration
    for i := 1; i < len(times); i++ {
        total += times[i].Sub(times[i-1])
    }

    avgInterval := total / time.Duration(len(times)-1)
    return avgInterval * time.Duration(position-1)
}

func (s *InMemStorage) NotifiedQueuePeekBack() (models.Subscription, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.notifiedQueue.Len() == 0 {
		return models.Subscription{}, false
	}
	return s.notifiedQueue.Back(), true
}

func (s *InMemStorage) NotifiedQueuePopBack() (models.Subscription, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.notifiedQueue.Len() == 0 {
		return models.Subscription{}, false
	}
	return s.notifiedQueue.PopBack(), true
}

func (s *InMemStorage) GetUserPosition(userID int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < s.wantToPark.Len(); i++ {
		if s.wantToPark.At(i) == userID {
			return i + 1, nil
		}
	}
	return -1, fmt.Errorf("user not in queue")
}
