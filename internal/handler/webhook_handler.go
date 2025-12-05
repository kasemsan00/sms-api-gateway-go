package handler

import (
	"context"
	"sync"
	"time"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/repository"
	"api-gateway-go/internal/service"
	"api-gateway-go/internal/socket"
	"api-gateway-go/pkg/logger"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/livekit/protocol/livekit"
	"github.com/rs/zerolog/log"
)

// TrackData stores track information for auto-recording
type TrackData struct {
	Room     string
	Video    *livekit.TrackInfo
	Audio    *livekit.TrackInfo
	Record   int
	Identity string
}

// WebhookHandler handles LiveKit webhook events
type WebhookHandler struct {
	roomService   *service.RoomService
	userService   *service.UserService
	recordService *service.RecordService
	recordRepo    *repository.RecordRepository
	livekitMgr    *config.LiveKitManager
	cfg           *config.Config
	socketHub     *socket.Hub

	// Track data management for auto-recording
	trackPublishedData map[string]*TrackData
	trackDataMu        sync.RWMutex
	participantTimers  map[string]*time.Timer
	timersMu           sync.RWMutex
	inactivityTimeout  time.Duration
	maxTrackDataSize   int
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
		roomService:        roomService,
		userService:        userService,
		recordService:      recordService,
		recordRepo:         recordRepo,
		livekitMgr:         livekitMgr,
		cfg:                cfg,
		trackPublishedData: make(map[string]*TrackData),
		participantTimers:  make(map[string]*time.Timer),
		inactivityTimeout:  30 * time.Minute,
		maxTrackDataSize:   1000,
	}
}

// SetSocketHub sets the socket hub for real-time updates
func (h *WebhookHandler) SetSocketHub(hub *socket.Hub) {
	h.socketHub = hub
}

// HandleLiveKitWebhook handles LiveKit webhook events
// POST /webhook/livekit
func (h *WebhookHandler) HandleLiveKitWebhook(c *fiber.Ctx) error {
	// Get authorization header for validation
	authHeader := c.Get("Authorization")
	body := c.Body()

	// Log webhook received
	log.Info().
		Str("authorization", authHeader).
		Int("bodyLength", len(body)).
		Msg("LiveKit webhook received")

	// Parse webhook event
	var event struct {
		Event       string                   `json:"event"`
		Room        *livekit.Room            `json:"room"`
		Participant *livekit.ParticipantInfo `json:"participant"`
		Track       *livekit.TrackInfo       `json:"track"`
		EgressInfo  *livekit.EgressInfo      `json:"egressInfo"`
	}

	if err := c.BodyParser(&event); err != nil {
		log.Error().Err(err).Msg("Failed to parse webhook event")
		return utils.BadRequestResponse(c, "Invalid event data")
	}

	log.Info().
		Str("event", event.Event).
		Msg("Processing webhook event")

	// Process the event asynchronously
	go h.processEvent(event.Event, event.Room, event.Participant, event.Track, event.EgressInfo)

	return c.SendStatus(fiber.StatusOK)
}

