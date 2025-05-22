package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getIntVariable(c *fiber.Ctx, key string) (int, error) {
	varStr := c.Params(key)

	id, err := strconv.ParseInt(varStr, 10, 32)
	if err != nil {
		return 0, errors.New("unable to parse user ID as integer")
	}

	return int(id), nil
}
