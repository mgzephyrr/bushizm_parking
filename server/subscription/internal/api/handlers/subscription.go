package handlers

import (
	"subscription/internal/api"
	"time"

	"github.com/gofiber/fiber/v2"
)

type NewSubRequest struct {
	ID int `json:"id"`
}

func CreateSubscription(queue api.Queue) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		now := time.Now()
		req := new(NewSubRequest)

		if err := c.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid request body: "+err.Error())
		}

		id := req.ID
		sub, err := queue.AddSubToEnd(id, now)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(sub)
	}
}
