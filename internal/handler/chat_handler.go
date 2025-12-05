package handler

import (
	"strconv"

	"api-gateway-go/internal/repository"
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ChatHandler handles chat routes
type ChatHandler struct {
	chatService *service.ChatService
}

// NewChatHandler creates a new ChatHandler
func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// GetHistory gets chat history
// GET /chat/history
func (h *ChatHandler) GetHistory(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	limitStr := c.Query("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	history, err := h.chatService.GetChatHistory(c.Context(), room, limit)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, history)
}

// GetNotification gets chat notifications
// GET /chat/notification
func (h *ChatHandler) GetNotification(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	notifications, err := h.chatService.GetChatNotification(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, notifications)
}

// AddMessage adds a new chat message
// POST /chat/message
func (h *ChatHandler) AddMessage(c *fiber.Ctx) error {
	type AddMessageRequest struct {
		Room             string `json:"room"`
		Identity         string `json:"identity"`
		ChatIdentity     string `json:"chatIdentity"`
		UserName         string `json:"userName"`
		UserType         string `json:"userType"`
		Text             string `json:"text"`
		Files            string `json:"files"`
		Color            string `json:"color"`
		ReplyToMessageID int    `json:"replyToMessageId"`
		ReplyToUserName  string `json:"replyToUserName"`
		ReplyToText      string `json:"replyToText"`
	}

	var req AddMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.Room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	id, err := h.chatService.AddChatMessage(c.Context(), repository.SaveMessageParams{
		Room:             req.Room,
		Identity:         req.Identity,
		ChatIdentity:     req.ChatIdentity,
		UserName:         req.UserName,
		UserType:         req.UserType,
		Text:             req.Text,
		Files:            req.Files,
		Color:            req.Color,
		ReplyToMessageID: req.ReplyToMessageID,
		ReplyToUserName:  req.ReplyToUserName,
		ReplyToText:      req.ReplyToText,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"created": true,
	})
}

// GetMessageCount gets message count for a room
// GET /chat/count
func (h *ChatHandler) GetMessageCount(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	count, err := h.chatService.GetMessageCount(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"room":  room,
		"count": count,
	})
}

// DeleteMessages deletes all messages for a room
// DELETE /chat/messages
func (h *ChatHandler) DeleteMessages(c *fiber.Ctx) error {
	room := c.Query("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room name required")
	}

	err := h.chatService.DeleteMessagesByRoom(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"deleted": true,
	})
}
