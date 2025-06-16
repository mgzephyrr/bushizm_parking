package routes

import (
	"subscription/internal/api"
	"subscription/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterSubsRoutes(router fiber.Router, queue api.Queue) {
	group := router.Group("/subscriptions")

	group.Post("/subscribe", handlers.CreateSubscription(queue))
	group.Get("/position", handlers.GetUserQueuePosition(queue))
}
