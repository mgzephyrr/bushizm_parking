package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"subscription/internal/api"
	"time"

	"github.com/gofiber/fiber/v2"
)

const AUTH_COOKIE = "access_token"

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
		authCookie := c.Cookies(AUTH_COOKIE)
		if authCookie == "" {
			return fiber.NewError(fiber.StatusUnauthorized)
		}

		spotsNumber, err := CheckAvailableSpots()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if spotsNumber > 0 {
			return fiber.NewError(fiber.StatusNotAcceptable, "There are free spots on parking!")
		}

		userID, err := GetUserID(authCookie)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = queue.AddSubToEnd(userID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		slog.Info(fmt.Sprintf("%v", queue.GetAllQueue()))
		return c.Status(fiber.StatusOK).SendString("Added to queue")
	}
}

type userIDResponse struct {
	UserID int `json:"json_id"`
}

func GetUserID(authCookie string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/api/v1/zones/", nil)
	if err != nil {
		slog.Error("Error while creating request", slog.String("error", err.Error()))
		return 0, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	defer client.CloseIdleConnections()

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error while requesting service", slog.String("error", err.Error()))
	}
	defer func() {
		if resp == nil {
			slog.Info("Could not get a response")
		}

		if err := resp.Body.Close(); err != nil {
			slog.Error("Error closing response body", slog.String("error", err.Error()))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Did not receive code 200 from service")
		return 0, fmt.Errorf("did not receive code 200 from service")
	}

	var response userIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("Error decoding JSON", slog.String("error", err.Error()))
		return 0, err
	}

	return response.UserID, nil
}
