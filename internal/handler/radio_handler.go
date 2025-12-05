package handler

import (
	"strconv"

	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// RadioHandler handles radio device and location routes
type RadioHandler struct {
	radioService *service.RadioService
}

// NewRadioHandler creates a new RadioHandler
func NewRadioHandler(radioService *service.RadioService) *RadioHandler {
	return &RadioHandler{
		radioService: radioService,
	}
}

// ListDevices lists radio devices
// GET /radio/devices
func (h *RadioHandler) ListDevices(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	devices, err := h.radioService.GetDevices(c.Context(), limit)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, devices)
}

// GetDeviceByID gets a device by ID
// GET /radio/device/:id
func (h *RadioHandler) GetDeviceByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid device ID")
	}

	device, err := h.radioService.GetDeviceByID(c.Context(), id)
	if err != nil {
		if err == service.ErrRadioDeviceNotFound {
			return utils.NotFoundResponse(c, "Device not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, device)
}

// GetDeviceByDeviceID gets a device by device ID
// GET /radio/device/deviceid/:deviceId
func (h *RadioHandler) GetDeviceByDeviceID(c *fiber.Ctx) error {
	deviceID := c.Params("deviceId")
	if deviceID == "" {
		return utils.BadRequestResponse(c, "Device ID required")
	}

	device, err := h.radioService.GetDeviceByDeviceID(c.Context(), deviceID)
	if err != nil {
		if err == service.ErrRadioDeviceNotFound {
			return utils.NotFoundResponse(c, "Device not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, device)
}

// CreateDevice creates a new radio device
// POST /radio/device
func (h *RadioHandler) CreateDevice(c *fiber.Ctx) error {
	type CreateRequest struct {
		DeviceID string `json:"deviceId"`
		RadioNo  string `json:"radioNo"`
		Channel  string `json:"channel"`
		Status   string `json:"status"`
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.DeviceID == "" {
		return utils.BadRequestResponse(c, "Device ID required")
	}

	id, err := h.radioService.CreateDevice(c.Context(), req.DeviceID, req.RadioNo, req.Channel, req.Status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"created": true,
	})
}

// UpdateDevice updates a radio device
// PUT /radio/device/:id
func (h *RadioHandler) UpdateDevice(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid device ID")
	}

	type UpdateRequest struct {
		RadioNo string `json:"radioNo"`
		Channel string `json:"channel"`
		Status  string `json:"status"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err = h.radioService.UpdateDevice(c.Context(), id, req.RadioNo, req.Channel, req.Status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// UpdateDeviceLocation updates device location
// PUT /radio/device/location
func (h *RadioHandler) UpdateDeviceLocation(c *fiber.Ctx) error {
	type UpdateRequest struct {
		DeviceID  string  `json:"deviceId"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.DeviceID == "" {
		return utils.BadRequestResponse(c, "Device ID required")
	}

	err := h.radioService.UpdateDeviceLocation(c.Context(), req.DeviceID, req.Latitude, req.Longitude)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// DeleteDevice deletes a radio device
// DELETE /radio/device/:id
func (h *RadioHandler) DeleteDevice(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid device ID")
	}

	err = h.radioService.DeleteDevice(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"deleted": true,
	})
}

// ListLocations lists radio locations
// GET /radio/locations
func (h *RadioHandler) ListLocations(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	locations, err := h.radioService.GetLocations(c.Context(), limit)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, locations)
}

// GetLocationByRadioNo gets location by radio number
// GET /radio/location/:radioNo
func (h *RadioHandler) GetLocationByRadioNo(c *fiber.Ctx) error {
	radioNo := c.Params("radioNo")
	if radioNo == "" {
		return utils.BadRequestResponse(c, "Radio number required")
	}

	location, err := h.radioService.GetLocationByRadioNo(c.Context(), radioNo)
	if err != nil {
		if err == service.ErrRadioLocationNotFound {
			return utils.NotFoundResponse(c, "Location not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, location)
}

// CreateLocation creates a new radio location
// POST /radio/location
func (h *RadioHandler) CreateLocation(c *fiber.Ctx) error {
	type CreateRequest struct {
		RadioNo   string  `json:"radioNo"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.RadioNo == "" {
		return utils.BadRequestResponse(c, "Radio number required")
	}

	id, err := h.radioService.CreateLocation(c.Context(), req.RadioNo, req.Latitude, req.Longitude)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"created": true,
	})
}
