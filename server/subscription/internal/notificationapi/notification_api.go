package notificationapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"subscription/internal/models"
	"time"
)

type NotificationAPI struct{}

func NewNotificationAPI() *NotificationAPI {
	return &NotificationAPI{}
}

func (*NotificationAPI) Notify(sub models.Subscription) error {
	err := sub.Notify()
	if err != nil {
		return fmt.Errorf("error while notificating %w", err)
	}
	
	data := map[string]string{"user_id": strconv.Itoa(sub.UserID)}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("Error marshaling JSON", slog.String("error", err.Error()))
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://notification:8001/queue",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		slog.Error("Error while creating request", slog.String("error", err.Error()))
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	defer client.CloseIdleConnections()

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error while requesting service", slog.String("error", err.Error()))
		return err
	}
	defer func() {
		if resp == nil || resp.Body == nil {
			slog.Info("Could not get a response")
			return
		}
		if err := resp.Body.Close(); err != nil {
			slog.Error("Error closing response body", slog.String("error", err.Error()))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Did not receive code 200 from service")
		return fmt.Errorf("did not receive code 200 from service")
	}

	return nil
}
