package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/idsulik/go-collections/deque"
)

func CreateSubscription(queue *deque.Deque[string]) func (c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sub := "Added new sub!"
		queue.PushFront(sub)

		return c.SendString("Added you to queue")
	}
}
