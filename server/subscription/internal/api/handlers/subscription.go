package handlers

import (
	"fmt"
	"log/slog"
	"subscription/internal/api"

	"github.com/gofiber/fiber/v2"
)

type NewSubRequest struct {
	ID int `json:"id"`
}

type ZoneInfo struct {
	Data    ZoneInfoData `json:"data"`
	Code    int          `json:"code"`
	Message string       `json:"message"`
}

type ZoneInfoData struct {
	ParkingZoneID  string `json:"parking_zone_id"`
	Name           string `json:"name"`
	AvailableSpots int    `json:"available_spots"`
	Comment        string `json:"comment"`
}

func CreateSubscription(queue api.Queue) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		spotsNumber, err := CheckAvailableSpots()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if spotsNumber > 0 {
			return fiber.NewError(fiber.StatusNotAcceptable, "There are free spots on parking!")
		}

		req := new(NewSubRequest)

		if err := c.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid request body: "+err.Error())
		}

		err = queue.AddSubToEnd(req.ID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		slog.Info(fmt.Sprintf("%v", queue.GetAllQueue()))
		return c.Status(fiber.StatusOK).SendString("Added to queue")
	}
}
