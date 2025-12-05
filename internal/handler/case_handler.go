package handler

import (
	"fmt"
	"strconv"

	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// CaseHandler handles case routes
type CaseHandler struct {
	caseService *service.CaseService
}

// NewCaseHandler creates a new CaseHandler
func NewCaseHandler(caseService *service.CaseService) *CaseHandler {
	return &CaseHandler{
		caseService: caseService,
	}
}

// CreateCase creates a new case
// POST /case/create
func (h *CaseHandler) CreateCase(c *fiber.Ctx) error {
	type CreateRequest struct {
		CaseID          int    `json:"caseId"`
		Service         int    `json:"service"`
		RoomID          int    `json:"roomId"`
		OperationNumber string `json:"operationNumber"`
		Status          string `json:"status"`
		HN              string `json:"hn"`
		PatientMobile   string `json:"patientMobile"`
		MobileCreated   string `json:"mobileCreated"`
		CaseType        string `json:"caseType"`
		UserName        string `json:"userName"`
		Organization    string `json:"organization"`
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	id, err := h.caseService.CreateCase(c.Context(), service.CreateCaseOptions{
		CaseID:          req.CaseID,
		Service:         req.Service,
		RoomID:          req.RoomID,
		OperationNumber: req.OperationNumber,
		Status:          req.Status,
		HN:              req.HN,
		PatientMobile:   req.PatientMobile,
		MobileCreated:   req.MobileCreated,
		CaseType:        req.CaseType,
		UserName:        req.UserName,
		Organization:    req.Organization,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"created": true,
	})
}

// GetCaseByID gets a case by ID
// GET /case/:id
func (h *CaseHandler) GetCaseByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	fmt.Println("Get Case Id", idStr)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid case ID")
	}

	caseData, err := h.caseService.GetCaseByID(c.Context(), uint(id))
	if err != nil {
		if err == service.ErrCaseNotFound {
			return utils.NotFoundResponse(c, "Case not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, caseData)
}

// GetCaseByCaseID gets a case by case ID
// GET /case/caseid/:caseId
func (h *CaseHandler) GetCaseByCaseID(c *fiber.Ctx) error {
	caseIDStr := c.Params("caseId")
	caseID, err := strconv.Atoi(caseIDStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid case ID")
	}

	caseData, err := h.caseService.GetCaseByCaseID(c.Context(), caseID)
	if err != nil {
		if err == service.ErrCaseNotFound {
			return utils.NotFoundResponse(c, "Case not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, caseData)
}

// GetCaseByRoomID gets a case by room ID
// GET /case/room/:roomId
func (h *CaseHandler) GetCaseByRoomID(c *fiber.Ctx) error {
	roomIDStr := c.Params("roomId")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid room ID")
	}

	caseData, err := h.caseService.GetCaseByRoomID(c.Context(), roomID)
	if err != nil {
		if err == service.ErrCaseNotFound {
			return utils.NotFoundResponse(c, "Case not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, caseData)
}

// UpdateCase updates a case
// PUT /case/:id
func (h *CaseHandler) UpdateCase(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid case ID")
	}

	type UpdateRequest struct {
		Status        string `json:"status"`
		HN            string `json:"hn"`
		PatientMobile string `json:"patientMobile"`
		CaseType      string `json:"caseType"`
		UserName      string `json:"userName"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err = h.caseService.UpdateCase(c.Context(), uint(id), req.Status, req.HN, req.PatientMobile, req.CaseType, req.UserName)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// UpdateCaseStatus updates case status
// PUT /case/status/:caseId
func (h *CaseHandler) UpdateCaseStatus(c *fiber.Ctx) error {
	caseIDStr := c.Params("caseId")
	caseID, err := strconv.Atoi(caseIDStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid case ID")
	}

	type UpdateRequest struct {
		Status string `json:"status"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err = h.caseService.UpdateCaseStatus(c.Context(), caseID, req.Status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// GetCaseHistory gets case history
// GET /case/history
func (h *CaseHandler) GetCaseHistory(c *fiber.Ctx) error {
	serviceStr := c.Query("service", "0")
	limitStr := c.Query("limit", "100")
	offsetStr := c.Query("offset", "0")

	service, _ := strconv.Atoi(serviceStr)
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	cases, err := h.caseService.GetCaseHistory(c.Context(), service, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, cases)
}

// GetCaseHistoryCount gets case history count
// GET /case/historycount
func (h *CaseHandler) GetCaseHistoryCount(c *fiber.Ctx) error {
	serviceStr := c.Query("service", "0")
	service, _ := strconv.Atoi(serviceStr)

	count, err := h.caseService.GetCaseHistoryCount(c.Context(), service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"count": count,
	})
}

// GetCasesByService gets cases by service
// GET /case/service/:service
func (h *CaseHandler) GetCasesByService(c *fiber.Ctx) error {
	serviceStr := c.Params("service")
	service, err := strconv.Atoi(serviceStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid service ID")
	}

	cases, err := h.caseService.GetCasesByService(c.Context(), service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, cases)
}

// GetRoomName gets room name by case ID
// GET /case/roomname
func (h *CaseHandler) GetRoomName(c *fiber.Ctx) error {
	caseIDStr := c.Query("caseId")
	serviceStr := c.Query("service")

	caseID, err := strconv.Atoi(caseIDStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid case ID")
	}

	service, err := strconv.Atoi(serviceStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid service ID")
	}

	roomName, err := h.caseService.GetRoomNameByCaseID(c.Context(), caseID, service)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"roomName": roomName,
	})
}

// DeleteCase deletes a case
// DELETE /case/:id
func (h *CaseHandler) DeleteCase(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid case ID")
	}

	err = h.caseService.DeleteCase(c.Context(), uint(id))
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"deleted": true,
	})
}
