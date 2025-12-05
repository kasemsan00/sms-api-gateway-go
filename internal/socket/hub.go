package socket

import (
	"context"
	"net/http"
	"sync"
	"time"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/repository"
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/logger"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// SocketConfig holds socket configuration
type SocketConfig struct {
	CleanupInterval     time.Duration
	JoinRoomRepeatDelay time.Duration
	DisconnectTimeout   time.Duration
	MaxRetries          int
	MaxMessageLength    int
	MaxConnectionsPerNS int
}

// DefaultSocketConfig returns default socket configuration
func DefaultSocketConfig() *SocketConfig {
	return &SocketConfig{
		CleanupInterval:     1 * time.Hour,
		JoinRoomRepeatDelay: 5 * time.Second,
		DisconnectTimeout:   1 * time.Second,
		MaxRetries:          3,
		MaxMessageLength:    1000,
		MaxConnectionsPerNS: 1000,
	}
}

// Hub manages all socket connections and namespaces
type Hub struct {
	server   *socketio.Server
	config   *SocketConfig
	cfg      *config.Config
	redisMgr *config.RedisManager
	eventMgr *CrossInstanceEventManager
	stateMgr *SocketStateManager

	// Services
	roomService *service.RoomService
	userService *service.UserService
	chatService *service.ChatService
	carService  *service.CarService
	linkService *service.LinkService

	// Repositories
	roomRepo *repository.RoomRepository
	userRepo *repository.UserRepository
	chatRepo *repository.ChatRepository

	// Registered namespaces
	namespaces   map[string]bool
	namespacesMu sync.RWMutex

	// Cleanup
	cleanupTicker *time.Ticker
	done          chan struct{}
}

// NewHub creates a new Socket.IO hub
func NewHub(
	cfg *config.Config,
	redisMgr *config.RedisManager,
	roomService *service.RoomService,
	userService *service.UserService,
	chatService *service.ChatService,
	carService *service.CarService,
	linkService *service.LinkService,
	roomRepo *repository.RoomRepository,
	userRepo *repository.UserRepository,
	chatRepo *repository.ChatRepository,
) (*Hub, error) {
	// Create Socket.IO server with custom transport options
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOrigin,
			},
			&websocket.Transport{
				CheckOrigin: allowOrigin,
			},
		},
	})

	hub := &Hub{
		server:      server,
		config:      DefaultSocketConfig(),
		cfg:         cfg,
		redisMgr:    redisMgr,
		eventMgr:    NewCrossInstanceEventManager(redisMgr),
		stateMgr:    NewSocketStateManager(redisMgr),
		roomService: roomService,
		userService: userService,
		chatService: chatService,
		carService:  carService,
		linkService: linkService,
		roomRepo:    roomRepo,
		userRepo:    userRepo,
		chatRepo:    chatRepo,
		namespaces:  make(map[string]bool),
		done:        make(chan struct{}),
	}

	// Initialize cross-instance event manager
	if err := hub.eventMgr.Initialize(); err != nil {
		logger.Error("Failed to initialize event manager: %v", err)
	}

	// Setup default namespace handlers
	hub.setupDefaultNamespace()

	return hub, nil
}

// allowOrigin allows all origins (configure for production)
func allowOrigin(r *http.Request) bool {
	return true
}

// setupDefaultNamespace sets up the default namespace handlers
func (h *Hub) setupDefaultNamespace() {
	h.server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logger.Debug("Socket connected: %s", s.ID())
		s.Emit("log", map[string]string{"message": "connection success"})
		return nil
	})

	h.server.OnError("/", func(s socketio.Conn, e error) {
		logger.Error("Socket error: %v", e)
	})

	h.server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logger.Debug("Socket disconnected: %s, reason: %s", s.ID(), reason)
	})
}

// InitRoomNamespace initializes a namespace for a specific room
func (h *Hub) InitRoomNamespace(room string) {
	h.namespacesMu.Lock()
	defer h.namespacesMu.Unlock()

	namespace := "/" + room
	if h.namespaces[namespace] {
		return // Already initialized
	}

	// Setup connection handler
	h.server.OnConnect(namespace, func(s socketio.Conn) error {
		return h.handleRoomConnection(s, room)
	})

	// Setup event handlers
	h.setupRoomEventHandlers(namespace, room)

	// Setup disconnect handler
	h.server.OnDisconnect(namespace, func(s socketio.Conn, reason string) {
		h.handleRoomDisconnect(s, room, reason)
	})

	// Setup cross-instance event handler
	h.setupCrossInstanceEvents(namespace, room)

	h.namespaces[namespace] = true
	logger.Info("Initialized namespace: %s", namespace)
}

