package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the global logger
func Init(environment string) {
	// Set timezone to Bangkok
	loc, _ := time.LoadLocation("Asia/Bangkok")
	time.Local = loc

	// Configure zerolog
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if environment == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Console writer for pretty output
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
}

// InitWithFile initializes logger with file output
func InitWithFile(environment, logPath string) error {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	time.Local = loc

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if environment == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Create log file
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Multi-writer for console and file
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}

	multi := io.MultiWriter(consoleWriter, file)
	log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()

	return nil
}

// Debug logs a debug message
func Debug(msg string, args ...interface{}) {
	log.Debug().Msgf(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	log.Info().Msgf(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	log.Warn().Msgf(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...interface{}) {
	log.Error().Msgf(msg, args...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, args ...interface{}) {
	log.Fatal().Msgf(msg, args...)
}

// WithFields returns a logger with additional fields
func WithFields(fields map[string]interface{}) zerolog.Logger {
	ctx := log.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return ctx.Logger()
}

// HTTP logs HTTP request info
func HTTP(method, path string, status int, latency time.Duration) {
	log.Info().
		Str("method", method).
		Str("path", path).
		Int("status", status).
		Dur("latency", latency).
		Msg("HTTP Request")
}

// Socket logs socket events
func Socket(event, namespace, socketID string, data interface{}) {
	log.Debug().
		Str("event", event).
		Str("namespace", namespace).
		Str("socketId", socketID).
		Interface("data", data).
		Msg("Socket Event")
}