// processEvent processes webhook events
func (h *WebhookHandler) processEvent(eventType string, room *livekit.Room, participant *livekit.ParticipantInfo, track *livekit.TrackInfo, egressInfo *livekit.EgressInfo) {
	ctx := context.Background()

	// Handle participant activity for inactivity timers
	if participant != nil {
		switch eventType {
		case "participant_joined":
			h.startInactivityTimer(participant.Sid)
		case "participant_left":
			h.clearInactivityTimer(participant.Sid)
		case "track_published", "track_subscribed":
			h.resetInactivityTimer(participant.Sid)
		}
	}

	switch eventType {
	case "room_started":
		if room != nil {
			log.Info().Str("room", room.Name).Msg("Room started")
			// Update room start time
			_ = h.roomService.UpdateStartStopRecord(ctx, room.Name, true)
		}

	case "room_finished":
		if room != nil {
			log.Info().Str("room", room.Name).Msg("Room finished")
			// Update room finish time
			_ = h.roomService.UpdateStartStopRecord(ctx, room.Name, false)
			// Stop any active recording
			if _, err := h.recordService.StopRecord(ctx, room.Name); err != nil {
				log.Debug().Err(err).Msg("No active recording to stop")
			}
			// Update room status to closed
			_ = h.roomService.UpdateRoomStatus(ctx, room.Name, "close")
		}

	case "participant_joined":
		if participant != nil && room != nil {
			log.Info().
				Str("room", room.Name).
				Str("identity", participant.Identity).
				Str("name", participant.Name).
				Msg("Participant joined")
			// Update user status
			_ = h.userService.UpdateUserStatus(ctx, room.Name, participant.Identity, "connect")

			// Check for auto-recording
			h.checkAutoRecord(ctx, room.Name, participant)
		}

	case "participant_left":
		if participant != nil && room != nil {
			log.Info().
				Str("room", room.Name).
				Str("identity", participant.Identity).
				Msg("Participant left")
			// Update user status
			_ = h.userService.UpdateUserStatus(ctx, room.Name, participant.Identity, "disconnect")

			// Clean up track data
			h.trackDataMu.Lock()
			delete(h.trackPublishedData, participant.Identity)
			h.trackDataMu.Unlock()
		}

	case "track_published":
		if participant != nil && room != nil && track != nil {
			log.Info().
				Str("room", room.Name).
				Str("identity", participant.Identity).
				Str("trackSource", track.Source.String()).
				Msg("Track published")

			h.handleTrackPublished(ctx, room.Name, participant, track)
		}

	case "egress_started":
		if egressInfo != nil {
			log.Info().
				Str("egressId", egressInfo.EgressId).
				Str("roomName", egressInfo.RoomName).
				Msg("Egress started")

			roomName := getRoomNameFromEgress(egressInfo)
			if roomName != "" {
				_ = h.roomService.UpdateRecordID(ctx, roomName, egressInfo.EgressId)
				_ = h.roomService.UpdateRecordStatus(ctx, roomName, 1)

				// Broadcast to socket
				if h.socketHub != nil {
					h.socketHub.BroadcastToRoom("/"+roomName, roomName, "room-record", map[string]interface{}{
						"egressId":       egressInfo.EgressId,
						"status":         "startRecord",
						"dtmStartRecord": time.Now().Format("2006-01-02 15:04:05"),
					})
				}
			}
		}

	case "egress_updated":
		if egressInfo != nil {
			log.Info().
				Str("egressId", egressInfo.EgressId).
				Str("status", egressInfo.Status.String()).
				Msg("Egress updated")
		}

	case "egress_ended":
		if egressInfo != nil {
			log.Info().
				Str("egressId", egressInfo.EgressId).
				Str("status", egressInfo.Status.String()).
				Msg("Egress ended")

			roomName := getRoomNameFromEgress(egressInfo)
			if roomName != "" {
				_ = h.roomService.UpdateRecordID(ctx, roomName, "")
				_ = h.roomService.UpdateRecordStatus(ctx, roomName, 2)

				// Update record info
				if egressInfo.Status == livekit.EgressStatus_EGRESS_COMPLETE {
					var fileName, filePath string
					var fileSize, duration int

					if len(egressInfo.GetFileResults()) > 0 {
						result := egressInfo.GetFileResults()[0]
						fileName = result.Filename
						filePath = result.Filename
						fileSize = int(result.Size)
						duration = int(result.Duration / 1000000000) // nanoseconds to seconds
					}

					_ = h.recordService.UpdateRecordByEgressID(ctx, egressInfo.EgressId, fileName, filePath, "complete", fileSize, duration)
				} else if egressInfo.Status == livekit.EgressStatus_EGRESS_FAILED {
					_ = h.recordService.UpdateRecordByEgressID(ctx, egressInfo.EgressId, "", "", "failed", 0, 0)
				}

				// Broadcast to socket
				if h.socketHub != nil {
					h.socketHub.BroadcastToRoom("/"+roomName, roomName, "room-record", map[string]interface{}{
						"egressId":      egressInfo.EgressId,
						"status":        "stopRecord",
						"dtmStopRecord": time.Now().Format("2006-01-02 15:04:05"),
					})
				}
			}
		}

	default:
		log.Warn().Str("event", eventType).Msg("Unknown webhook event")
	}
}

// checkAutoRecord checks if auto-recording should start
func (h *WebhookHandler) checkAutoRecord(ctx context.Context, roomName string, participant *livekit.ParticipantInfo) {
	if participant.Metadata == "" {
		return
	}

	roomDetail, err := h.roomService.GetRoomDetail(ctx, roomName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get room detail for auto-record check")
		return
	}

	// Check if auto-record is enabled
	if !roomDetail.AutoRecord.Valid || roomDetail.AutoRecord.Int32 != 1 {
		return
	}

	// Check if already recording
	if roomDetail.RecordStatus.Valid && roomDetail.RecordStatus.Int32 == 1 {
		return
	}

	recordType := ""
	if roomDetail.RecordType.Valid {
		recordType = roomDetail.RecordType.String
	}

	// Initialize track data for TrackComposite recording
	if recordType == "TrackComposite" && !isEgressParticipant(participant.Identity) {
		h.trackDataMu.Lock()
		h.trackPublishedData[participant.Identity] = &TrackData{
			Room:     roomName,
			Identity: participant.Identity,
			Record:   0,
		}
		h.trackDataMu.Unlock()
	}

	// Start room composite recording
	if recordType == "RoomCompositeVideoAudio" || recordType == "RoomCompositeAudio" {
		log.Info().Str("room", roomName).Str("recordType", recordType).Msg("Starting auto-record")

		opts := service.StartRecordOptions{
			Room:       roomName,
			RecordType: recordType,
			FilePath:   h.cfg.RecordPath,
		}

		result, err := h.recordService.StartRecord(ctx, opts)
		if err != nil {
			log.Error().Err(err).Msg("Failed to start auto-record")
			return
		}

		// Broadcast to socket
		if h.socketHub != nil && result != nil {
			h.socketHub.BroadcastToRoom("/"+roomName, roomName, "room-record", map[string]interface{}{
				"egressId": result.EgressId,
				"status":   "startRecord",
			})
		}
	}
}

