package handlers

import (
	"context"
	"log/slog"
	"subscription/internal/api"
	"subscription/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CarEvent struct {
	Data    CarEventData `json:"data"`
	Code    int          `json:"code"`
	Message string       `json:"message"`
}

type CarEventData struct {
	EventID               string           `json:"event_id"`
	AccessPointID         string           `json:"access_point_id"`
	Direction             models.Direction `json:"direction"`
	CreatedAt             time.Time        `json:"timestamp"`
	LicensePlate          string           `json:"license_plate"`
	RecognitionConfidence float32          `json:"recognition_confidence"`
}

func HandleCarEvent(ctx context.Context, queue api.Queue) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		carEvent := new(CarEvent)

		now := time.Now()

		if err := c.BodyParser(carEvent); err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid request body: "+err.Error())
		}

		if carEvent.Data.Direction == models.DirectionIn {
			return handleCarIn(ctx, queue, carEvent, now)
		}

		return handleCarOut(ctx, queue, carEvent, now)
	}
}

func handleCarIn(ctx context.Context, queue api.Queue, carEvent *CarEvent, now time.Time) error {
	slog.Info("Processing Car In...")
	queue.NotifiedQueuePopBack(ctx)

	return nil
}

func handleCarOut(ctx context.Context, queue api.Queue, carEvent *CarEvent, now time.Time) error {
	slog.Info("Processing Car Out...")
	err := queue.MoveToNotificationQueue(ctx, now)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"Internal error while adding entry to notification queue: "+err.Error(),
		)
	}

	return nil
}
