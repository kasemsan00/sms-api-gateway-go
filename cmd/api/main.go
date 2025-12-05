package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	// In Docker/production, prefer environment variables passed to the container
	// but still allow .env file if mounted as volume
	env := os.Getenv("ENVIRONMENT")
	if err := godotenv.Load(); err != nil {
		// Only warn in development mode, in production it's expected to use env vars
		if env == "" || env == "development" {
			log.Println("Warning: .env file not found, using system environment variables")
		}
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "SMS API Gateway",
	})

	// Add middleware
	app.Use(logger.New())

	// Setup routes
	setupRoutes(app)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	listenAddr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Server starting on http://localhost%s\n", listenAddr)

	if err := app.Listen(listenAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(app *fiber.App) {
	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to SMS API Gateway",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// API v1 routes group
	api := app.Group("/api/v1")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API v1",
		})
	})
}
