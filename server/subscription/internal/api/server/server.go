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
	server  *fiber.App
	queue   api.Queue
	parking api.Parking
}

func NewAPIServer(queue api.Queue, parking api.Parking) *APIServer {
	api := &APIServer{
		server:  fiber.New(),
		queue:   queue,
		parking: parking,
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

	routes.RegisterSubsRoutes(apiVersion, api.queue, api.parking)
	routes.RegisterCarEventsRoutes(apiVersion, api.queue)
	routes.RegisterParkingZonesRoutes(apiVersion, api.parking)
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
