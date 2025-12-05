package handler

import (
	"api-gateway-go/internal/config"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// TestHandler handles test routes
type TestHandler struct {
	cfg        *config.Config
	db         *config.Database
	redisMgr   *config.RedisManager
	livekitMgr *config.LiveKitManager
}

// NewTestHandler creates a new TestHandler
func NewTestHandler(
	cfg *config.Config,
	db *config.Database,
	redisMgr *config.RedisManager,
	livekitMgr *config.LiveKitManager,
) *TestHandler {
	return &TestHandler{
		cfg:        cfg,
		db:         db,
		redisMgr:   redisMgr,
		livekitMgr: livekitMgr,
	}
}

// Ping returns a pong response
// GET /test/ping
func (h *TestHandler) Ping(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.Map{
		"message": "pong",
	})
}

// Echo echoes back the request body
// POST /test/echo
func (h *TestHandler) Echo(c *fiber.Ctx) error {
	var body interface{}
	if err := c.BodyParser(&body); err != nil {
		body = string(c.Body())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"echo":    body,
		"headers": c.GetReqHeaders(),
		"query":   c.Queries(),
	})
}

// TestDatabase tests database connection
// GET /test/database
func (h *TestHandler) TestDatabase(c *fiber.Ctx) error {
	if h.db == nil {
		return utils.ErrorResponse(c, "Database not configured")
	}

	err := h.db.Health(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, "Database connection failed: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"status":  "connected",
		"message": "Database connection successful",
	})
}

// TestRedis tests Redis connection
// GET /test/redis
func (h *TestHandler) TestRedis(c *fiber.Ctx) error {
	if h.redisMgr == nil {
		return utils.ErrorResponse(c, "Redis not configured")
	}

	err := h.redisMgr.Health(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, "Redis connection failed: "+err.Error())
	}

	// Test set/get
	testKey := "test:connection"
	err = h.redisMgr.Set(c.Context(), testKey, "ok", 10)
	if err != nil {
		return utils.ErrorResponse(c, "Redis set failed: "+err.Error())
	}

	value, err := h.redisMgr.Get(c.Context(), testKey)
	if err != nil {
		return utils.ErrorResponse(c, "Redis get failed: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"status":    "connected",
		"message":   "Redis connection successful",
		"testValue": value,
	})
}

// TestLiveKit tests LiveKit connection
// GET /test/livekit
func (h *TestHandler) TestLiveKit(c *fiber.Ctx) error {
	if h.livekitMgr == nil {
		return utils.ErrorResponse(c, "LiveKit not configured")
	}

	err := h.livekitMgr.Health(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, "LiveKit connection failed: "+err.Error())
	}

	// List rooms to verify connection
	rooms, err := h.livekitMgr.ListRooms(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, "LiveKit list rooms failed: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"status":    "connected",
		"message":   "LiveKit connection successful",
		"roomCount": len(rooms),
	})
}

// TestAll tests all connections
// GET /test/all
func (h *TestHandler) TestAll(c *fiber.Ctx) error {
	results := fiber.Map{}

	// Test Database
	if h.db != nil {
		if err := h.db.Health(c.Context()); err == nil {
			results["database"] = "connected"
		} else {
			results["database"] = "failed: " + err.Error()
		}
	} else {
		results["database"] = "not configured"
	}

	// Test Redis
	if h.redisMgr != nil {
		if err := h.redisMgr.Health(c.Context()); err == nil {
			results["redis"] = "connected"
		} else {
			results["redis"] = "failed: " + err.Error()
		}
	} else {
		results["redis"] = "not configured"
	}

	// Test LiveKit
	if h.livekitMgr != nil {
		if err := h.livekitMgr.Health(c.Context()); err == nil {
			results["livekit"] = "connected"
		} else {
			results["livekit"] = "failed: " + err.Error()
		}
	} else {
		results["livekit"] = "not configured"
	}

	return utils.SuccessResponse(c, results)
}

