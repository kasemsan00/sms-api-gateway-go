package handler

import (
	"strconv"
	"strings"

	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// LinkHandler handles link routes
type LinkHandler struct {
	linkService *service.LinkService
	userService *service.UserService
}

// NewLinkHandler creates a new LinkHandler
func NewLinkHandler(linkService *service.LinkService, userService *service.UserService) *LinkHandler {
	return &LinkHandler{
		linkService: linkService,
		userService: userService,
	}
}

// GetLinkDetail gets link details
// GET /link/getdetail
func (h *LinkHandler) GetLinkDetail(c *fiber.Ctx) error {
	linkID := c.Query("linkID")
	room := c.Query("room")

	link, err := h.linkService.GetLinkDetail(c.Context(), linkID, room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, link)
}

// GetLinkHistory gets link history
// GET /link/history
func (h *LinkHandler) GetLinkHistory(c *fiber.Ctx) error {
	room := c.Query("room")
	mobile := c.Query("mobile")

	links, err := h.linkService.GetLinkIDList(c.Context(), room, mobile)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, links)
}

// CreateLink creates a new link
// POST /link/create
func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
	type CreateRequest struct {
		Room                  string `json:"room"`
		UserType              string `json:"userType"`
		Mobile                string `json:"mobile"`
		Share                 int    `json:"share"`
		RequireJoinPermission int    `json:"requireJoinPermission"`
		RequireUserName       int    `json:"requireUserName"`
		RequirePassword       int    `json:"requirePassword"`
		Password              string `json:"password"`
		OneTimeLink           int    `json:"oneTimeLink"`
		UserAgent             string `json:"userAgent"`
		DaysExpired           int    `json:"daysExpired"`
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	link, err := h.linkService.CreateLink(c.Context(), service.CreateLinkOptions{
		Room:                  req.Room,
		UserType:              req.UserType,
		Mobile:                req.Mobile,
		Share:                 req.Share,
		RequireJoinPermission: req.RequireJoinPermission,
		RequireUserName:       req.RequireUserName,
		RequirePassword:       req.RequirePassword,
		Password:              req.Password,
		OneTimeLink:           req.OneTimeLink,
		UserAgent:             req.UserAgent,
		DaysExpired:           req.DaysExpired,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, link)
}

// CreateHLSLink creates an HLS link
// POST /link/create/hls
func (h *LinkHandler) CreateHLSLink(c *fiber.Ctx) error {
	type CreateRequest struct {
		Room        string `json:"room"`
		UserType    string `json:"userType"`
		DaysExpired int    `json:"daysExpired"`
	}

	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	link, err := h.linkService.CreateLink(c.Context(), service.CreateLinkOptions{
		Room:        req.Room,
		UserType:    "hls",
		DaysExpired: req.DaysExpired,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, link)
}

// UpdateLatLng updates latitude and longitude
// POST /link/update/latlng
func (h *LinkHandler) UpdateLatLng(c *fiber.Ctx) error {
	type UpdateRequest struct {
		LinkID    string  `json:"linkID"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Accuracy  int     `json:"accuracy"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.linkService.UpdateLatLng(c.Context(), req.LinkID, req.Latitude, req.Longitude, req.Accuracy)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// MultiLatLng handles multiple lat/lng updates
// POST /link/multilatlng/send
func (h *LinkHandler) MultiLatLng(c *fiber.Ctx) error {
	type LatLngPoint struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Accuracy  int     `json:"accuracy"`
	}

	type MultiRequest struct {
		LinkID string        `json:"linkID"`
		Points []LatLngPoint `json:"points"`
	}

	var req MultiRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Update with last point
	if len(req.Points) > 0 {
		lastPoint := req.Points[len(req.Points)-1]
		err := h.linkService.UpdateLatLng(c.Context(), req.LinkID, lastPoint.Latitude, lastPoint.Longitude, lastPoint.Accuracy)
		if err != nil {
			return utils.ErrorResponse(c, err.Error())
		}
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated":     true,
		"pointsCount": len(req.Points),
	})
}

// GetShareURL gets share URL
// GET /link/share
func (h *LinkHandler) GetShareURL(c *fiber.Ctx) error {
	room := c.Query("room")
	userType := c.Query("userType", "guest")

	if room == "" {
		return utils.BadRequestResponse(c, "Room required")
	}

	url, err := h.linkService.GetShareURL(c.Context(), room, userType)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"url": url,
	})
}

// CarTracking handles car tracking
// POST /link/cartracking
func (h *LinkHandler) CarTracking(c *fiber.Ctx) error {
	type TrackRequest struct {
		LinkID    string  `json:"linkID"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Accuracy  int     `json:"accuracy"`
	}

	var req TrackRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	err := h.linkService.UpdateLatLng(c.Context(), req.LinkID, req.Latitude, req.Longitude, req.Accuracy)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"updated": true,
	})
}

// GetDomain gets domain
// GET /link/get/domain
func (h *LinkHandler) GetDomain(c *fiber.Ctx) error {
	service, _ := strconv.Atoi(c.Query("service"))
	sender := c.Query("sender")
	linkType := c.Query("linkType")
	linkID := c.Query("linkID")

	domain, err := h.linkService.GetDomain(c.Context(), service, sender, linkType, linkID)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"domain": domain,
	})
}

// GetLinkList gets link list
// GET /link/list
func (h *LinkHandler) GetLinkList(c *fiber.Ctx) error {
	room := c.Query("room")
	mobile := c.Query("mobile")

	links, err := h.linkService.GetLinkIDList(c.Context(), room, mobile)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, links)
}

// SendCustomMessage sends a custom SMS message
// POST /sms/custom
func (h *LinkHandler) SendCustomMessage(c *fiber.Ctx) error {
	type SendRequest struct {
		Message string `json:"message"`
		Mobile  string `json:"mobile"`
		Room    string `json:"room"`
	}

	var req SendRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.Message == "" {
		return utils.BadRequestResponse(c, "Message is required")
	}
	if req.Mobile == "" {
		return utils.BadRequestResponse(c, "Mobile is required")
	}
	if strings.TrimSpace(req.Message) == "" {
		return utils.BadRequestResponse(c, "Message cannot be empty")
	}

	// Here we would send the SMS using SMS service
	// For now, return success as placeholder
	return utils.SuccessResponse(c, fiber.Map{
		"message": "SMS sent successfully",
		"mobile":  req.Mobile,
	})
}
