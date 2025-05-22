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

func CreateSubscription(queue api.Queue) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := new(NewSubRequest)

		if err := c.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid request body: "+err.Error())
		}

		err := queue.AddSubToEnd(req.ID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		
		slog.Info(fmt.Sprintf("%v", queue.GetAllQueue()))
		return c.Status(fiber.StatusOK).SendString("Added to queue")
	}
}
