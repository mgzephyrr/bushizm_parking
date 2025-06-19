package routes

import (
	"subscription/internal/api"
	"subscription/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterParkingZonesRoutes(router fiber.Router, parking api.Parking) {
	router.Get("/spotsnumber", handlers.GetAvailableSpots(parking))
}
