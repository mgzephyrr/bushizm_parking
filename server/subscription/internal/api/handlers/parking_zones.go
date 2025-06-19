package handlers

import (
	"subscription/internal/api"

	"github.com/gofiber/fiber/v2"
)

type SpotsResponse struct {
	SpotsNumber int `json:"spots_number"`
}

func GetAvailableSpots(parking api.ParkingService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		spots, err := parking.CheckAvailableSpots()
		if err != nil {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				"Error while retreiving number of available spots",
			)
		}

		return c.Status(fiber.StatusOK).JSON(SpotsResponse{
			SpotsNumber: spots,
		})
	}
}
