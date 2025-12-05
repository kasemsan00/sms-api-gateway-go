package handler

import (
	"strconv"

	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// CarHandler handles car tracking routes
type CarHandler struct {
	carService *service.CarService
}

// NewCarHandler creates a new CarHandler
func NewCarHandler(carService *service.CarService) *CarHandler {
	return &CarHandler{
		carService: carService,
	}
}

// CreateTask creates a new car tracking task
// POST /car/task
func (h *CarHandler) CreateTask(c *fiber.Ctx) error {
	type CreateRequest struct {
		UID      string `json:"uid"`
		Status   string `json:"status"`
		Mobile   string `json:"mobile"`
		UserName string `json:"userName"`
		Room     string `json:"room"`
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	id, err := h.carService.CreateTask(c.Context(), service.CreateCarTaskOptions{
		UID:      req.UID,
		Status:   req.Status,
		Mobile:   req.Mobile,
		UserName: req.UserName,
		Room:     req.Room,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"created": true,
	})
}

// GetTaskDetail gets task detail by ID
// GET /car/task/:id
func (h *CarHandler) GetTaskDetail(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid task ID")
	}

	task, err := h.carService.GetTaskDetail(c.Context(), id)
	if err != nil {
		if err == service.ErrCarTaskNotFound {
			return utils.NotFoundResponse(c, "Task not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, task)
}

// GetTaskByUID gets task by UID
// GET /car/uid/:uid
func (h *CarHandler) GetTaskByUID(c *fiber.Ctx) error {
	uid := c.Params("uid")
	if uid == "" {
		return utils.BadRequestResponse(c, "UID required")
	}

	task, err := h.carService.GetTaskByUID(c.Context(), uid)
	if err != nil {
		if err == service.ErrCarTaskNotFound {
			return utils.NotFoundResponse(c, "Task not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, task)
}

// GetTaskByRoom gets task by room
// GET /car/room/:room
func (h *CarHandler) GetTaskByRoom(c *fiber.Ctx) error {
	room := c.Params("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room required")
	}

	task, err := h.carService.GetTaskByRoom(c.Context(), room)
	if err != nil {
		if err == service.ErrCarTaskNotFound {
			return utils.NotFoundResponse(c, "Task not found")
		}
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, task)
}

// ListTasks lists car tasks
// GET /car/list
func (h *CarHandler) ListTasks(c *fiber.Ctx) error {
	status := c.Query("status", "")
	limitStr := c.Query("limit", "100")
	offsetStr := c.Query("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	tasks, err := h.carService.GetTaskList(c.Context(), status, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, tasks)
}

// UpdateTask updates a car task
// PUT /car/task/:id
func (h *CarHandler) UpdateTask(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid task ID")
	}

	type UpdateRequest struct {
		Status string `json:"status"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err = h.carService.UpdateTask(c.Context(), id, req.Status)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// UpdatePosition updates car position
// POST /car/position
func (h *CarHandler) UpdatePosition(c *fiber.Ctx) error {
	type UpdateRequest struct {
		UserName         string  `json:"userName"`
		Room             string  `json:"room"`
		Latitude         float64 `json:"latitude"`
		Longitude        float64 `json:"longitude"`
		Accuracy         float64 `json:"accuracy"`
		Altitude         float64 `json:"altitude"`
		AltitudeAccuracy float64 `json:"altitudeAccuracy"`
		Speed            int     `json:"speed"`
		Heading          int     `json:"heading"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.carService.UpdatePosition(c.Context(),
		req.UserName,
		req.Room,
		req.Latitude,
		req.Longitude,
		req.Accuracy,
		req.Altitude,
		req.AltitudeAccuracy,
		req.Speed,
		req.Heading,
	)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// GetCarPosition gets car position by room
// GET /car/position/:room
func (h *CarHandler) GetCarPosition(c *fiber.Ctx) error {
	room := c.Params("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room required")
	}

	position, err := h.carService.GetCarPosition(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, position)
}

// GetUserLatLng gets user lat/lng by room
// GET /car/latlng/:room
func (h *CarHandler) GetUserLatLng(c *fiber.Ctx) error {
	room := c.Params("room")
	if room == "" {
		return utils.BadRequestResponse(c, "Room required")
	}

	position, err := h.carService.GetUserLatLng(c.Context(), room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, position)
}

// DeleteTask deletes a car task
// DELETE /car/task/:id
func (h *CarHandler) DeleteTask(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid task ID")
	}

	err = h.carService.DeleteTask(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"deleted": true,
	})
}
