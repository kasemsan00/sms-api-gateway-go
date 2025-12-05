package handler

import (
	"strconv"

	"api-gateway-go/internal/repository"
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// NotificationHandler handles notification routes
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// ListNotifications lists all notifications
// GET /notification/list
func (h *NotificationHandler) ListNotifications(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "100")
	offsetStr := c.Query("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	notifications, err := h.notificationService.GetAllNotifications(c.Context(), limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, notifications)
}

// GetByUserName gets notifications by username
// GET /notification/user
func (h *NotificationHandler) GetByUserName(c *fiber.Ctx) error {
	userName := c.Query("userName")
	if userName == "" {
		return utils.BadRequestResponse(c, "User name required")
	}

	notifications, err := h.notificationService.GetNotificationsByUserName(c.Context(), userName)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, notifications)
}

// GetUnread gets unread notifications
// GET /notification/unread
func (h *NotificationHandler) GetUnread(c *fiber.Ctx) error {
	notifications, err := h.notificationService.GetUnreadNotifications(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, notifications)
}

// GetUnreadCount gets count of unread notifications
// GET /notification/unreadcount
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	count, err := h.notificationService.GetUnreadCount(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"count": count,
	})
}

// GetByID gets a notification by ID
// GET /notification/:id
func (h *NotificationHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid notification ID")
	}

	notification, err := h.notificationService.GetNotificationByID(c.Context(), id)
	if err != nil {
		if err == service.ErrNotificationNotFound {
			return utils.NotFoundResponse(c, "Notification not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, notification)
}

// Create creates a new notification
// POST /notification/create
func (h *NotificationHandler) Create(c *fiber.Ctx) error {
	type CreateRequest struct {
		UserName         string `json:"userName"`
		Mobile           string `json:"mobile"`
		Message          string `json:"message"`
		CaseID           int    `json:"caseId"`
		NotificationType string `json:"notificationType"`
		RelatedURL       string `json:"relatedUrl"`
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	id, err := h.notificationService.CreateNotification(c.Context(), repository.CreateNotificationParams{
		UserName:         req.UserName,
		Mobile:           req.Mobile,
		Message:          req.Message,
		CaseID:           req.CaseID,
		NotificationType: req.NotificationType,
		RelatedURL:       req.RelatedURL,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"created": true,
	})
}

// MarkAsRead marks a notification as read
// PUT /notification/read/:id
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid notification ID")
	}

	err = h.notificationService.UpdateReadStatus(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// MarkAllAsRead marks all notifications as read
// PUT /notification/readall
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	err := h.notificationService.MarkAllAsRead(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// Delete deletes a notification
// DELETE /notification/:id
func (h *NotificationHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid notification ID")
	}

	err = h.notificationService.DeleteNotification(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"deleted": true,
	})
}
