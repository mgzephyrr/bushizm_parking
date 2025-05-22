package inmem

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

const TIMEOUT time.Duration = time.Second

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

				if !ok || sub.ExpiresAt.Before(now) {
					break
				}

				err := sub.Notify()
				if err != nil {
					slog.Error("Error while notifying", slog.String("error", err.Error()))

					if sub.AttemptsLeft > 0 {
						slog.Info(fmt.Sprintf("Sub with userID = %d lost 1 attempt", sub.UserID))
						sub.AttemptsLeft -= 1
						break
					}
				}

				w.storage.NotifiedQueuePopBack()
			}
		case <-ctx.Done():
			w.Shutdown()
		}
	}
}

func (w *QueueWorker) Shutdown() {
	w.ticker.Stop()
	slog.Info(fmt.Sprintf("QueueWorker %d was successfully shutdown", w.WorkerID))
}
