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
	// System routes
	app.Get("/", handlers.System.Root)
	app.Get("/health", handlers.System.HealthCheck)
	app.Get("/status", handlers.System.GetStatus)
	app.Get("/service", handlers.System.GetServiceInfo)
	app.Get("/namespace", handlers.System.GetNamespaces)
	app.Post("/log", handlers.System.AddLog)

	// Auth routes
	auth := app.Group("/auth")
	auth.Get("/create", handlers.Auth.CreateToken)
	auth.Get("/verify", handlers.Auth.VerifyToken)
	auth.Post("/verifyuser", handlers.Auth.VerifyUser)

	// Room routes
	room := app.Group("/room")
	room.Get("/detail", handlers.Room.GetRoomDetail)
	room.Get("/listrooms", handlers.Room.ListRooms)
	room.Get("/checkexpired", handlers.Room.CheckExpired)
	room.Get("/verifytoken", handlers.Room.VerifyToken)
	room.Get("/picture", handlers.Room.GetRoomPicture)
	room.Post("/create", handlers.Room.CreateRoom)
	room.Post("/updateuser", handlers.Room.UpdateUser)
	room.Post("/deleteroom", handlers.Room.DeleteRoom)
	room.Put("/updatetype", handlers.Room.UpdateType)
	room.Put("/updatestatus", handlers.Room.UpdateStatus)
	room.Put("/close", handlers.Room.CloseRoom)
	room.Put("/recordstatus", handlers.Room.UpdateRecordStatus)

	// User routes
	user := app.Group("/user")
	user.Get("/getuseralreadyinroom", handlers.User.GetUserAlreadyInRoom)
	user.Get("/getuserdetail", handlers.User.GetUserDetail)
	user.Get("/listparticipants", handlers.User.ListParticipants)
	user.Get("/log", handlers.User.GetUserLog)
	user.Post("/generate", handlers.User.GenerateUser)
	user.Post("/joingenerate", handlers.User.JoinGenerate)
	user.Post("/generateChatUser", handlers.User.GenerateChatUser)
	user.Post("/updateparticipants", handlers.User.UpdateParticipants)
	user.Post("/mutepublishedtrack", handlers.User.MutePublishedTrack)
	user.Post("/removeParticipant", handlers.User.RemoveParticipant)
	user.Put("/handle/track", handlers.User.HandleTrack)

	// Link routes
	link := app.Group("/link")
	link.Get("/getdetail", handlers.Link.GetLinkDetail)
	link.Get("/history", handlers.Link.GetLinkHistory)
	link.Get("/share", handlers.Link.GetShareURL)
	link.Get("/get/domain", handlers.Link.GetDomain)
	link.Get("/list", handlers.Link.GetLinkList)
	link.Post("/create", handlers.Link.CreateLink)
	link.Post("/create/hls", handlers.Link.CreateHLSLink)
	link.Post("/update/latlng", handlers.Link.UpdateLatLng)
	link.Post("/multilatlng/send", handlers.Link.MultiLatLng)
	link.Post("/cartracking", handlers.Link.CarTracking)

	// Chat routes
	chat := app.Group("/chat")
	chat.Get("/history", handlers.Chat.GetHistory)
	chat.Get("/notification", handlers.Chat.GetNotification)
	chat.Get("/count", handlers.Chat.GetMessageCount)
	chat.Post("/message", handlers.Chat.AddMessage)
	chat.Delete("/messages", handlers.Chat.DeleteMessages)

	// Notification routes
	notification := app.Group("/notification")
	notification.Get("/list", handlers.Notification.ListNotifications)
	notification.Get("/user", handlers.Notification.GetByUserName)
	notification.Get("/unread", handlers.Notification.GetUnread)
	notification.Get("/unreadcount", handlers.Notification.GetUnreadCount)
	notification.Get("/:id", handlers.Notification.GetByID)
	notification.Post("/create", handlers.Notification.Create)
	notification.Put("/read/:id", handlers.Notification.MarkAsRead)
	notification.Put("/readall", handlers.Notification.MarkAllAsRead)
	notification.Delete("/:id", handlers.Notification.Delete)

	// Record routes
	record := app.Group("/record")
	record.Get("/listegress", handlers.Record.ListEgress)
	record.Get("/available", handlers.Record.CheckEgressAvailable)
	record.Get("/queue", handlers.Record.GetRecordQueue)
	record.Get("/activecount", handlers.Record.GetActiveRecordCount)
	record.Get("/filehistory", handlers.Record.GetFileHistory)
	record.Get("/room", handlers.Record.GetRecordByRoom)
	record.Get("/detail/:id", handlers.Record.GetRecordDetail)
	record.Post("/start", handlers.Record.StartRecord)
	record.Post("/stop", handlers.Record.StopRecord)
	record.Post("/stopall", handlers.Record.StopAllActive)

	// Car tracking routes
	car := app.Group("/car")
	car.Get("/list", handlers.Car.ListTasks)
	car.Get("/task/:id", handlers.Car.GetTaskDetail)
	car.Get("/uid/:uid", handlers.Car.GetTaskByUID)
	car.Get("/room/:room", handlers.Car.GetTaskByRoom)
	car.Get("/position/:room", handlers.Car.GetCarPosition)
	car.Get("/latlng/:room", handlers.Car.GetUserLatLng)
	car.Post("/task", handlers.Car.CreateTask)
	car.Post("/position", handlers.Car.UpdatePosition)
	car.Put("/task/:id", handlers.Car.UpdateTask)
	car.Delete("/task/:id", handlers.Car.DeleteTask)

	// Case routes
	caseRoutes := app.Group("/case")
	caseRoutes.Get("/history", handlers.Case.GetCaseHistory)
	caseRoutes.Get("/historycount", handlers.Case.GetCaseHistoryCount)
	caseRoutes.Get("/roomname", handlers.Case.GetRoomName)
	caseRoutes.Get("/service/:service", handlers.Case.GetCasesByService)
	caseRoutes.Get("/caseid/:caseId", handlers.Case.GetCaseByCaseID)
	caseRoutes.Get("/room/:roomId", handlers.Case.GetCaseByRoomID)
	caseRoutes.Get("/:id", handlers.Case.GetCaseByID)
	caseRoutes.Post("/create", handlers.Case.CreateCase)
	caseRoutes.Put("/status/:caseId", handlers.Case.UpdateCaseStatus)
	caseRoutes.Put("/:id", handlers.Case.UpdateCase)
	caseRoutes.Delete("/:id", handlers.Case.DeleteCase)

	// Radio routes
	radio := app.Group("/radio")
	radio.Get("/devices", handlers.Radio.ListDevices)
	radio.Get("/device/:id", handlers.Radio.GetDeviceByID)
	radio.Get("/device/deviceid/:deviceId", handlers.Radio.GetDeviceByDeviceID)
	radio.Get("/locations", handlers.Radio.ListLocations)
	radio.Get("/location/:radioNo", handlers.Radio.GetLocationByRadioNo)
	radio.Post("/device", handlers.Radio.CreateDevice)
	radio.Post("/location", handlers.Radio.CreateLocation)
	radio.Put("/device/:id", handlers.Radio.UpdateDevice)
	radio.Put("/device/location", handlers.Radio.UpdateDeviceLocation)
	radio.Delete("/device/:id", handlers.Radio.DeleteDevice)

	// Stats routes
	stats := app.Group("/stats")
	stats.Get("/summary", handlers.Stats.GetSummary)
	stats.Get("/device", handlers.Stats.GetDeviceStats)
	stats.Get("/type", handlers.Stats.GetTypeStats)
	stats.Get("/user", handlers.Stats.GetUserStats)
	stats.Get("/case", handlers.Stats.GetCaseStats)
	stats.Get("/daily", handlers.Stats.GetDailyStats)
	stats.Get("/monthly", handlers.Stats.GetMonthlyStats)
	stats.Get("/all", handlers.Stats.GetAll)

	// Upload routes
	upload := app.Group("/upload")
	upload.Post("/file", handlers.Upload.UploadFile)
	upload.Post("/image", handlers.Upload.UploadImage)
	upload.Post("/video", handlers.Upload.UploadVideo)
	upload.Post("/multiple", handlers.Upload.UploadMultiple)
	upload.Get("/exists", handlers.Upload.CheckFileExists)
	upload.Delete("/file", handlers.Upload.DeleteFile)
	// upload.Delete("/file", middleware.AuthMiddleware(authService), handlers.Upload.DeleteFile)

	// Webhook routes
	webhook := app.Group("/webhook")
	webhook.Post("/livekit", handlers.Webhook.HandleLiveKitWebhook)
	webhook.Post("/generic", handlers.Webhook.HandleGenericWebhook)

	// Test routes
	test := app.Group("/test")
	test.Get("/ping", handlers.Test.Ping)
	test.Post("/echo", handlers.Test.Echo)
	test.Get("/database", handlers.Test.TestDatabase)
	test.Get("/redis", handlers.Test.TestRedis)
	test.Get("/livekit", handlers.Test.TestLiveKit)
	test.Get("/all", handlers.Test.TestAll)
	test.Get("/config", handlers.Test.GetConfig)

	// Static file serving
	app.Static("/logo", "./logo")
	app.Static("/videos", "./uploads/videos")
	app.Static("/images", "./uploads/images")
	app.Static("/thumbnails", "./uploads/thumbnails")
	app.Static("/files", "./uploads/files")
	app.Static("/record", "./record-file")
}
