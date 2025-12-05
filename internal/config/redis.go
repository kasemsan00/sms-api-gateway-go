package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// RedisManager manages multiple Redis connections
type RedisManager struct {
	clients map[int]*redis.Client
	cfg     *Config
}

var redisManager *RedisManager

// InitRedis initializes Redis connections
func InitRedis(cfg *Config) (*RedisManager, error) {
	redisManager = &RedisManager{
		clients: make(map[int]*redis.Client),
		cfg:     cfg,
	}

	// Initialize default client (DB 0)
	if err := redisManager.initClient(0); err != nil {
		return nil, err
	}

	// Initialize adapter client
	if err := redisManager.initClient(cfg.Redis.AdapterDB); err != nil {
		return nil, err
	}

	// Initialize state client
	if err := redisManager.initClient(cfg.Redis.StateDB); err != nil {
		return nil, err
	}

	log.Info().Msgf("Connected to Redis: %s:%s", cfg.Redis.Host, cfg.Redis.Port)

	return redisManager, nil
}

// initClient initializes a Redis client for a specific database
func (rm *RedisManager) initClient(db int) error {
	if _, exists := rm.clients[db]; exists {
		return nil
	}

	addr := fmt.Sprintf("%s:%s", rm.cfg.Redis.Host, rm.cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        rm.cfg.Redis.Password,
		DB:              db,
		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond,
		MaxRetryBackoff: 1 * time.Second,
		PoolSize:        10,
		MinIdleConns:    5,
		PoolTimeout:     30 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis DB %d: %w", db, err)
	}

	rm.clients[db] = client
	return nil
}

// GetRedis returns the Redis manager
func GetRedis() *RedisManager {
	return redisManager
}

// Client returns the default Redis client (DB 0)
func (rm *RedisManager) Client() *redis.Client {
	return rm.clients[0]
}

// GetClient returns a Redis client for a specific database
func (rm *RedisManager) GetClient(db int) *redis.Client {
	if client, exists := rm.clients[db]; exists {
		return client
	}
	// Initialize on demand
	if err := rm.initClient(db); err != nil {
		log.Error().Err(err).Msgf("Failed to initialize Redis client for DB %d", db)
		return nil
	}
	return rm.clients[db]
}

// AdapterClient returns the Socket.IO adapter Redis client
func (rm *RedisManager) AdapterClient() *redis.Client {
	return rm.GetClient(rm.cfg.Redis.AdapterDB)
}

// StateClient returns the Socket.IO state Redis client
func (rm *RedisManager) StateClient() *redis.Client {
	return rm.GetClient(rm.cfg.Redis.StateDB)
}

// Close closes all Redis connections
func (rm *RedisManager) Close() error {
	log.Info().Msg("Closing Redis connections")
	for db, client := range rm.clients {
		if err := client.Close(); err != nil {
			log.Error().Err(err).Msgf("Failed to close Redis client for DB %d", db)
		}
	}
	return nil
}

// Health checks Redis health
func (rm *RedisManager) Health(ctx context.Context) error {
	for db, client := range rm.clients {
		if err := client.Ping(ctx).Err(); err != nil {
			return fmt.Errorf("Redis DB %d health check failed: %w", db, err)
		}
	}
	return nil
}

// Publish publishes a message to a channel
func (rm *RedisManager) Publish(ctx context.Context, channel string, message interface{}) error {
	return rm.Client().Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to channels
func (rm *RedisManager) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return rm.Client().Subscribe(ctx, channels...)
}

// Set sets a key-value pair with optional expiration
func (rm *RedisManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rm.Client().Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (rm *RedisManager) Get(ctx context.Context, key string) (string, error) {
	return rm.Client().Get(ctx, key).Result()
}

// Delete deletes keys
func (rm *RedisManager) Delete(ctx context.Context, keys ...string) error {
	return rm.Client().Del(ctx, keys...).Err()
}

// Exists checks if keys exist
func (rm *RedisManager) Exists(ctx context.Context, keys ...string) (int64, error) {
	return rm.Client().Exists(ctx, keys...).Result()
}
