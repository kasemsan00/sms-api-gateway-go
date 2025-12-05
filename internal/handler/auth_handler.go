package handler

import (
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication routes
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// CreateToken creates a new JWT token
// GET /auth/create
func (h *AuthHandler) CreateToken(c *fiber.Ctx) error {
	userName := c.Query("userName")
	room := c.Query("room")
	identity := c.Query("identity")
	userType := c.Query("userType")

	token, err := h.authService.CreateToken(c.Context(), userName, room, identity, userType, 0)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"token": token,
	})
}

// VerifyToken verifies a JWT token
// GET /auth/verify
func (h *AuthHandler) VerifyToken(c *fiber.Ctx) error {
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

// VerifyUser verifies a user by username and password
// POST /auth/verifyuser
func (h *AuthHandler) VerifyUser(c *fiber.Ctx) error {
	type VerifyRequest struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}

	var req VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// This is a placeholder - actual implementation would check against a user database
	// For now, we just return success
	return utils.SuccessResponse(c, fiber.Map{
		"verified": true,
		"userName": req.UserName,
	})
}