// handleRoomConnection handles a new connection to a room namespace
func (h *Hub) handleRoomConnection(s socketio.Conn, room string) error {
	// Get query parameters from URL
	urlObj := s.URL()
	identity := urlObj.Query().Get("identity")
	linkID := urlObj.Query().Get("linkID")

	logger.Debug("Socket %s connecting to room %s with identity %s", s.ID(), room, identity)

	// Store socket state
	ctx := s.Context()
	if ctx == nil {
		ctx = &SocketState{
			SocketID: s.ID(),
			Room:     room,
			Identity: identity,
		}
		s.SetContext(ctx)
	}

	// Join the room
	s.Join(room)

	// Emit connection success
	s.Emit("log", map[string]string{"message": "connection success"})

	// Get and emit user details if identity exists
	if identity != "" {
		// In a real implementation, you would fetch user details from the database
		// and emit user-connection event to all clients in the room
		h.server.BroadcastToRoom(room, room, "user-connection", map[string]interface{}{
			"identity": identity,
			"status":   "connect",
		})
	}

	// Handle link connection time
	if linkID != "" {
		// Update link connect time via service
		// This would be implemented with proper service call
	}

	return nil
}

// setupRoomEventHandlers sets up all event handlers for a room namespace
func (h *Hub) setupRoomEventHandlers(namespace, room string) {
	// Chat handlers
	h.server.OnEvent(namespace, "joinChat", func(s socketio.Conn, roomName string) {
		s.Join(roomName)
		logger.Debug("User joined chat room: %s", roomName)
	})

	h.server.OnEvent(namespace, "chat-message", func(s socketio.Conn, data map[string]interface{}) {
		h.handleChatMessage(s, room, data)
	})

	h.server.OnEvent(namespace, "get-chat-history", func(s socketio.Conn, data map[string]interface{}) {
		h.handleGetChatHistory(s, room, data)
	})

	// Position handlers
	h.server.OnEvent(namespace, "user-position", func(s socketio.Conn, data map[string]interface{}) {
		h.handleUserPosition(s, room, data)
	})

	h.server.OnEvent(namespace, "position", func(s socketio.Conn, data map[string]interface{}) {
		h.handlePosition(s, room, data)
	})

	// User handlers
	h.server.OnEvent(namespace, "get-user-detail", func(s socketio.Conn, data map[string]interface{}) {
		h.handleGetUserDetail(s, room, data)
	})

	h.server.OnEvent(namespace, "get-user-list", func(s socketio.Conn, data map[string]interface{}) {
		h.handleGetUserList(s, room)
	})

	h.server.OnEvent(namespace, "update-username", func(s socketio.Conn, data map[string]interface{}) {
		h.handleUpdateUsername(s, room, data)
	})

	// Conference handlers
	h.server.OnEvent(namespace, "user-conference", func(s socketio.Conn, data map[string]interface{}) {
		h.handleUserConference(s, room, data)
	})

	h.server.OnEvent(namespace, "auth-join-conference", func(s socketio.Conn, data map[string]interface{}) {
		h.handleAuthJoinConference(s, room, data)
	})

	h.server.OnEvent(namespace, "auth-join-conference-answer", func(s socketio.Conn, data map[string]interface{}) {
		h.handleAuthJoinConferenceAnswer(s, room, data)
	})

	// Room handlers
	h.server.OnEvent(namespace, "room-record", func(s socketio.Conn) {
		h.handleRoomRecord(s, room)
	})

	h.server.OnEvent(namespace, "case-data", func(s socketio.Conn, data map[string]interface{}) {
		h.server.BroadcastToRoom(namespace, room, "case-data", data)
	})
}

// handleRoomDisconnect handles a disconnect from a room namespace
func (h *Hub) handleRoomDisconnect(s socketio.Conn, room, reason string) {
	logger.Debug("Socket %s disconnected from room %s: %s", s.ID(), room, reason)

	// Get user context
	ctx := s.Context()
	if ctx == nil {
		return
	}

	state, ok := ctx.(*SocketState)
	if !ok {
		return
	}

	// Broadcast user disconnect
	h.server.BroadcastToRoom("/"+room, room, "user-disconnect", map[string]interface{}{
		"status":   "disconnect",
		"identity": state.Identity,
		"userName": state.UserName,
	})

	// Publish disconnect event for cross-instance
	h.eventMgr.PublishEvent(context.Background(), room, EventUserDisconnect, map[string]interface{}{
		"identity": state.Identity,
		"userName": state.UserName,
	})
}

