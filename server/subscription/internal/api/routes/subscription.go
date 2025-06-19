package routes

import (
	"context"
	"subscription/internal/api"
	"subscription/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterSubsRoutes(ctx context.Context, router fiber.Router, queue api.Queue, parking api.ParkingService) {
	group := router.Group("/subscriptions")

	group.Post("/subscribe", handlers.CreateSubscription(ctx, queue, parking))
	group.Get("/position", handlers.GetUserQueuePosition(ctx, queue))
}
