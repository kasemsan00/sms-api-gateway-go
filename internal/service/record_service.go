package service

import (
	"context"
	"database/sql"
	"errors"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/models"
	"api-gateway-go/internal/repository"

	"github.com/livekit/protocol/livekit"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEgressLimit    = errors.New("egress limit reached")
)

// StartRecordOptions holds options for starting a recording
type StartRecordOptions struct {
	Room       string
	RecordType string
	FilePath   string
}

// RecordService handles recording business logic
type RecordService struct {
	recordRepo *repository.RecordRepository
	roomRepo   *repository.RoomRepository
	livekitMgr *config.LiveKitManager
	cfg        *config.Config
}

// NewRecordService creates a new RecordService
func NewRecordService(recordRepo *repository.RecordRepository, roomRepo *repository.RoomRepository, livekitMgr *config.LiveKitManager, cfg *config.Config) *RecordService {
	return &RecordService{
		recordRepo: recordRepo,
		roomRepo:   roomRepo,
		livekitMgr: livekitMgr,
		cfg:        cfg,
	}
}

// StartRecord starts a recording
func (s *RecordService) StartRecord(ctx context.Context, opts StartRecordOptions) (*livekit.EgressInfo, error) {
	if s.livekitMgr == nil || s.livekitMgr.EgressClient() == nil {
		return nil, errors.New("LiveKit not configured")
	}

	// Check egress limit
	available, err := s.recordRepo.CheckEgressAvailable(ctx, s.cfg.EgressLimit)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrEgressLimit
	}

	// Start room composite egress
	filePath := opts.FilePath
	if filePath == "" {
		filePath = s.cfg.RecordPath + "/" + opts.Room + ".mp4"
	}

	output := &livekit.EncodedFileOutput{
		Filepath: filePath,
	}

	info, err := s.livekitMgr.StartRoomCompositeEgress(ctx, opts.Room, output)
	if err != nil {
		return nil, err
	}

	// Update room record status
	_ = s.roomRepo.UpdateRecordStatus(ctx, opts.Room, 1)
	_ = s.roomRepo.UpdateRecordID(ctx, opts.Room, info.EgressId)

	// Create record entry in database
	_, _ = s.recordRepo.Create(ctx, repository.CreateRecordParams{
		EgressID:   info.EgressId,
		Room:       opts.Room,
		FilePath:   filePath,
		RecordType: opts.RecordType,
		Status:     "recording",
	})

	return info, nil
}

// StopRecord stops a recording
func (s *RecordService) StopRecord(ctx context.Context, recordID string) (*livekit.EgressInfo, error) {
	if s.livekitMgr == nil || s.livekitMgr.EgressClient() == nil {
		return nil, errors.New("LiveKit not configured")
	}

	info, err := s.livekitMgr.StopEgress(ctx, recordID)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// ListEgress lists active egresses
func (s *RecordService) ListEgress(ctx context.Context, room string) ([]*livekit.EgressInfo, error) {
	if s.livekitMgr == nil || s.livekitMgr.EgressClient() == nil {
		return nil, errors.New("LiveKit not configured")
	}

	return s.livekitMgr.ListEgress(ctx, room)
}

// StopAllActiveRecords stops all active recordings
func (s *RecordService) StopAllActiveRecords(ctx context.Context) ([]string, error) {
	if s.livekitMgr == nil || s.livekitMgr.EgressClient() == nil {
		return nil, errors.New("LiveKit not configured")
	}

	// List all active egresses
	egresses, err := s.livekitMgr.ListEgress(ctx, "")
	if err != nil {
		return nil, err
	}

	var stopped []string
	for _, egress := range egresses {
		if egress.Status == livekit.EgressStatus_EGRESS_ACTIVE {
			_, err := s.livekitMgr.StopEgress(ctx, egress.EgressId)
			if err == nil {
				stopped = append(stopped, egress.EgressId)
			}
		}
	}

	return stopped, nil
}

// GetRecordDetail gets record detail by ID
func (s *RecordService) GetRecordDetail(ctx context.Context, id int) (*models.RecordMedia, error) {
	record, err := s.recordRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	if record == nil {
		return nil, ErrRecordNotFound
	}
	return record, nil
}

// GetRecordByRoom gets records by room
func (s *RecordService) GetRecordByRoom(ctx context.Context, room string) ([]models.RecordMedia, error) {
	return s.recordRepo.GetByRoom(ctx, room)
}

// GetFileHistory gets file history
func (s *RecordService) GetFileHistory(ctx context.Context, room string) ([]models.RecordMedia, error) {
	return s.recordRepo.GetFileHistory(ctx, room)
}

// CheckEgressAvailable checks if egress is available
func (s *RecordService) CheckEgressAvailable(ctx context.Context) (bool, error) {
	return s.recordRepo.CheckEgressAvailable(ctx, s.cfg.EgressLimit)
}

// CreateRecord creates a new record entry
func (s *RecordService) CreateRecord(ctx context.Context, params repository.CreateRecordParams) (int64, error) {
	return s.recordRepo.Create(ctx, params)
}

// UpdateRecordByEgressID updates record by egress ID
func (s *RecordService) UpdateRecordByEgressID(ctx context.Context, egressID, fileName, filePath, status string, fileSize, duration int) error {
	return s.recordRepo.UpdateByEgressID(ctx, egressID, fileName, filePath, status, fileSize, duration)
}

// GetRecordQueue gets the record queue
func (s *RecordService) GetRecordQueue(ctx context.Context) ([]models.RecordMedia, error) {
	return s.recordRepo.GetRecordQueue(ctx)
}

// GetActiveRecordCount gets count of active recordings
func (s *RecordService) GetActiveRecordCount(ctx context.Context) (int, error) {
	return s.recordRepo.GetActiveRecordCount(ctx)
}
