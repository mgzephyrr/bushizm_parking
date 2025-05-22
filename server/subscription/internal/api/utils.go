package api

import (
	"subscription/internal/api/models"
	"time"
)

type Queue interface {
	AddSubToEnd(int, time.Time) (models.Subscription, error)
}
