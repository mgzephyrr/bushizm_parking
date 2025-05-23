package inmem

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

const TIMEOUT time.Duration = 2 * time.Second

type QueueWorker struct {
	WorkerID int
	ticker   *time.Ticker
	storage  *InMemStorage
}

func NewQueueWorker(workerID int, storage *InMemStorage) *QueueWorker {
	return &QueueWorker{
		WorkerID: workerID,
		ticker:   time.NewTicker(TIMEOUT),
		storage:  storage,
	}
}

func (w *QueueWorker) Process(ctx context.Context) {
	for {
		now := time.Now()
		select {
		case <-w.ticker.C:
			for {
				sub, ok := w.storage.NotifiedQueuePeekBack()
				if !ok {
					break
				}
				slog.Info("User ID", slog.Int("ID", sub.UserID))
				if sub.ExpiresAt.Before(now) {
					w.storage.NotifiedQueuePopBack()
					w.storage.MoveToNotificationQueue(now)
					break
				}

				err := sub.Notify()
				if err != nil {
					slog.Error("Error while notifying", slog.String("error", err.Error()))
				}

				w.storage.NotifiedQueuePopBack()
			}
		case <-ctx.Done():
			w.Shutdown()
			return
		}
	}
}

func (w *QueueWorker) Shutdown() {
	w.ticker.Stop()
	slog.Info(fmt.Sprintf("QueueWorker %d was successfully shutdown", w.WorkerID))
}
