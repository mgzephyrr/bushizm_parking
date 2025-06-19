package server

import (
	"context"
	"fmt"
	"log/slog"
	"subscription/internal/api"
	"subscription/internal/api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type APIServer struct {
	server              *fiber.App
	queue               api.Queue
	parkingService      api.ParkingService
}

func NewAPIServer(ctx context.Context, queue api.Queue, parking api.ParkingService) *APIServer {
	api := &APIServer{
		server:              fiber.New(),
		queue:               queue,
		parkingService:      parking,
	}

	api.server.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:5174, http://localhost:8000, http://localhost:8001",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	api.server.Use(func(c *fiber.Ctx) error {
		bodyBytes := c.Body()
		slog.Info(fmt.Sprintf("Request Body: %s", string(bodyBytes)))
		return c.Next()
	})

	apiVersion := api.server.Group("/api/v1")

	routes.RegisterSubsRoutes(ctx, apiVersion, api.queue, api.parkingService)
	routes.RegisterCarEventsRoutes(ctx, apiVersion, api.queue)
	routes.RegisterParkingZonesRoutes(apiVersion, api.parkingService)
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
