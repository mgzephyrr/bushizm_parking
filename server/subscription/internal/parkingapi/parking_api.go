package parkingapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const parkingZoneID = "672681d0-5267-49de-b853-885576d0201a"

type ParkingAPI struct{}

func NewParkingAPI() *ParkingAPI {
	return &ParkingAPI{}
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

func (pAPI ParkingAPI) CheckAvailableSpots() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://192.168.1.28:8999/api/v1/zones/"+parkingZoneID, nil)
	if err != nil {
		slog.Error("Error while creating request", slog.String("error", err.Error()))
		return 0, err
	}
	slog.Info("TOKEN: " + os.Getenv("BEARER_TOKEN"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))

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

	var zoneInfo ZoneInfo
	if err := json.NewDecoder(resp.Body).Decode(&zoneInfo); err != nil {
		slog.Error("Error decoding JSON", slog.String("error", err.Error()))
		return 0, err
	}

	return zoneInfo.Data.AvailableSpots, nil
}
