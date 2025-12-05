package handler

import (
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user routes
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserAlreadyInRoom checks if user is already in room
// GET /user/getuseralreadyinroom
func (h *UserHandler) GetUserAlreadyInRoom(c *fiber.Ctx) error {
	room := c.Query("room")
	identity := c.Query("identity")

	if room == "" || identity == "" {
		return utils.BadRequestResponse(c, "Room and identity required")
	}

	inRoom, err := h.userService.GetUserAlreadyInRoom(c.Context(), room, identity)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"inRoom": inRoom,
	})
}

// GetUserDetail gets user details
// GET /user/getuserdetail
func (h *UserHandler) GetUserDetail(c *fiber.Ctx) error {
	room := c.Query("room")
	identity := c.Query("identity")
	socketID := c.Query("socketId")

	user, err := h.userService.GetUserDetail(c.Context(), room, identity, socketID)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, user)
}

// ListParticipants lists participants in a room
// GET /user/listparticipants
func (h *UserHandler) ListParticipants(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room required")
	}

	participants, err := h.userService.ListParticipants(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, participants)
}

// GenerateUser generates a new user
// POST /user/generate
func (h *UserHandler) GenerateUser(c *fiber.Ctx) error {
	type GenerateRequest struct {
		LinkID    string `json:"linkID"`
		Room      string `json:"room"`
		UserName  string `json:"userName"`
		UserType  string `json:"userType"`
		UserAgent string `json:"userAgent"`
		ServiceID int    `json:"serviceId"`
	}

	var req GenerateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	result, err := h.userService.GenerateUser(c.Context(), service.GenerateUserOptions{
		LinkID:    req.LinkID,
		Room:      req.Room,
		UserName:  req.UserName,
		UserType:  req.UserType,
		UserAgent: req.UserAgent,
		ServiceID: req.ServiceID,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, result)
}

// JoinGenerate generates a user for joining
// POST /user/joingenerate
func (h *UserHandler) JoinGenerate(c *fiber.Ctx) error {
	type JoinRequest struct {
		Room      string `json:"room"`
		UserName  string `json:"userName"`
		UserType  string `json:"userType"`
		UserAgent string `json:"userAgent"`
	}

	var req JoinRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	result, err := h.userService.GenerateUser(c.Context(), service.GenerateUserOptions{
		Room:      req.Room,
		UserName:  req.UserName,
		UserType:  req.UserType,
		UserAgent: req.UserAgent,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, result)
}

// GenerateChatUser generates a chat user
// POST /user/generateChatUser
func (h *UserHandler) GenerateChatUser(c *fiber.Ctx) error {
	type ChatUserRequest struct {
		Room     string `json:"room"`
		UserName string `json:"userName"`
	}

	var req ChatUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	result, err := h.userService.GenerateUser(c.Context(), service.GenerateUserOptions{
		Room:     req.Room,
		UserName: req.UserName,
		UserType: "chat",
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, result)
}

// UpdateParticipants updates participants
// POST /user/updateparticipants
func (h *UserHandler) UpdateParticipants(c *fiber.Ctx) error {
	type UpdateRequest struct {
		Room     string `json:"room"`
		Identity string `json:"identity"`
		Status   string `json:"status"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.userService.UpdateUserStatus(c.Context(), req.Room, req.Identity, req.Status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// MutePublishedTrack mutes a published track
// POST /user/mutepublishedtrack
func (h *UserHandler) MutePublishedTrack(c *fiber.Ctx) error {
	type MuteRequest struct {
		Room     string `json:"room"`
		Identity string `json:"identity"`
		TrackSid string `json:"trackSid"`
		Muted    bool   `json:"muted"`
	}

	var req MuteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.userService.MutePublishedTrack(c.Context(), req.Room, req.Identity, req.TrackSid, req.Muted)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"muted": req.Muted,
	})
}

// RemoveParticipant removes a participant
// POST /user/removeParticipant
func (h *UserHandler) RemoveParticipant(c *fiber.Ctx) error {
	type RemoveRequest struct {
		Room     string `json:"room"`
		Identity string `json:"identity"`
	}

	var req RemoveRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.userService.RemoveParticipant(c.Context(), req.Room, req.Identity)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"removed": true,
	})
}

// GetUserLog gets user log
// GET /user/log
func (h *UserHandler) GetUserLog(c *fiber.Ctx) error {
	room := c.Query("room")
	status := c.Query("status", "connect")

	users, err := h.userService.ListUsersInRoom(c.Context(), room, status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, users)
}

// HandleTrack handles track updates
// PUT /user/handle/track
func (h *UserHandler) HandleTrack(c *fiber.Ctx) error {
	type TrackRequest struct {
		Identity   string `json:"identity"`
		Camera     *bool  `json:"camera"`
		Microphone *bool  `json:"microphone"`
	}

	var req TrackRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.Camera != nil {
		err := h.userService.UpdateUserCamera(c.Context(), req.Identity, *req.Camera)
		if err != nil {
			return utils.ErrorResponse(c, err.Error())
		}
	}

	if req.Microphone != nil {
		err := h.userService.UpdateUserMicrophone(c.Context(), req.Identity, *req.Microphone)
		if err != nil {
			return utils.ErrorResponse(c, err.Error())
		}
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}
