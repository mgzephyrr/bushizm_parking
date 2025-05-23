package routes

import (
	"subscription/internal/api"
	"subscription/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterCarEventsRoutes(router fiber.Router, queue api.Queue) {
	router.Post("/carevents", handlers.HandleCarEvent(queue))
}
