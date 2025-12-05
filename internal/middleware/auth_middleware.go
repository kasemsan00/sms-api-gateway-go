package middleware

import (
	"strings"

	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		token := ""

		if authHeader != "" {
			// Check for Bearer token
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				token = authHeader
			}
		}

		// Try query parameter
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			return utils.UnauthorizedResponse(c, "Invalid token")
		}

		// Verify token
		claims, err := authService.VerifyToken(c.Context(), token)
		if err != nil {
			return utils.ErrorResponseWithStatus(c, fiber.StatusUnauthorized, err.Error())
		}

		// Store claims in context
		c.Locals("user", claims)
		c.Locals("room", claims.Room)
		c.Locals("identity", claims.Identity)
		c.Locals("userName", claims.UserName)
		c.Locals("userType", claims.UserType)

		return c.Next()
	}
}

// OptionalAuthMiddleware creates an optional authentication middleware
// It doesn't fail if no token is provided, but parses it if present
func OptionalAuthMiddleware(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		token := ""

		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				token = authHeader
			}
		}

		if token == "" {
			token = c.Query("token")
		}

		if token != "" {
			claims, err := authService.VerifyToken(c.Context(), token)
			if err == nil {
				c.Locals("user", claims)
				c.Locals("room", claims.Room)
				c.Locals("identity", claims.Identity)
				c.Locals("userName", claims.UserName)
				c.Locals("userType", claims.UserType)
			}
		}

		return c.Next()
	}
}

// GetUserFromContext gets user claims from context
func GetUserFromContext(c *fiber.Ctx) *service.Claims {
	user := c.Locals("user")
	if user == nil {
		return nil
	}
	claims, ok := user.(*service.Claims)
	if !ok {
		return nil
	}
	return claims
}

// GetRoomFromContext gets room from context
func GetRoomFromContext(c *fiber.Ctx) string {
	room := c.Locals("room")
	if room == nil {
		return ""
	}
	return room.(string)
}

// GetIdentityFromContext gets identity from context
func GetIdentityFromContext(c *fiber.Ctx) string {
	identity := c.Locals("identity")
	if identity == nil {
		return ""
	}
	return identity.(string)
}
