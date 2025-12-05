package handler

import (
	"strconv"

	"api-gateway-go/internal/repository"
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// RoomHandler handles room routes
type RoomHandler struct {
	roomService *service.RoomService
	authService *service.AuthService
}

// NewRoomHandler creates a new RoomHandler
func NewRoomHandler(roomService *service.RoomService, authService *service.AuthService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
		authService: authService,
	}
}

// GetRoomDetail gets room details
// GET /room/detail
func (h *RoomHandler) GetRoomDetail(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	roomDetail, err := h.roomService.GetRoomDetail(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, roomDetail)
}

// ListRooms lists all rooms by status
// GET /room/listrooms
func (h *RoomHandler) ListRooms(c *fiber.Ctx) error {
	status := c.Query("status", "open")

	rooms, err := h.roomService.GetRoomsByStatus(c.Context(), status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, rooms)
}

// CheckExpired checks if a room has expired
// GET /room/checkexpired
func (h *RoomHandler) CheckExpired(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	expired, err := h.roomService.CheckRoomExpired(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"room":    room,
		"expired": expired,
	})
}

// VerifyToken verifies room token
// GET /room/verifytoken
func (h *RoomHandler) VerifyToken(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		token = c.Query("token")
	}

	if token == "" {
		return utils.ErrorResponseWithStatus(c, fiber.StatusUnauthorized, "Token required")
	}

	claims, err := h.authService.VerifyToken(c.Context(), token)
	if err != nil {
		return utils.ErrorResponseWithStatus(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, claims)
}

// UpdateUser updates room user
// POST /room/updateuser
func (h *RoomHandler) UpdateUser(c *fiber.Ctx) error {
	type UpdateRequest struct {
		Room     string `json:"room"`
		Identity string `json:"identity"`
		Status   string `json:"status"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.roomService.UpdateRoomStatus(c.Context(), req.Room, req.Status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// DeleteRoom deletes a room
// POST /room/deleteroom
func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	type DeleteRequest struct {
		Room string `json:"room"`
	}

	var req DeleteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.roomService.DeleteRoom(c.Context(), req.Room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"deleted": true,
	})
}

// UpdateType updates room type
// PUT /room/updatetype
func (h *RoomHandler) UpdateType(c *fiber.Ctx) error {
	type UpdateRequest struct {
		Room         string `json:"room"`
		RoomType     string `json:"roomType"`
		AutoRecord   int    `json:"autoRecord"`
		ChatEnabled  int    `json:"chatEnabled"`
		WebSocketURL string `json:"webSocketURL"`
		UserAgent    string `json:"userAgent"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.roomService.UpdateRoomType(c.Context(), req.Room, repository.UpdateRoomTypeParams{
		RoomType:     req.RoomType,
		AutoRecord:   req.AutoRecord,
		ChatEnabled:  req.ChatEnabled,
		WebSocketURL: req.WebSocketURL,
		UserAgent:    req.UserAgent,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// UpdateStatus updates room status
// PUT /room/updatestatus
func (h *RoomHandler) UpdateStatus(c *fiber.Ctx) error {
	type UpdateRequest struct {
		Room   string `json:"room"`
		Status string `json:"status"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.roomService.UpdateRoomStatus(c.Context(), req.Room, req.Status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// CloseRoom closes a room
// PUT /room/close
func (h *RoomHandler) CloseRoom(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		type CloseRequest struct {
			Room string `json:"room"`
		}
		var req CloseRequest
		if err := c.BodyParser(&req); err == nil {
			room = req.Room
		}
	}

	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	err := h.roomService.CloseRoom(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"closed": true,
	})
}

// CreateRoom creates a new room
// POST /room/create
func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	type CreateRequest struct {
		Service               int    `json:"service"`
		RoomType              string `json:"roomType"`
		AutoRecord            int    `json:"autoRecord"`
		RecordType            string `json:"recordType"`
		EncodingOptionsPreset string `json:"encodingOptionsPreset"`
		ChatEnabled           int    `json:"chatEnabled"`
		WebSocketURL          string `json:"webSocketURL"`
		UserAgent             string `json:"userAgent"`
		DaysExpired           int    `json:"daysExpired"`
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	room, err := h.roomService.CreateRoom(c.Context(), service.CreateRoomOptions{
		Service:               req.Service,
		RoomType:              req.RoomType,
		AutoRecord:            req.AutoRecord,
		RecordType:            req.RecordType,
		EncodingOptionsPreset: req.EncodingOptionsPreset,
		ChatEnabled:           req.ChatEnabled,
		WebSocketURL:          req.WebSocketURL,
		UserAgent:             req.UserAgent,
		DaysExpired:           req.DaysExpired,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, room)
}

// GetRoomPicture gets room picture
// GET /room/picture
func (h *RoomHandler) GetRoomPicture(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	// This is a placeholder - actual implementation would fetch room picture
	return utils.SuccessResponse(c, fiber.Map{
		"room":    room,
		"picture": nil,
	})
}

// UpdateRecordStatus updates room record status
// PUT /room/recordstatus
func (h *RoomHandler) UpdateRecordStatus(c *fiber.Ctx) error {
	room := c.Query("room")
	statusStr := c.Query("status", "0")
	status, _ := strconv.Atoi(statusStr)

	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	err := h.roomService.UpdateRecordStatus(c.Context(), room, status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}
