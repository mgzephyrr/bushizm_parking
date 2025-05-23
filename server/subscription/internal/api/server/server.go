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
	server *fiber.App
	queue  api.Queue
}

func NewAPIServer(queue api.Queue) *APIServer {
	api := &APIServer{
		server: fiber.New(),
		queue:  queue,
	}

	api.server.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	api.server.Use(func(c *fiber.Ctx) error {
		bodyBytes := c.Body()
		slog.Info(fmt.Sprintf("Request Body: %s", string(bodyBytes)))
		return c.Next()
	})

	apiVersion := api.server.Group("/api/v1")

	routes.RegisterSubsRoutes(apiVersion, api.queue)
	routes.RegisterCarEventsRoutes(apiVersion, api.queue)
	routes.RegisterParkingZonesRoutes(apiVersion)
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
