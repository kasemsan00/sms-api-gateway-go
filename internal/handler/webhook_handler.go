package handler

import (
	"api-gateway-go/internal/config"
	"api-gateway-go/internal/repository"
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/livekit/protocol/livekit"
	"github.com/rs/zerolog/log"
)

// WebhookHandler handles LiveKit webhook events
type WebhookHandler struct {
	roomService   *service.RoomService
	userService   *service.UserService
	recordService *service.RecordService
	recordRepo    *repository.RecordRepository
	livekitMgr    *config.LiveKitManager
	cfg           *config.Config
}

// NewWebhookHandler creates a new WebhookHandler
func NewWebhookHandler(
	roomService *service.RoomService,
	userService *service.UserService,
	recordService *service.RecordService,
	recordRepo *repository.RecordRepository,
	livekitMgr *config.LiveKitManager,
	cfg *config.Config,
) *WebhookHandler {
	return &WebhookHandler{
		roomService:   roomService,
		userService:   userService,
		recordService: recordService,
		recordRepo:    recordRepo,
		livekitMgr:    livekitMgr,
		cfg:           cfg,
	}
}

// HandleLiveKitWebhook handles LiveKit webhook events
// POST /webhook/livekit
func (h *WebhookHandler) HandleLiveKitWebhook(c *fiber.Ctx) error {
	// Get authorization header for validation
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.ErrorResponseWithStatus(c, fiber.StatusUnauthorized, "Authorization required")
	}

	// Parse webhook event
	body := c.Body()

	// Log webhook received
	log.Info().
		Str("authorization", authHeader).
		Int("bodyLength", len(body)).
		Msg("LiveKit webhook received")

	// TODO: Implement proper webhook token validation using LiveKit SDK
	// For now, just parse the event from body

	// Parse the event type from body (simplified)
	var event struct {
		Event       string                   `json:"event"`
		Room        *livekit.Room            `json:"room"`
		Participant *livekit.ParticipantInfo `json:"participant"`
		EgressInfo  *livekit.EgressInfo      `json:"egressInfo"`
	}

	if err := c.BodyParser(&event); err != nil {
		log.Error().Err(err).Msg("Failed to parse webhook event")
		return utils.BadRequestResponse(c, "Invalid event data")
	}

	log.Info().
		Str("event", event.Event).
		Msg("Processing webhook event")

	ctx := c.Context()

	switch event.Event {
	case "room_started":
		if event.Room != nil {
			log.Info().
				Str("room", event.Room.Name).
				Msg("Room started")
		}

	case "room_finished":
		if event.Room != nil {
			log.Info().
				Str("room", event.Room.Name).
				Msg("Room finished")
			// Update room status to closed
			_ = h.roomService.UpdateRoomStatus(ctx, event.Room.Name, "close")
		}

	case "participant_joined":
		if event.Participant != nil && event.Room != nil {
			log.Info().
				Str("room", event.Room.Name).
				Str("identity", event.Participant.Identity).
				Str("name", event.Participant.Name).
				Msg("Participant joined")
			// Update user status
			_ = h.userService.UpdateUserStatus(ctx, event.Room.Name, event.Participant.Identity, "connect")
		}

	case "participant_left":
		if event.Participant != nil && event.Room != nil {
			log.Info().
				Str("room", event.Room.Name).
				Str("identity", event.Participant.Identity).
				Msg("Participant left")
			// Update user status
			_ = h.userService.UpdateUserStatus(ctx, event.Room.Name, event.Participant.Identity, "disconnect")
		}

	case "track_published":
		if event.Participant != nil && event.Room != nil {
			log.Info().
				Str("room", event.Room.Name).
				Str("identity", event.Participant.Identity).
				Msg("Track published")
		}

	case "track_unpublished":
		if event.Participant != nil && event.Room != nil {
			log.Info().
				Str("room", event.Room.Name).
				Str("identity", event.Participant.Identity).
				Msg("Track unpublished")
		}

	case "egress_started":
		if event.EgressInfo != nil {
			log.Info().
				Str("egressId", event.EgressInfo.EgressId).
				Str("roomName", event.EgressInfo.RoomName).
				Msg("Egress started")
		}

	case "egress_updated":
		if event.EgressInfo != nil {
			log.Info().
				Str("egressId", event.EgressInfo.EgressId).
				Str("status", event.EgressInfo.Status.String()).
				Msg("Egress updated")
		}

	case "egress_ended":
		if event.EgressInfo != nil {
			log.Info().
				Str("egressId", event.EgressInfo.EgressId).
				Str("status", event.EgressInfo.Status.String()).
				Msg("Egress ended")

			// Update record status in database
			if event.EgressInfo.Status == livekit.EgressStatus_EGRESS_COMPLETE {
				// Get file info from egress results
				var fileName, filePath string
				var fileSize, duration int

				if len(event.EgressInfo.GetFileResults()) > 0 {
					result := event.EgressInfo.GetFileResults()[0]
					fileName = result.Filename
					filePath = result.Filename
					fileSize = int(result.Size)
					duration = int(result.Duration / 1000000000) // Convert nanoseconds to seconds
				}

				_ = h.recordService.UpdateRecordByEgressID(
					ctx,
					event.EgressInfo.EgressId,
					fileName,
					filePath,
					"complete",
					fileSize,
					duration,
				)

				// Update room record status
				_ = h.roomService.UpdateRecordStatus(ctx, event.EgressInfo.RoomName, 0)
			} else if event.EgressInfo.Status == livekit.EgressStatus_EGRESS_FAILED {
				_ = h.recordService.UpdateRecordByEgressID(
					ctx,
					event.EgressInfo.EgressId,
					"",
					"",
					"failed",
					0,
					0,
				)
			}
		}

	default:
		log.Warn().
			Str("event", event.Event).
			Msg("Unknown webhook event")
	}

	return c.SendStatus(fiber.StatusOK)
}

// HandleGenericWebhook handles generic webhook events
// POST /webhook/generic
func (h *WebhookHandler) HandleGenericWebhook(c *fiber.Ctx) error {
	// Log webhook received
	log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Int("bodyLength", len(c.Body())).
		Msg("Generic webhook received")

	// Just acknowledge receipt
	return utils.SuccessResponse(c, fiber.Map{
		"received": true,
	})
}
