package middleware

import (
	"time"

	"api-gateway-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware creates a logging middleware
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request
		logger.HTTP(
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			latency,
		)

		return err
	}
}
