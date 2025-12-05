package handler

import (
	"strconv"

	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// RecordHandler handles recording routes
type RecordHandler struct {
	recordService *service.RecordService
}

// NewRecordHandler creates a new RecordHandler
func NewRecordHandler(recordService *service.RecordService) *RecordHandler {
	return &RecordHandler{
		recordService: recordService,
	}
}

// StartRecord starts a recording
// POST /record/start
func (h *RecordHandler) StartRecord(c *fiber.Ctx) error {
	type StartRequest struct {
		Room       string `json:"room"`
		RecordType string `json:"recordType"`
		FilePath   string `json:"filePath"`
	}

	var req StartRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.Room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	info, err := h.recordService.StartRecord(c.Context(), service.StartRecordOptions{
		Room:       req.Room,
		RecordType: req.RecordType,
		FilePath:   req.FilePath,
	})
	if err != nil {
		if err == service.ErrEgressLimit {
			return utils.ErrorResponseWithStatus(c, fiber.StatusTooManyRequests, "Egress limit reached")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"egressId": info.EgressId,
		"room":     req.Room,
		"status":   info.Status.String(),
	})
}

// StopRecord stops a recording
// POST /record/stop
func (h *RecordHandler) StopRecord(c *fiber.Ctx) error {
	type StopRequest struct {
		RecordID string `json:"recordId"`
		EgressID string `json:"egressId"`
	}

	var req StopRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	recordID := req.RecordID
	if recordID == "" {
		recordID = req.EgressID
	}

	if recordID == "" {
		return utils.BadRequestResponse(c, "Record ID or Egress ID required")
	}

	info, err := h.recordService.StopRecord(c.Context(), recordID)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"egressId": info.EgressId,
		"status":   info.Status.String(),
		"stopped":  true,
	})
}

// ListEgress lists active egresses
// GET /record/listegress
func (h *RecordHandler) ListEgress(c *fiber.Ctx) error {
	room := c.Query("room", "")

	egresses, err := h.recordService.ListEgress(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, egresses)
}

// StopAllActive stops all active recordings
// POST /record/stopall
func (h *RecordHandler) StopAllActive(c *fiber.Ctx) error {
	stopped, err := h.recordService.StopAllActiveRecords(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"stopped": stopped,
		"count":   len(stopped),
	})
}

// GetRecordDetail gets record detail by ID
// GET /record/detail/:id
func (h *RecordHandler) GetRecordDetail(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid record ID")
	}

	record, err := h.recordService.GetRecordDetail(c.Context(), id)
	if err != nil {
		if err == service.ErrRecordNotFound {
			return utils.NotFoundResponse(c, "Record not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, record)
}

// GetRecordByRoom gets records by room
// GET /record/room
func (h *RecordHandler) GetRecordByRoom(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	records, err := h.recordService.GetRecordByRoom(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, records)
}

// GetFileHistory gets recording file history
// GET /record/filehistory
func (h *RecordHandler) GetFileHistory(c *fiber.Ctx) error {
	room := c.Query("room", "")

	records, err := h.recordService.GetFileHistory(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, records)
}

// CheckEgressAvailable checks if egress is available
// GET /record/available
func (h *RecordHandler) CheckEgressAvailable(c *fiber.Ctx) error {
	available, err := h.recordService.CheckEgressAvailable(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"available": available,
	})
}

// GetRecordQueue gets the record queue
// GET /record/queue
func (h *RecordHandler) GetRecordQueue(c *fiber.Ctx) error {
	queue, err := h.recordService.GetRecordQueue(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, queue)
}

// GetActiveRecordCount gets count of active recordings
// GET /record/activecount
func (h *RecordHandler) GetActiveRecordCount(c *fiber.Ctx) error {
	count, err := h.recordService.GetActiveRecordCount(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"count": count,
	})
}
