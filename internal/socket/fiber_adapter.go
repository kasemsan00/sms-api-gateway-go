package socket

import (
	"net/http"

	"api-gateway-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// SetupSocketIO sets up Socket.IO routes with Fiber
func SetupSocketIO(app *fiber.App, hub *Hub) {
	// Get the Socket.IO server
	server := hub.Server()

	// Create HTTP handlers for Socket.IO
	socketHandler := func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	}

	// Use adaptor to convert net/http handler to Fiber handler
	app.Use("/socket.io/", adaptor.HTTPHandler(http.HandlerFunc(socketHandler)))

	logger.Info("Socket.IO routes configured")
}

// SocketNamespaceMiddleware is a middleware to initialize room namespaces dynamically
func SocketNamespaceMiddleware(hub *Hub) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if this is a room-related request that needs namespace initialization
		room := c.Query("room")
		if room != "" {
			// Initialize the room namespace if not already done
			hub.InitRoomNamespace(room)
		}
		return c.Next()
	}
}
