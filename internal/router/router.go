package router

import (
	"api-gateway-go/internal/config"
	"api-gateway-go/internal/handler"
	"api-gateway-go/internal/repository"
	"api-gateway-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

// Handlers holds all handler instances
type Handlers struct {
	Auth         *handler.AuthHandler
	Room         *handler.RoomHandler
	User         *handler.UserHandler
	Link         *handler.LinkHandler
	System       *handler.SystemHandler
	Chat         *handler.ChatHandler
	Notification *handler.NotificationHandler
	Record       *handler.RecordHandler
	Car          *handler.CarHandler
	Case         *handler.CaseHandler
	Radio        *handler.RadioHandler
	Stats        *handler.StatsHandler
	Upload       *handler.UploadHandler
	Webhook      *handler.WebhookHandler
	Test         *handler.TestHandler
}

// SetupRoutes configures all routes for the application
// Routes are aligned with Node.js API Gateway
func SetupRoutes(
	app *fiber.App,
	handlers *Handlers,
	authService *service.AuthService,
	cfg *config.Config,
	db *config.Database,
	redisMgr *config.RedisManager,
	livekitMgr *config.LiveKitManager,
	recordRepo *repository.RecordRepository,
) {
	// ============================================
	// System routes (from index.routes.js)
	// ============================================
	app.Get("/", handlers.System.Root)
	app.Get("/health", handlers.System.HealthCheck)
	app.Post("/log", handlers.System.AddLog)
	app.Get("/status", handlers.System.GetStatus)
	app.Get("/status/cron", handlers.System.GetCronJobStatus)
	app.Get("/service", handlers.System.GetServiceInfo)
	app.Post("/webhook", handlers.Webhook.HandleGenericWebhook)
	app.Post("/sms/custom", handlers.Link.SendCustomMessage)
	app.Get("/namespace", handlers.System.GetNamespaces)

	// ============================================
	// Auth routes (from auth.routes.js)
	// Node.js only has: POST /auth/verifyuser
	// ============================================
	auth := app.Group("/auth")
	auth.Post("/verifyuser", handlers.Auth.VerifyUser)

	// ============================================
	// Room routes (from room.routes.js)
	// ============================================
	room := app.Group("/room")
	room.Get("/detail", handlers.Room.GetRoomDetail)
	room.Get("/listrooms", handlers.Room.ListRooms)
	room.Get("/checkexpired", handlers.Room.CheckExpired)
	room.Get("/verifytoken", handlers.Room.VerifyToken)
	room.Get("/picture", handlers.Room.GetRoomPicture)
	room.Post("/updateuser", handlers.Room.UpdateUser)
	room.Post("/deleteroom", handlers.Room.DeleteRoom)
	room.Put("/updatetype", handlers.Room.UpdateType)
	room.Put("/updatestatus", handlers.Room.UpdateStatus)
	room.Put("/close", handlers.Room.CloseRoom)
	room.Post("/verifyuser", handlers.Room.VerifyUser)

	// ============================================
	// User routes (from user.routes.js)
	// ============================================
	user := app.Group("/user")
	user.Get("/getuseralreadyinroom", handlers.User.GetUserAlreadyInRoom)
	user.Get("/getuserdetail", handlers.User.GetUserDetail)
	user.Get("/listparticipants", handlers.User.ListParticipants)
	user.Post("/generate", handlers.User.GenerateUser)
	user.Post("/joingenerate", handlers.User.JoinGenerate)
	user.Post("/generateChatUser", handlers.User.GenerateChatUser)
	user.Post("/updateparticipants", handlers.User.UpdateParticipants)
	user.Post("/mutepublishedtrack", handlers.User.MutePublishedTrack)
	user.Post("/removeParticipant", handlers.User.RemoveParticipant)
	user.Get("/log", handlers.User.GetUserLog)
	user.Put("/handle/track", handlers.User.HandleTrack)

	// ============================================
	// Link routes (from link.routes.js)
	// ============================================
	link := app.Group("/link")
	link.Get("/getdetail", handlers.Link.GetLinkDetail)
	link.Get("/history", handlers.Link.GetLinkHistory)
	link.Post("/create", handlers.Link.CreateLink)
	link.Post("/create/hls", handlers.Link.CreateHLSLink)
	link.Post("/update/latlng", handlers.Link.UpdateLatLng)
	link.Post("/multilatlng/send", handlers.Link.MultiLatLng)
	link.Get("/share", handlers.Link.GetShareURL)
	link.Post("/cartracking", handlers.Link.CarTracking)
	link.Get("/get/domain", handlers.Link.GetDomain)
	link.Get("/list", handlers.Link.GetLinkList)

	// ============================================
	// Chat routes (from chat.routes.js)
	// Node.js only has: GET /chat/history, GET /chat/notification
	// ============================================
	chat := app.Group("/chat")
	chat.Get("/history", handlers.Chat.GetHistory)
	chat.Get("/notification", handlers.Chat.GetNotification)

	// ============================================
	// Notification routes (from notification.routes.js)
	// ============================================
	notification := app.Group("/notification")
	notification.Get("/events", handlers.Notification.ListNotifications)
	notification.Put("/update/:notificationId", handlers.Notification.MarkAsRead)
	notification.Get("/unread", handlers.Notification.GetUnread)
	notification.Get("/:notificationId", handlers.Notification.GetByID)
	notification.Post("/", handlers.Notification.Create)

	// ============================================
	// Record routes (from record.routes.js)
	// ============================================
	record := app.Group("/record")
	record.Get("/request", handlers.Record.StartRecord)
	record.Get("/list", handlers.Record.ListEgress)
	record.Get("/stopall", handlers.Record.StopAllActive)
	record.Get("/file", handlers.Record.GetFileHistory)
	record.Get("/check", handlers.Record.CheckEgressAvailable)

	// ============================================
	// Car tracking routes (from car.routes.js)
	// ============================================
	car := app.Group("/car")
	car.Get("/task", handlers.Car.GetTaskDetail)
	car.Put("/task", handlers.Car.UpdateTask)
	car.Post("/task", handlers.Car.CreateTask)
	car.Get("/list", handlers.Car.ListTasks)

	// ============================================
	// Case routes (from case.routes.js)
	// ============================================
	caseRoutes := app.Group("/case")
	caseRoutes.Post("/create", handlers.Case.CreateCase)
	caseRoutes.Get("/get", handlers.Case.GetCaseByID)
	caseRoutes.Get("/history", handlers.Case.GetCaseHistory)
	caseRoutes.Put("/update", handlers.Case.UpdateCase)

	// ============================================
	// Radio routes (from radio.routes.js)
	// ============================================
	radio := app.Group("/radio")
	radio.Get("/device", handlers.Radio.ListDevices)
	radio.Get("/device/:id", handlers.Radio.GetDeviceByID)
	radio.Get("/location", handlers.Radio.ListLocations)
	radio.Get("/location/:radioNo", handlers.Radio.GetLocationByRadioNo)

	// ============================================
	// Stats routes (from stats.routes.js)
	// ============================================
	stats := app.Group("/stats")
	stats.Get("/summary", handlers.Stats.GetSummary)
	stats.Get("/device", handlers.Stats.GetDeviceStats)
	stats.Get("/type", handlers.Stats.GetTypeStats)
	stats.Get("/gen", handlers.Stats.GetDailyStats)
	stats.Get("/generate", handlers.Stats.GetDailyStats)
	stats.Get("/user", handlers.Stats.GetUserStats)
	stats.Get("/case", handlers.Stats.GetCaseStats)

	// ============================================
	// Upload routes (from upload.routes.js)
	// ============================================
	upload := app.Group("/upload")
	upload.Post("/file", handlers.Upload.UploadFile)
	upload.Post("/video", handlers.Upload.UploadVideo)
	upload.Get("/list", handlers.Upload.VideoList)
	upload.Post("/sms/send", handlers.Upload.SendSMS)

	// ============================================
	// Service routes (from service.routes.js)
	// ============================================
	serviceRoutes := app.Group("/service")
	serviceRoutes.Get("/get", handlers.System.GetServiceByRoom)
	serviceRoutes.Put("/update", handlers.System.UpdateService)

	// ============================================
	// Webhook routes
	// ============================================
	webhook := app.Group("/webhook")
	webhook.Post("/livekit", handlers.Webhook.HandleLiveKitWebhook)

	// ============================================
	// Test routes (from index.routes.js)
	// ============================================
	test := app.Group("/test")
	test.Get("/", handlers.Test.Ping)
	test.Get("/get/namespace", handlers.Test.GetAllNamespaces)
	test.Get("/redis/connection", handlers.Test.TestRedis)
	test.Get("/redis/operations", handlers.Test.TestRedis)
	test.Delete("/redis/clear", handlers.Test.ClearRedisTestData)
	test.Get("/mp4/queue", handlers.Test.GetMP4ProcessingQueue)
	test.Delete("/mp4/queue/:recordId", handlers.Test.RemoveFromMP4ProcessingQueue)
	test.Delete("/mp4/queue", handlers.Test.ClearMP4ProcessingQueue)

	// Static file serving
	app.Static("/logo", "./logo")
	app.Static("/videos", "./uploads/videos")
	app.Static("/images", "./uploads/images")
	app.Static("/thumbnails", "./uploads/thumbnails")
	app.Static("/files", "./uploads/files")
	app.Static("/record", "./record-file")
}
