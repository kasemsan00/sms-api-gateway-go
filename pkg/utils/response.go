package utils

import "github.com/gofiber/fiber/v2"

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Status: "OK",
		Data:   data,
	})
}

// SuccessResponseWithMessage sends a success response with message
func SuccessResponseWithMessage(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(Response{
		Status:  "OK",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, message string) error {
	return c.JSON(Response{
		Status:  "FAIL",
		Message: message,
	})
}

// ErrorResponseWithData sends an error response with data
func ErrorResponseWithData(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Status: "FAIL",
		Data:   data,
	})
}

// ErrorResponseWithStatus sends an error response with HTTP status code
func ErrorResponseWithStatus(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Status:  "FAIL",
		Message: message,
	})
}

// NotFoundResponse sends a 404 not found response
func NotFoundResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Status:  "FAIL",
		Message: message,
	})
}

// UnauthorizedResponse sends a 401 unauthorized response
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Status:  "FAIL",
		Message: message,
	})
}

// BadRequestResponse sends a 400 bad request response
func BadRequestResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Bad request"
	}
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Status:  "FAIL",
		Message: message,
	})
}

// InternalServerErrorResponse sends a 500 internal server error response
func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Status:  "FAIL",
		Message: message,
	})
}