// setupCrossInstanceEvents sets up cross-instance event handling for a room
func (h *Hub) setupCrossInstanceEvents(namespace, room string) {
	handler := func(event EventType, data interface{}) {
		switch event {
		case EventUserDisconnect:
			h.server.BroadcastToRoom(namespace, room, "user-disconnect", data)
		case EventUserConnect:
			h.server.BroadcastToRoom(namespace, room, "user-connection", data)
		case EventCarPosition:
			h.server.BroadcastToRoom(namespace, room, "agentCar", data)
		case EventChatMessage:
			h.server.BroadcastToRoom(namespace, room, "chat-message", data)
		case EventRoomRecord:
			h.server.BroadcastToRoom(namespace, room, "room-record", data)
		}
	}

	h.eventMgr.RegisterHandler(room, handler)
}

// Chat event handlers
func (h *Hub) handleChatMessage(s socketio.Conn, room string, data map[string]interface{}) {
	// Add timestamp
	data["dtmcreated"] = time.Now().Format("2006-01-02 15:04:05")

	// Broadcast to room
	h.server.BroadcastToRoom("/"+room, room, "chat-message", data)

	// Publish for cross-instance
	h.eventMgr.PublishEvent(context.Background(), room, EventChatMessage, data)
}

func (h *Hub) handleGetChatHistory(s socketio.Conn, room string, data map[string]interface{}) {
	// Get chat history from service
	// This would use h.chatService.GetChatHistory
	s.Emit("chat-history", []interface{}{})
}

// Position event handlers
func (h *Hub) handleUserPosition(s socketio.Conn, room string, data map[string]interface{}) {
	h.server.BroadcastToRoom("/"+room, room, "user-position", data)
}

func (h *Hub) handlePosition(s socketio.Conn, room string, data map[string]interface{}) {
	h.server.BroadcastToRoom("/"+room, room, "position", data)

	// Publish for cross-instance
	h.eventMgr.PublishEvent(nil, room, EventCarPosition, data)
}

// User event handlers
func (h *Hub) handleGetUserDetail(s socketio.Conn, room string, data map[string]interface{}) {
	// Get user details and emit
	s.Emit("user-detail", map[string]interface{}{})
}

func (h *Hub) handleGetUserList(s socketio.Conn, room string) {
	// Get user list and emit
	s.Emit("user-list", []interface{}{})
}

func (h *Hub) handleUpdateUsername(s socketio.Conn, room string, data map[string]interface{}) {
	// Update username and broadcast
	h.server.BroadcastToRoom("/"+room, room, "update-username", data)
}

// Conference event handlers
func (h *Hub) handleUserConference(s socketio.Conn, room string, data map[string]interface{}) {
	h.server.BroadcastToRoom("/"+room, room, "user-conference", data)
}

func (h *Hub) handleAuthJoinConference(s socketio.Conn, room string, data map[string]interface{}) {
	// Handle auth join conference request
	h.server.BroadcastToRoom("/"+room, room, "auth-join-conference", data)
}

func (h *Hub) handleAuthJoinConferenceAnswer(s socketio.Conn, room string, data map[string]interface{}) {
	h.server.BroadcastToRoom("/"+room, room, "auth-join-conference-answer", data)
}

// Room event handlers
func (h *Hub) handleRoomRecord(s socketio.Conn, room string) {
	// Handle room recording toggle
	response := map[string]interface{}{
		"status": "recording",
		"room":   room,
	}

	h.server.BroadcastToRoom("/"+room, room, "room-record", response)
	h.eventMgr.PublishEvent(context.Background(), room, EventRoomRecord, response)
}

// InitMobileNamespace initializes the mobile namespace
func (h *Hub) InitMobileNamespace() {
	namespace := "/mobile"

	h.server.OnConnect(namespace, func(s socketio.Conn) error {
		urlObj := s.URL()
		taskID := urlObj.Query().Get("id")
		if taskID == "" {
			return nil // Will be disconnected
		}

		s.SetContext(map[string]string{"taskId": taskID})

		s.Emit("connection", map[string]interface{}{
			"message":   "connection success",
			"case":      taskID,
			"status":    "open || arrive || cancel || complete",
			"latitude":  13.734,
			"longitude": 100.567,
			"accuracy":  100,
		})

		return nil
	})

	h.server.OnEvent(namespace, "location", func(s socketio.Conn, data map[string]interface{}) {
		// Handle location update from mobile
		logger.Debug("Mobile location update: %v", data)
	})

	h.server.OnDisconnect(namespace, func(s socketio.Conn, reason string) {
		logger.Debug("Mobile disconnected: %s", reason)
	})

	h.namespaces[namespace] = true
	logger.Info("Initialized mobile namespace")
}

