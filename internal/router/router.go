package router

import (
	"api-gateway-go/internal/handler"
	"api-gateway-go/internal/middleware"
	"api-gateway-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

// Handlers holds all handler instances
type Handlers struct {
	Auth   *handler.AuthHandler
	Room   *handler.RoomHandler
	User   *handler.UserHandler
	Link   *handler.LinkHandler
	System *handler.SystemHandler
}

// SetupRoutes configures all routes for the application
func SetupRoutes(app *fiber.App, handlers *Handlers, authService *service.AuthService) {
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
	room.Post("/create", middleware.AuthMiddleware(authService), handlers.Room.CreateRoom)
	room.Post("/updateuser", middleware.AuthMiddleware(authService), handlers.Room.UpdateUser)
	room.Post("/deleteroom", middleware.AuthMiddleware(authService), handlers.Room.DeleteRoom)
	room.Put("/updatetype", middleware.AuthMiddleware(authService), handlers.Room.UpdateType)
	room.Put("/updatestatus", handlers.Room.UpdateStatus)
	room.Put("/close", handlers.Room.CloseRoom)
	room.Put("/recordstatus", handlers.Room.UpdateRecordStatus)

	// User routes
	user := app.Group("/user")
	user.Get("/getuseralreadyinroom", handlers.User.GetUserAlreadyInRoom)
	user.Get("/getuserdetail", handlers.User.GetUserDetail)
	user.Get("/listparticipants", handlers.User.ListParticipants)
	user.Get("/log", handlers.User.GetUserLog)
	user.Post("/generate", middleware.OptionalAuthMiddleware(authService), handlers.User.GenerateUser)
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
	link.Post("/create", middleware.AuthMiddleware(authService), handlers.Link.CreateLink)
	link.Post("/create/hls", middleware.AuthMiddleware(authService), handlers.Link.CreateHLSLink)
	link.Post("/update/latlng", handlers.Link.UpdateLatLng)
	link.Post("/multilatlng/send", handlers.Link.MultiLatLng)
	link.Post("/cartracking", handlers.Link.CarTracking)

	// Static file serving
	app.Static("/logo", "./logo")
	app.Static("/videos", "./uploads/videos")
	app.Static("/images", "./uploads/images")
	app.Static("/thumbnails", "./uploads/thumbnails")
	app.Static("/files", "./uploads/files")
	app.Static("/record", "./record-file")
}
