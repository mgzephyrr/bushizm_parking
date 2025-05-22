package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"subscription/internal/api/server"
	"subscription/internal/storage/inmem"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := inmem.NewInMemStorage(ctx)

	server := server.NewAPIServer(queue)
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
