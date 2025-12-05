package handler

import (
	"context"
	"time"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// SystemHandler handles system routes
type SystemHandler struct {
	db      *config.Database
	redis   *config.RedisManager
	livekit *config.LiveKitManager
	crontab *service.CrontabService
	cfg     *config.Config
}

// NewSystemHandler creates a new SystemHandler
func NewSystemHandler(db *config.Database, redis *config.RedisManager, livekit *config.LiveKitManager, crontab *service.CrontabService, cfg *config.Config) *SystemHandler {
	return &SystemHandler{
		db:      db,
		redis:   redis,
		livekit: livekit,
		crontab: crontab,
		cfg:     cfg,
	}
}

// Root handles root route
// GET /
func (h *SystemHandler) Root(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to SMS API Gateway",
		"version": "1.0.0",
		"status":  "running",
	})
}

// HealthCheck handles health check
// GET /health
func (h *SystemHandler) HealthCheck(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	status := "healthy"
	details := make(map[string]string)

	// Check database
	if h.db != nil {
		if err := h.db.Health(ctx); err != nil {
			status = "unhealthy"
			details["database"] = err.Error()
		} else {
			details["database"] = "connected"
		}
	}

	// Check Redis
	if h.redis != nil {
		if err := h.redis.Health(ctx); err != nil {
			status = "unhealthy"
			details["redis"] = err.Error()
		} else {
			details["redis"] = "connected"
		}
	}

	// Check LiveKit
	if h.livekit != nil {
		if err := h.livekit.Health(ctx); err != nil {
			details["livekit"] = err.Error()
		} else {
			details["livekit"] = "connected"
		}
	}

	return c.JSON(fiber.Map{
		"status":  status,
		"details": details,
	})
}

// GetStatus handles status route
// GET /status
func (h *SystemHandler) GetStatus(c *fiber.Ctx) error {
	cronStatus := h.crontab.GetStatus()

	return c.JSON(fiber.Map{
		"status":      "running",
		"environment": h.cfg.Environment,
		"port":        h.cfg.Port,
		"cron":        cronStatus,
	})
}

// GetServiceInfo handles service info
// GET /service
func (h *SystemHandler) GetServiceInfo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":        "SMS API Gateway",
		"version":     "1.0.0",
		"environment": h.cfg.Environment,
		"apiURL":      h.cfg.APIURL,
		"features": fiber.Map{
			"smsEnabled":    h.cfg.SMSEnable,
			"autoCloseRoom": h.cfg.AutoCloseRoom,
			"egressLimit":   h.cfg.EgressLimit,
		},
	})
}

// GetNamespaces gets socket namespaces
// GET /namespace
func (h *SystemHandler) GetNamespaces(c *fiber.Ctx) error {
	// Placeholder for socket namespaces
	return utils.SuccessResponse(c, fiber.Map{
		"namespaces": []string{
			"/",
			"/mobile",
			"/notification",
			"/queue",
			"/newqueue",
		},
	})
}

// AddLog adds a log entry
// POST /log
func (h *SystemHandler) AddLog(c *fiber.Ctx) error {
	type LogRequest struct {
		Level   string      `json:"level"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	var req LogRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Log is handled by the logging system
	return utils.SuccessResponse(c, fiber.Map{
		"logged": true,
	})
}