// GetConfig returns configuration info (sanitized)
// GET /test/config
func (h *TestHandler) GetConfig(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.Map{
		"environment": h.cfg.Environment,
		"port":        h.cfg.Port,
		"database": fiber.Map{
			"host": h.cfg.MySQL.Host,
			"port": h.cfg.MySQL.Port,
			"name": h.cfg.MySQL.Database,
		},
		"redis": fiber.Map{
			"host": h.cfg.Redis.Host,
			"port": h.cfg.Redis.Port,
		},
		"livekit": fiber.Map{
			"host":       h.cfg.LiveKitHost,
			"configured": h.cfg.LiveKitAPIKey != "",
		},
		"features": fiber.Map{
			"roomDayDefaultTimeout": h.cfg.RoomDayDefaultTimeout,
			"egressLimit":           h.cfg.EgressLimit,
		},
	})
}

// GetAllNamespaces returns all socket namespaces
// GET /test/get/namespace
func (h *TestHandler) GetAllNamespaces(c *fiber.Ctx) error {
	// Placeholder - actual implementation would get from socket.io manager
	namespaces := []string{
		"/",
		"/mobile",
		"/notification",
		"/queue",
		"/newqueue",
	}

	return utils.SuccessResponse(c, namespaces)
}

// ClearRedisTestData clears Redis test data
// DELETE /test/redis/clear
func (h *TestHandler) ClearRedisTestData(c *fiber.Ctx) error {
	if h.redisMgr == nil {
		return utils.ErrorResponse(c, "Redis not configured")
	}

	// Clean up test data
	testKeys := []string{"test:key", "test:user:1", "test:messages"}
	for _, key := range testKeys {
		h.redisMgr.Delete(c.Context(), key)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Test data cleared from Redis",
	})
}

// GetMP4ProcessingQueue gets MP4 encoding queue status
// GET /test/mp4/queue
func (h *TestHandler) GetMP4ProcessingQueue(c *fiber.Ctx) error {
	if h.redisMgr == nil {
		return utils.ErrorResponse(c, "Redis not configured")
	}

	// Get queue items
	queueItems, err := h.redisMgr.Client().LRange(c.Context(), "mp4_encode_queue", 0, -1).Result()
	if err != nil {
		return utils.ErrorResponse(c, "Failed to get MP4 encoding queue status")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"queueSize":  len(queueItems),
		"queueItems": queueItems,
		"message":    "Queue items waiting for processing",
	})
}

// RemoveFromMP4ProcessingQueue removes a record from MP4 processing queue
// DELETE /test/mp4/queue/:recordId
func (h *TestHandler) RemoveFromMP4ProcessingQueue(c *fiber.Ctx) error {
	if h.redisMgr == nil {
		return utils.ErrorResponse(c, "Redis not configured")
	}

	recordID := c.Params("recordId")
	if recordID == "" {
		return utils.BadRequestResponse(c, "Record ID is required")
	}

	// Remove from processing set
	processingKey := "mp4_processing_set:" + recordID
	h.redisMgr.Delete(c.Context(), processingKey)

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Record " + recordID + " removed from processing set",
		"note":    "Cannot remove from FIFO queue directly - queue will be processed by workers",
	})
}

// ClearMP4ProcessingQueue clears the MP4 encoding queue
// DELETE /test/mp4/queue
func (h *TestHandler) ClearMP4ProcessingQueue(c *fiber.Ctx) error {
	if h.redisMgr == nil {
		return utils.ErrorResponse(c, "Redis not configured")
	}

	// Get current queue items for display
	queueItems, _ := h.redisMgr.Client().LRange(c.Context(), "mp4_encode_queue", 0, -1).Result()
	queueSize := len(queueItems)

	// Clear the entire FIFO queue
	h.redisMgr.Delete(c.Context(), "mp4_encode_queue")

	return utils.SuccessResponse(c, fiber.Map{
		"message":      "Cleared " + string(rune(queueSize)) + " records from MP4 encoding queue",
		"clearedItems": queueItems,
	})
}
