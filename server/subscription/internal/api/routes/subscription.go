package routes

import (
	"subscription/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/idsulik/go-collections/deque"
)

func RegisterSubsRoutes(router fiber.Router, queue *deque.Deque[string]) {
	group := router.Group("/subscriptions")

	group.Post("/subscribe", handlers.CreateSubscription(queue))
}
