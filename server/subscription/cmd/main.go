package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"subscription/internal/api/server"
	"subscription/internal/notificationapi"
	"subscription/internal/parkingapi"
	"subscription/internal/storage/inmem"
	"subscription/internal/storage/workerpool"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	queueSize    = 10
	workersCount = 1
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notifService := notificationapi.NewNotificationAPI()
	queue := inmem.NewInMemStorage(ctx, queueSize)
	parking := parkingapi.NewParkingAPI()

	for i := range workersCount {
		go workerpool.NewQueueWorker(i+1, queue, notifService).Process(ctx)
	}

	server := server.NewAPIServer(queue, parking)
	serverErr := make(chan error, 1)

	go func() {
		if err := server.Run("8080"); err != nil {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		slog.Info("Received shutdown signal")
	case err := <-serverErr:
		slog.Error(fmt.Sprintf("Server error: %s", err.Error()))
	}
	cancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error(fmt.Sprintf("Shutdown error: %s", err.Error()))
	}

	slog.Info("Server stopped")
}
