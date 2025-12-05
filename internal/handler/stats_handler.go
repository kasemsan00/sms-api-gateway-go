package handler

import (
	"strconv"

	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// StatsHandler handles statistics routes
type StatsHandler struct {
	statsService *service.StatsService
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetSummary gets statistics summary
// GET /stats/summary
func (h *StatsHandler) GetSummary(c *fiber.Ctx) error {
	serviceStr := c.Query("service", "0")
	service, _ := strconv.Atoi(serviceStr)

	summary, err := h.statsService.GetSummary(c.Context(), service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, summary)
}

// GetDeviceStats gets device statistics
// GET /stats/device
func (h *StatsHandler) GetDeviceStats(c *fiber.Ctx) error {
	serviceStr := c.Query("service", "0")
	service, _ := strconv.Atoi(serviceStr)

	stats, err := h.statsService.GetDeviceStats(c.Context(), service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, stats)
}

// GetTypeStats gets type statistics
// GET /stats/type
func (h *StatsHandler) GetTypeStats(c *fiber.Ctx) error {
	serviceStr := c.Query("service", "0")
	service, _ := strconv.Atoi(serviceStr)

	stats, err := h.statsService.GetTypeStats(c.Context(), service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, stats)
}

// GetUserStats gets user statistics
// GET /stats/user
func (h *StatsHandler) GetUserStats(c *fiber.Ctx) error {
	serviceStr := c.Query("service", "0")
	service, _ := strconv.Atoi(serviceStr)

	stats, err := h.statsService.GetUserStats(c.Context(), service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, stats)
}

// GetCaseStats gets case statistics
// GET /stats/case
func (h *StatsHandler) GetCaseStats(c *fiber.Ctx) error {
	serviceStr := c.Query("service", "0")
	service, _ := strconv.Atoi(serviceStr)

	stats, err := h.statsService.GetCaseStats(c.Context(), service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, stats)
}

// GetDailyStats gets daily statistics
// GET /stats/daily
func (h *StatsHandler) GetDailyStats(c *fiber.Ctx) error {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	serviceStr := c.Query("service", "0")
	service, _ := strconv.Atoi(serviceStr)

	if startDate == "" || endDate == "" {
		return utils.BadRequestResponse(c, "Start date and end date required")
	}

	stats, err := h.statsService.GetDailyStats(c.Context(), startDate, endDate, service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, stats)
}

// GetMonthlyStats gets monthly statistics
// GET /stats/monthly
func (h *StatsHandler) GetMonthlyStats(c *fiber.Ctx) error {
	yearStr := c.Query("year")
	serviceStr := c.Query("service", "0")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid year")
	}

	service, _ := strconv.Atoi(serviceStr)

	stats, err := h.statsService.GetMonthlyStats(c.Context(), year, service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, stats)
}

// GetAll gets all statistics
// GET /stats/all
func (h *StatsHandler) GetAll(c *fiber.Ctx) error {
	serviceStr := c.Query("service", "0")
	service, _ := strconv.Atoi(serviceStr)

	summary, err := h.statsService.GetSummary(c.Context(), service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	deviceStats, _ := h.statsService.GetDeviceStats(c.Context(), service)
	typeStats, _ := h.statsService.GetTypeStats(c.Context(), service)
	userStats, _ := h.statsService.GetUserStats(c.Context(), service)
	caseStats, _ := h.statsService.GetCaseStats(c.Context(), service)

	return utils.SuccessResponse(c, fiber.Map{
		"summary": summary,
		"device":  deviceStats,
		"type":    typeStats,
		"user":    userStats,
		"case":    caseStats,
	})
}
