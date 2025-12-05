package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/models"
	"api-gateway-go/internal/repository"
	"api-gateway-go/pkg/utils"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrRoomClosed   = errors.New("room is closed")
	ErrRoomExpired  = errors.New("room has expired")
)

// CreateRoomOptions holds options for creating a room
type CreateRoomOptions struct {
	Service               int
	RoomType              string
	AutoRecord            int
	RecordType            string
	EncodingOptionsPreset string
	ChatEnabled           int
	WebSocketURL          string
	UserAgent             string
	DaysExpired           int
}

// RoomService handles room business logic
type RoomService struct {
	roomRepo   *repository.RoomRepository
	livekitMgr *config.LiveKitManager
	cfg        *config.Config
}

// NewRoomService creates a new RoomService
func NewRoomService(roomRepo *repository.RoomRepository, livekitMgr *config.LiveKitManager, cfg *config.Config) *RoomService {
	return &RoomService{
		roomRepo:   roomRepo,
		livekitMgr: livekitMgr,
		cfg:        cfg,
	}
}

// CreateRoom creates a new room
func (s *RoomService) CreateRoom(ctx context.Context, opts CreateRoomOptions) (*models.RoomConference, error) {
	// Generate room name
	roomName := utils.GenerateRoomName()
	now := utils.FormatDateTimeNow()

	// Calculate expiration
	daysExpired := opts.DaysExpired
	if daysExpired == 0 {
		daysExpired = s.cfg.RoomDayDefaultTimeout
	}
	expiredAt := utils.AddDays(time.Now(), daysExpired)

	// Create room in database
	params := repository.CreateRoomParams{
		Room:                  roomName,
		Service:               opts.Service,
		Status:                "open",
		RoomType:              opts.RoomType,
		AutoRecord:            opts.AutoRecord,
		RecordType:            opts.RecordType,
		EncodingOptionsPreset: opts.EncodingOptionsPreset,
		ChatEnabled:           opts.ChatEnabled,
		WebSocketURL:          opts.WebSocketURL,
		UserAgent:             opts.UserAgent,
		DtmCreated:            now,
		DtmUpdated:            now,
		DtmExpired:            expiredAt,
	}

	roomID, err := s.roomRepo.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	// Create room in LiveKit if available
	if s.livekitMgr != nil && s.livekitMgr.RoomClient() != nil {
		_, err := s.livekitMgr.CreateRoom(ctx, roomName, 0, 0)
		if err != nil {
			// Log error but don't fail - room can be created on-demand
		}
	}

	// Get created room
	room, err := s.roomRepo.GetByID(ctx, uint(roomID))
	if err != nil {
		return nil, err
	}

	return room, nil
}

// GetRoomDetail gets room details by name
func (s *RoomService) GetRoomDetail(ctx context.Context, roomName string) (*models.RoomConference, error) {
	room, err := s.roomRepo.GetByRoom(ctx, roomName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}
	return room, nil
}

// GetRoomByID gets room by ID
func (s *RoomService) GetRoomByID(ctx context.Context, id uint) (*models.RoomConference, error) {
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}
	return room, nil
}

// GetRoomsByStatus gets all rooms with a specific status
func (s *RoomService) GetRoomsByStatus(ctx context.Context, status string) ([]models.RoomConference, error) {
	return s.roomRepo.GetByStatus(ctx, status)
}

// CloseRoom closes a room
func (s *RoomService) CloseRoom(ctx context.Context, roomName string) error {
	// Close in database
	_, err := s.roomRepo.CloseRoom(ctx, roomName)
	if err != nil {
		return err
	}

	// Delete from LiveKit if available
	if s.livekitMgr != nil && s.livekitMgr.RoomClient() != nil {
		err := s.livekitMgr.DeleteRoom(ctx, roomName)
		if err != nil {
			// Log error but don't fail
		}
	}

	return nil
}

