package routes

import (
	"context"
	"subscription/internal/api"
	"subscription/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterCarEventsRoutes(ctx context.Context, router fiber.Router, queue api.Queue) {
	router.Post("/carevents", handlers.HandleCarEvent(ctx, queue))
}
