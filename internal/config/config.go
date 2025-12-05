package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	MySQL MySQLConfig

	// Redis
	Redis RedisConfig

	// API
	APIURL string

	// SMS
	SMSEnable bool
	SMSAPIURL string

	// File
	RecordPath    string
	FileSizeLimit int64

	// Room
	JoinRoomRepeatDelay   time.Duration
	AutoCloseRoom         bool
	RoomDayDefaultTimeout int

	// LiveKit
	LiveKitAPIKey    string
	LiveKitAPISecret string
	LiveKitHost      string
	EgressLimit      int

	// Radio API
	RadioLocationAPIURL                string
	RadioLocationAPICredentialsUser    string
	RadioLocationAPICredentialsPassword string

	// Encode API
	EncodeAPI string

	// Custom
	CustomCharset string
}

// MySQLConfig holds MySQL database configuration
type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host       string
	Port       string
	Password   string
	AdapterDB  int
	StateDB    int
}

var cfg *Config

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	env := os.Getenv("ENVIRONMENT")
	if err := godotenv.Load(); err != nil {
		// Only warn in development mode
		if env == "" || env == "development" {
			// .env file not found, using system environment variables
		}
	}

	cfg = &Config{
		// Server
		Port:        getEnv("PORT", "5500"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Database
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_PASSWORD", ""),
			Database: getEnv("MYSQL_DATABASE", "conference"),
		},

		// Redis
		Redis: RedisConfig{
			Host:      getEnv("REDIS_HOST", "localhost"),
			Port:      getEnv("REDIS_PORT", "6379"),
			Password:  getEnv("REDIS_PASS", ""),
			AdapterDB: getEnvAsInt("REDIS_ADAPTER_DB", 1),
			StateDB:   getEnvAsInt("REDIS_STATE_DB", 2),
		},

		// API
		APIURL: getEnv("API_URL", "http://localhost:5500"),

		// SMS
		SMSEnable: getEnvAsBool("SMS_ENABLE", false),
		SMSAPIURL: getEnv("SMS_API_URL", ""),

		// File
		RecordPath:    getEnv("RECORD_PATH", "./record-file"),
		FileSizeLimit: getEnvAsInt64("FILE_SIZE_LIMIT", 524288000),

		// Room
		JoinRoomRepeatDelay:   time.Duration(getEnvAsInt("JOIN_ROOM_REPEAT_DELAY", 5000)) * time.Millisecond,
		AutoCloseRoom:         getEnvAsBool("AUTO_CLOSE_ROOM", false),
		RoomDayDefaultTimeout: getEnvAsInt("ROOM_DAY_DEFAULT_TIMEOUT", 365),

		// LiveKit
		LiveKitAPIKey:    getEnv("LIVEKIT_API_KEY", ""),
		LiveKitAPISecret: getEnv("LIVEKIT_API_SECRET", ""),
		LiveKitHost:      getEnv("LIVEKIT_HOST", ""),
		EgressLimit:      getEnvAsInt("EGRESS_LIMIT", 4),

		// Radio API
		RadioLocationAPIURL:                 getEnv("RADIO_LOCATION_API_URL", ""),
		RadioLocationAPICredentialsUser:     getEnv("RADIO_LOCATION_API_CREDENTIALS_USERNAME", ""),
		RadioLocationAPICredentialsPassword: getEnv("RADIO_LOCATION_API_CREDENTIALS_PASSWORD", ""),

		// Encode API
		EncodeAPI: getEnv("ENCODE_API", "http://encode-api:5600"),

		// Custom
		CustomCharset: getEnv("CUSTOM_CHARSET", "ABCDEFGHIJKLMOPQRSTUVWXYZabcdefghijklmopqrstuvwxyz"),
	}

	return cfg, nil
}

// Get returns the loaded configuration
func Get() *Config {
	if cfg == nil {
		cfg, _ = Load()
	}
	return cfg
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
