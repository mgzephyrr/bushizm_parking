package server

import (
	"context"
	"fmt"
	"log/slog"
	"subscription/internal/api"
	"subscription/internal/api/routes"

	"github.com/gofiber/fiber/v2"
)

type APIServer struct {
	server *fiber.App
	queue  api.Queue
}

func NewAPIServer(queue api.Queue) *APIServer {
	api := &APIServer{
		server: fiber.New(),
		queue:  queue,
	}

	apiVersion := api.server.Group("/api/v1")

	routes.RegisterSubsRoutes(apiVersion, api.queue)
	return api
}

func (a *APIServer) Run(port string) error {
	return a.server.Listen(":" + port)
}

func (a *APIServer) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down HTTP server...")
	if err := a.server.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}

	return nil
}
