package workerpool

import (
	"context"
	"fmt"
	"log/slog"
	"subscription/internal/api"
	"time"
)

const TIMEOUT time.Duration = 2 * time.Second

type QueueWorker struct {
	WorkerID            int
	ticker              *time.Ticker
	storage             api.Queue
	notificationService api.NotificationService
}

func NewQueueWorker(workerID int, storage api.Queue, notif api.NotificationService) *QueueWorker {
	return &QueueWorker{
		WorkerID:            workerID,
		ticker:              time.NewTicker(TIMEOUT),
		storage:             storage,
		notificationService: notif,
	}
}

func (w *QueueWorker) Process(ctx context.Context) {
	for {
		now := time.Now()
		select {
		case <-w.ticker.C:
			for {
				sub, ok := w.storage.NotifiedQueuePeekBack(ctx)
				if !ok {
					break
				}

				if sub.ExpiresAt.Before(now) {
					w.storage.NotifiedQueuePopBack(ctx)
					w.storage.MoveToNotificationQueue(ctx, now)
					break
				}

				err := w.notificationService.Notify(sub)
				if err != nil {
					if err.Error() == "max notification attempts reached" {
						slog.Info("Max notify attempts reached, removing subscription", slog.Int("userID", sub.UserID))
						w.storage.NotifiedQueuePopBack(ctx)
					} else {
						slog.Error("Error while notifying", slog.String("error", err.Error()))
					}
					break
				}

				w.storage.NotifiedQueuePopBack(ctx)
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