// InitQueueNamespace initializes the queue namespace
func (h *Hub) InitQueueNamespace() {
	namespace := "/queue"

	h.server.OnConnect(namespace, func(s socketio.Conn) error {
		s.Join("queue")
		s.Emit("log", map[string]string{"message": "connected to queue"})
		return nil
	})

	h.server.OnEvent(namespace, "queue", func(s socketio.Conn, data map[string]interface{}) {
		h.server.BroadcastToRoom(namespace, "queue", "queue", data)
	})

	h.server.OnDisconnect(namespace, func(s socketio.Conn, reason string) {
		logger.Debug("Queue client disconnected: %s", reason)
	})

	h.namespaces[namespace] = true
	logger.Info("Initialized queue namespace")
}

// InitNewQueueNamespace initializes the newqueue namespace
func (h *Hub) InitNewQueueNamespace() {
	namespace := "/newqueue"

	h.server.OnConnect(namespace, func(s socketio.Conn) error {
		s.Join("newqueue")
		s.Emit("log", map[string]string{"message": "connected to newqueue"})
		return nil
	})

	h.server.OnEvent(namespace, "queue", func(s socketio.Conn, data map[string]interface{}) {
		h.server.BroadcastToRoom(namespace, "newqueue", "queue", data)
	})

	h.server.OnEvent(namespace, "case-data", func(s socketio.Conn, data map[string]interface{}) {
		h.server.BroadcastToRoom(namespace, "newqueue", "case-data", data)
	})

	h.server.OnDisconnect(namespace, func(s socketio.Conn, reason string) {
		logger.Debug("New queue client disconnected: %s", reason)
	})

	h.namespaces[namespace] = true
	logger.Info("Initialized newqueue namespace")
}

// InitNotificationNamespace initializes the notification namespace
func (h *Hub) InitNotificationNamespace() {
	namespace := "/notification"

	h.server.OnConnect(namespace, func(s socketio.Conn) error {
		s.Join("notification")
		s.Emit("log", map[string]string{"message": "connected to notifications"})
		return nil
	})

	h.server.OnEvent(namespace, "all", func(s socketio.Conn, data map[string]interface{}) {
		h.server.BroadcastToRoom(namespace, "notification", "all", data)
	})

	h.server.OnEvent(namespace, "unread", func(s socketio.Conn, data map[string]interface{}) {
		h.server.BroadcastToRoom(namespace, "notification", "unread", data)
	})

	h.server.OnEvent(namespace, "new", func(s socketio.Conn, data map[string]interface{}) {
		h.server.BroadcastToRoom(namespace, "notification", "new", data)
	})

	h.server.OnDisconnect(namespace, func(s socketio.Conn, reason string) {
		logger.Debug("Notification client disconnected: %s", reason)
	})

	h.namespaces[namespace] = true
	logger.Info("Initialized notification namespace")
}

// Start starts the socket server
func (h *Hub) Start() error {
	// Initialize standard namespaces
	h.InitMobileNamespace()
	h.InitQueueNamespace()
	h.InitNewQueueNamespace()
	h.InitNotificationNamespace()

	// Start cleanup task
	h.startCleanupTask()

	// Start serving
	go func() {
		if err := h.server.Serve(); err != nil {
			logger.Error("Socket.IO serve error: %v", err)
		}
	}()

	logger.Info("Socket.IO hub started")
	return nil
}

// startCleanupTask starts the periodic cleanup task
func (h *Hub) startCleanupTask() {
	h.cleanupTicker = time.NewTicker(h.config.CleanupInterval)

	go func() {
		for {
			select {
			case <-h.done:
				return
			case <-h.cleanupTicker.C:
				if err := h.stateMgr.CleanupExpiredData(context.Background()); err != nil {
					logger.Error("Cleanup task error: %v", err)
				}
			}
		}
	}()
}

// Stop stops the socket server
func (h *Hub) Stop() error {
	close(h.done)

	if h.cleanupTicker != nil {
		h.cleanupTicker.Stop()
	}

	h.eventMgr.Cleanup()

	if err := h.server.Close(); err != nil {
		return err
	}

	logger.Info("Socket.IO hub stopped")
	return nil
}

// Server returns the underlying Socket.IO server for HTTP handler
func (h *Hub) Server() *socketio.Server {
	return h.server
}

// BroadcastToRoom broadcasts a message to all clients in a room
func (h *Hub) BroadcastToRoom(namespace, room, event string, data interface{}) {
	h.server.BroadcastToRoom(namespace, room, event, data)
}

// GetRegisteredNamespaces returns list of registered namespaces
func (h *Hub) GetRegisteredNamespaces() []string {
	h.namespacesMu.RLock()
	defer h.namespacesMu.RUnlock()

	namespaces := make([]string, 0, len(h.namespaces))
	for ns := range h.namespaces {
		namespaces = append(namespaces, ns)
	}
	return namespaces
}