// handleTrackPublished handles track published events for TrackComposite recording
func (h *WebhookHandler) handleTrackPublished(ctx context.Context, roomName string, participant *livekit.ParticipantInfo, track *livekit.TrackInfo) {
	h.trackDataMu.Lock()
	trackData, exists := h.trackPublishedData[participant.Identity]
	if !exists {
		h.trackDataMu.Unlock()
		return
	}

	// Update track data based on source
	switch track.Source {
	case livekit.TrackSource_CAMERA:
		trackData.Video = track
	case livekit.TrackSource_MICROPHONE:
		trackData.Audio = track
	}

	// Check if we can start TrackComposite recording
	if trackData.Record == 0 && trackData.Video != nil && trackData.Audio != nil {
		trackData.Record = 1
		h.trackDataMu.Unlock()

		log.Info().Str("room", roomName).Str("identity", participant.Identity).Msg("Starting TrackComposite auto-record")

		// Start track composite record would be implemented here
		// For now, this is a placeholder
	} else {
		h.trackDataMu.Unlock()
	}
}

// Timer management functions
func (h *WebhookHandler) startInactivityTimer(participantSid string) {
	h.timersMu.Lock()
	defer h.timersMu.Unlock()

	// Clear existing timer
	if timer, exists := h.participantTimers[participantSid]; exists {
		timer.Stop()
	}

	// Create new timer
	h.participantTimers[participantSid] = time.AfterFunc(h.inactivityTimeout, func() {
		logger.Info("Participant %s timed out due to inactivity", participantSid)
		h.timersMu.Lock()
		delete(h.participantTimers, participantSid)
		h.timersMu.Unlock()
	})
}

func (h *WebhookHandler) resetInactivityTimer(participantSid string) {
	h.startInactivityTimer(participantSid)
}

func (h *WebhookHandler) clearInactivityTimer(participantSid string) {
	h.timersMu.Lock()
	defer h.timersMu.Unlock()

	if timer, exists := h.participantTimers[participantSid]; exists {
		timer.Stop()
		delete(h.participantTimers, participantSid)
	}
}

// Cleanup cleans up all resources
func (h *WebhookHandler) Cleanup() {
	// Clear all timers
	h.timersMu.Lock()
	for _, timer := range h.participantTimers {
		timer.Stop()
	}
	h.participantTimers = make(map[string]*time.Timer)
	h.timersMu.Unlock()

	// Clear track data
	h.trackDataMu.Lock()
	h.trackPublishedData = make(map[string]*TrackData)
	h.trackDataMu.Unlock()

	logger.Info("WebhookHandler cleaned up")
}

// GetStatus returns the status of the webhook handler
func (h *WebhookHandler) GetStatus() map[string]interface{} {
	h.timersMu.RLock()
	timerCount := len(h.participantTimers)
	h.timersMu.RUnlock()

	h.trackDataMu.RLock()
	trackDataCount := len(h.trackPublishedData)
	h.trackDataMu.RUnlock()

	return map[string]interface{}{
		"activeTimers":     timerCount,
		"trackDataEntries": trackDataCount,
		"maxTrackDataSize": h.maxTrackDataSize,
		"healthy":          trackDataCount < h.maxTrackDataSize,
	}
}

// HandleGenericWebhook handles generic webhook events
// POST /webhook/generic
func (h *WebhookHandler) HandleGenericWebhook(c *fiber.Ctx) error {
	log.Info().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Int("bodyLength", len(c.Body())).
		Msg("Generic webhook received")

	return utils.SuccessResponse(c, fiber.Map{
		"received": true,
	})
}

// Helper functions
func getRoomNameFromEgress(info *livekit.EgressInfo) string {
	if info.RoomName != "" {
		return info.RoomName
	}
	if rc := info.GetRoomComposite(); rc != nil {
		return rc.RoomName
	}
	if tc := info.GetTrackComposite(); tc != nil {
		return tc.RoomName
	}
	return ""
}

func isEgressParticipant(identity string) bool {
	return len(identity) > 3 && identity[:3] == "EG_"
}
