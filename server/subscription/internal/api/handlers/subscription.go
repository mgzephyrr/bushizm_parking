package handlers

import (
	"subscription/internal/api"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateSubscription(queue api.Queue) func (c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		now := time.Now()
		
		id, err := getIntVariable(c, "id")
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}

		sub, err := queue.AddSubToEnd(id, now)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(sub)
	}
}
