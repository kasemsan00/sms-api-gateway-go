package socket

import (
	"context"
	"encoding/json"
	"sync"

	"api-gateway-go/internal/config"
	"api-gateway-go/pkg/logger"

	"github.com/redis/go-redis/v9"
)

// EventType represents cross-instance event types
type EventType string

const (
	EventUserDisconnect EventType = "user_disconnect"
	EventUserConnect    EventType = "user_connect"
	EventCarPosition    EventType = "car_position"
	EventChatMessage    EventType = "chat_message"
	EventRoomRecord     EventType = "room_record"
	EventQueueUpdate    EventType = "queue_update"
	EventNewCase        EventType = "newcase"
)

// CrossInstanceEvent represents an event to be published across instances
type CrossInstanceEvent struct {
	Namespace string      `json:"namespace"`
	Event     EventType   `json:"event"`
	Data      interface{} `json:"data"`
}

// EventHandler is a function that handles cross-instance events
type EventHandler func(event EventType, data interface{})

// CrossInstanceEventManager manages cross-instance event communication via Redis pub/sub
type CrossInstanceEventManager struct {
	redis       *config.RedisManager
	handlers    map[string]EventHandler
	handlerMu   sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	pubsub      *redis.PubSub
	channelName string
	initialized bool
	initMu      sync.Mutex
}

// NewCrossInstanceEventManager creates a new CrossInstanceEventManager
func NewCrossInstanceEventManager(redisMgr *config.RedisManager) *CrossInstanceEventManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &CrossInstanceEventManager{
		redis:       redisMgr,
		handlers:    make(map[string]EventHandler),
		ctx:         ctx,
		cancel:      cancel,
		channelName: "socket:events",
	}
}

// Initialize initializes the event manager and starts listening for events
func (m *CrossInstanceEventManager) Initialize() error {
	m.initMu.Lock()
	defer m.initMu.Unlock()

	if m.initialized {
		return nil
	}

	if m.redis == nil {
		logger.Warn("Redis not configured, cross-instance events disabled")
		m.initialized = true
		return nil
	}

	// Subscribe to the events channel using adapter client
	client := m.redis.AdapterClient()
	if client == nil {
		logger.Warn("Redis adapter client not available, cross-instance events disabled")
		m.initialized = true
		return nil
	}

	m.pubsub = client.Subscribe(m.ctx, m.channelName)

	// Start listening for messages in a goroutine
	go m.listenForMessages()

	m.initialized = true
	logger.Info("Cross-instance event manager initialized")
	return nil
}

// listenForMessages listens for messages on the Redis channel
func (m *CrossInstanceEventManager) listenForMessages() {
	if m.pubsub == nil {
		return
	}

	ch := m.pubsub.Channel()
	for {
		select {
		case <-m.ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			m.handleMessage(msg.Payload)
		}
	}
}

// handleMessage processes a received message
func (m *CrossInstanceEventManager) handleMessage(payload string) {
	var event CrossInstanceEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		logger.Error("Failed to unmarshal cross-instance event: %v", err)
		return
	}

	m.handlerMu.RLock()
	handler, exists := m.handlers[event.Namespace]
	m.handlerMu.RUnlock()

	if exists {
		handler(event.Event, event.Data)
	} else {
		logger.Debug("No handler registered for namespace: %s", event.Namespace)
	}
}

// PublishEvent publishes an event to all instances
func (m *CrossInstanceEventManager) PublishEvent(ctx context.Context, namespace string, event EventType, data interface{}) error {
	if m.redis == nil {
		return nil
	}

	client := m.redis.AdapterClient()
	if client == nil {
		return nil
	}

	eventData := CrossInstanceEvent{
		Namespace: namespace,
		Event:     event,
		Data:      data,
	}

	payload, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	return client.Publish(ctx, m.channelName, payload).Err()
}

// RegisterHandler registers an event handler for a namespace
func (m *CrossInstanceEventManager) RegisterHandler(namespace string, handler EventHandler) {
	m.handlerMu.Lock()
	defer m.handlerMu.Unlock()
	m.handlers[namespace] = handler
	logger.Debug("Registered event handler for namespace: %s", namespace)
}

// UnregisterHandler removes an event handler for a namespace
func (m *CrossInstanceEventManager) UnregisterHandler(namespace string) {
	m.handlerMu.Lock()
	defer m.handlerMu.Unlock()
	delete(m.handlers, namespace)
	logger.Debug("Unregistered event handler for namespace: %s", namespace)
}

// Cleanup cleans up resources
func (m *CrossInstanceEventManager) Cleanup() error {
	m.cancel()

	if m.pubsub != nil {
		if err := m.pubsub.Close(); err != nil {
			logger.Error("Error closing pubsub: %v", err)
			return err
		}
	}

	m.handlerMu.Lock()
	m.handlers = make(map[string]EventHandler)
	m.handlerMu.Unlock()

	logger.Info("Cross-instance event manager cleaned up")
	return nil
}