// CloseAllRooms closes all open rooms
func (s *RoomService) CloseAllRooms(ctx context.Context) error {
	return s.roomRepo.CloseAllRooms(ctx)
}

// UpdateRoomStatus updates room status
func (s *RoomService) UpdateRoomStatus(ctx context.Context, roomName, status string) error {
	return s.roomRepo.UpdateStatus(ctx, roomName, status)
}

// UpdateRoomType updates room type and related fields
func (s *RoomService) UpdateRoomType(ctx context.Context, roomName string, opts repository.UpdateRoomTypeParams) error {
	opts.Room = roomName
	opts.DtmUpdated = utils.FormatDateTimeNow()
	return s.roomRepo.UpdateRoomType(ctx, opts)
}

// UpdateRecordStatus updates room record status
func (s *RoomService) UpdateRecordStatus(ctx context.Context, roomName string, status int) error {
	return s.roomRepo.UpdateRecordStatus(ctx, roomName, status)
}

// CheckRoomExpired checks if a room has expired
func (s *RoomService) CheckRoomExpired(ctx context.Context, roomName string) (bool, error) {
	return s.roomRepo.CheckRoomExpired(ctx, roomName)
}

// AutoRoomExpiredClose closes all expired rooms
func (s *RoomService) AutoRoomExpiredClose(ctx context.Context) error {
	rooms, err := s.roomRepo.GetExpiredRooms(ctx)
	if err != nil {
		return err
	}

	for _, room := range rooms {
		if room.Room.Valid {
			err := s.CloseRoom(ctx, room.Room.String)
			if err != nil {
				// Log but continue
			}
		}
	}

	return nil
}

// DeleteRoom deletes a room
func (s *RoomService) DeleteRoom(ctx context.Context, roomName string) error {
	// Delete from LiveKit first
	if s.livekitMgr != nil && s.livekitMgr.RoomClient() != nil {
		_ = s.livekitMgr.DeleteRoom(ctx, roomName)
	}

	return s.roomRepo.Delete(ctx, roomName)
}

// UpdateExpired updates room expiration time
func (s *RoomService) UpdateExpired(ctx context.Context, roomName, expiredAt string) error {
	return s.roomRepo.UpdateExpired(ctx, roomName, expiredAt)
}

// UpdateRecordID updates the record ID
func (s *RoomService) UpdateRecordID(ctx context.Context, roomName, recordID string) error {
	return s.roomRepo.UpdateRecordID(ctx, roomName, recordID)
}

// GetServiceID gets the service ID for a room
func (s *RoomService) GetServiceID(ctx context.Context, roomName string) (int, error) {
	return s.roomRepo.GetServiceID(ctx, roomName)
}

// UpdateTime updates room dtmUpdated
func (s *RoomService) UpdateTime(ctx context.Context, roomName string) error {
	return s.roomRepo.UpdateTime(ctx, roomName)
}

// UpdateMessageUnread updates message unread status
func (s *RoomService) UpdateMessageUnread(ctx context.Context, roomName string, messageUnread int) error {
	return s.roomRepo.UpdateMessageUnread(ctx, roomName, messageUnread)
}

// GetCountUnreadMessage gets count of unread messages
func (s *RoomService) GetCountUnreadMessage(ctx context.Context, service int) (int, error) {
	return s.roomRepo.GetCountUnreadMessage(ctx, service)
}

// UpdateRoomStartedFinished updates room started/finished timestamps
func (s *RoomService) UpdateRoomStartedFinished(ctx context.Context, roomName string, started, finished bool) error {
	return s.roomRepo.UpdateRoomStartedFinished(ctx, roomName, started, finished)
}

// UpdateStartStopRecord updates start/stop record timestamps
func (s *RoomService) UpdateStartStopRecord(ctx context.Context, roomName string, isStart bool) error {
	return s.roomRepo.UpdateStartStopRecord(ctx, roomName, isStart)
}
