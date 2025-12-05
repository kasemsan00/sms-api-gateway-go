package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/handler"
	"api-gateway-go/internal/middleware"
	"api-gateway-go/internal/repository"
	"api-gateway-go/internal/router"
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize logger
	logger.Init(cfg.Environment)
	log.Info().Msg("Starting SMS API Gateway")

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Initialize Redis
	redis, err := config.InitRedis(cfg)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to Redis, some features may be unavailable")
	}
	if redis != nil {
		defer redis.Close()
	}

	// Initialize LiveKit
	livekit, err := config.InitLiveKit(cfg)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to initialize LiveKit, video features may be unavailable")
	}

	// Initialize repositories
	roomRepo := repository.NewRoomRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	linkRepo := repository.NewLinkRepository(db.DB)
	chatRepo := repository.NewChatRepository(db.DB)
	notificationRepo := repository.NewNotificationRepository(db.DB)
	recordRepo := repository.NewRecordRepository(db.DB)
	carRepo := repository.NewCarRepository(db.DB)
	caseRepo := repository.NewCaseRepository(db.DB)
	radioRepo := repository.NewRadioRepository(db.DB)
	statsRepo := repository.NewStatsRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(cfg.LiveKitAPISecret)
	roomService := service.NewRoomService(roomRepo, livekit, cfg)
	userService := service.NewUserService(userRepo, roomRepo, livekit, cfg)
	linkService := service.NewLinkService(linkRepo, roomRepo, cfg)
	chatService := service.NewChatService(chatRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	recordService := service.NewRecordService(recordRepo, roomRepo, livekit, cfg)
	carService := service.NewCarService(carRepo)
	caseService := service.NewCaseService(caseRepo)
	radioService := service.NewRadioService(radioRepo)
	statsService := service.NewStatsService(statsRepo)
	smsService := service.NewSMSService(cfg)
	fileService := service.NewFileService(cfg)

	// Initialize crontab service
	crontabService := service.NewCrontabService(roomService, linkService, livekit, cfg)
	if err := crontabService.InitCronJobs(); err != nil {
		log.Warn().Err(err).Msg("Failed to initialize cron jobs")
	}
	defer crontabService.Stop()

	// Suppress unused variable warnings - these will be used when more handlers are added
	_ = smsService

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "SMS API Gateway",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		BodyLimit:    int(cfg.FileSizeLimit),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status":  "FAIL",
				"message": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.LoggerMiddleware())

	// Initialize handlers
	handlers := &router.Handlers{
		Auth:         handler.NewAuthHandler(authService),
		Room:         handler.NewRoomHandler(roomService, authService),
		User:         handler.NewUserHandler(userService),
		Link:         handler.NewLinkHandler(linkService, userService),
		System:       handler.NewSystemHandler(db, redis, livekit, crontabService, cfg),
		Chat:         handler.NewChatHandler(chatService),
		Notification: handler.NewNotificationHandler(notificationService),
		Record:       handler.NewRecordHandler(recordService),
		Car:          handler.NewCarHandler(carService),
		Case:         handler.NewCaseHandler(caseService),
		Radio:        handler.NewRadioHandler(radioService),
		Stats:        handler.NewStatsHandler(statsService),
		Upload:       handler.NewUploadHandler(fileService),
		Webhook:      handler.NewWebhookHandler(roomService, userService, recordService, recordRepo, livekit, cfg),
		Test:         handler.NewTestHandler(cfg, db, redis, livekit),
	}

	// Setup routes
	router.SetupRoutes(app, handlers, authService, cfg, db, redis, livekit, recordRepo)

	// Start server
	go func() {
		listenAddr := fmt.Sprintf(":%s", cfg.Port)
		log.Info().Msgf("🚀 Server starting on http://localhost%s", listenAddr)

		if err := app.Listen(listenAddr); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown
	gracefulShutdown(app, db, redis, crontabService)
}

// gracefulShutdown handles graceful shutdown of the application
func gracefulShutdown(app *fiber.App, db *config.Database, redis *config.RedisManager, crontab *service.CrontabService) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info().Msg("Starting graceful shutdown...")

	// Phase 1: Stop cron jobs
	log.Info().Msg("Stopping cron jobs...")
	crontab.Stop()

	// Phase 2: Close Redis connections
	if redis != nil {
		log.Info().Msg("Closing Redis connections...")
		redis.Close()
	}

	// Phase 3: Close database pool
	if db != nil {
		log.Info().Msg("Closing database connection pool...")
		db.Close()
	}

	// Phase 4: Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info().Msg("Shutting down HTTP server...")
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server shutdown complete")
}
