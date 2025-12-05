package socket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"api-gateway-go/internal/config"
	"api-gateway-go/pkg/logger"
)

// SocketState represents the state of a socket connection
type SocketState struct {
	SocketID  string    `json:"socketId"`
	Room      string    `json:"room"`
	Identity  string    `json:"identity"`
	UserName  string    `json:"userName"`
	UserType  string    `json:"userType"`
	CreatedAt time.Time `json:"createdAt"`
}

// SocketStateManager manages socket state in Redis
type SocketStateManager struct {
	redis      *config.RedisManager
	localState sync.Map // fallback for when Redis is not available
}

// NewSocketStateManager creates a new SocketStateManager
func NewSocketStateManager(redisMgr *config.RedisManager) *SocketStateManager {
	return &SocketStateManager{
		redis: redisMgr,
	}
}

// SetUserSession stores a user session in Redis
func (m *SocketStateManager) SetUserSession(ctx context.Context, room, identity string, state *SocketState) error {
	if m.redis == nil {
		// Use local state as fallback
		key := room + ":" + identity
		m.localState.Store(key, state)
		return nil
	}

	key := "socket:session:" + room + ":" + identity
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return m.redis.StateClient().Set(ctx, key, data, 24*time.Hour).Err()
}

// GetUserSession retrieves a user session from Redis
func (m *SocketStateManager) GetUserSession(ctx context.Context, room, identity string) (*SocketState, error) {
	if m.redis == nil {
		// Use local state as fallback
		key := room + ":" + identity
		if val, ok := m.localState.Load(key); ok {
			return val.(*SocketState), nil
		}
		return nil, nil
	}

	key := "socket:session:" + room + ":" + identity
	data, err := m.redis.StateClient().Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var state SocketState
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// RemoveUserSession removes a user session from Redis
func (m *SocketStateManager) RemoveUserSession(ctx context.Context, room, identity string) error {
	if m.redis == nil {
		key := room + ":" + identity
		m.localState.Delete(key)
		return nil
	}

	key := "socket:session:" + room + ":" + identity
	return m.redis.StateClient().Del(ctx, key).Err()
}

// GetRoomSessions gets all sessions in a room
func (m *SocketStateManager) GetRoomSessions(ctx context.Context, room string) ([]*SocketState, error) {
	if m.redis == nil {
		var sessions []*SocketState
		m.localState.Range(func(key, value interface{}) bool {
			state := value.(*SocketState)
			if state.Room == room {
				sessions = append(sessions, state)
			}
			return true
		})
		return sessions, nil
	}

	pattern := "socket:session:" + room + ":*"
	keys, err := m.redis.StateClient().Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var sessions []*SocketState
	for _, key := range keys {
		data, err := m.redis.StateClient().Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var state SocketState
		if err := json.Unmarshal([]byte(data), &state); err != nil {
			continue
		}
		sessions = append(sessions, &state)
	}

	return sessions, nil
}

// CleanupExpiredData cleans up expired session data
func (m *SocketStateManager) CleanupExpiredData(ctx context.Context) error {
	if m.redis == nil {
		return nil
	}

	pattern := "socket:session:*"
	keys, err := m.redis.StateClient().Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	cleaned := 0
	for _, key := range keys {
		ttl, err := m.redis.StateClient().TTL(ctx, key).Result()
		if err != nil {
			continue
		}

		// If TTL is -1, set a new TTL (24 hours)
		if ttl == -1 {
			m.redis.StateClient().Expire(ctx, key, 24*time.Hour)
			cleaned++
		}
	}

	if cleaned > 0 {
		logger.Info("Cleaned up %d socket sessions without TTL", cleaned)
	}

	return nil
}

// GetActiveRooms gets list of active rooms
func (m *SocketStateManager) GetActiveRooms(ctx context.Context) ([]string, error) {
	if m.redis == nil {
		roomSet := make(map[string]bool)
		m.localState.Range(func(key, value interface{}) bool {
			state := value.(*SocketState)
			roomSet[state.Room] = true
			return true
		})

		var rooms []string
		for room := range roomSet {
			rooms = append(rooms, room)
		}
		return rooms, nil
	}

	pattern := "socket:session:*"
	keys, err := m.redis.StateClient().Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	roomSet := make(map[string]bool)
	for _, key := range keys {
		// Key format: socket:session:{room}:{identity}
		// Extract room from key
		data, err := m.redis.StateClient().Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var state SocketState
		if err := json.Unmarshal([]byte(data), &state); err != nil {
			continue
		}
		roomSet[state.Room] = true
	}

	var rooms []string
	for room := range roomSet {
		rooms = append(rooms, room)
	}

	return rooms, nil
}
