package routes

import (
	"subscription/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterParkingZonesRoutes(router fiber.Router) {
	router.Get("/spotsnumber", handlers.GetAvailableSpots)
}
