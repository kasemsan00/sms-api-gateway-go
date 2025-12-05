package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// RecordRepository handles record database operations
type RecordRepository struct {
	db *sqlx.DB
}

// NewRecordRepository creates a new RecordRepository
func NewRecordRepository(db *sqlx.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

// CreateRecordParams holds parameters for creating a record
type CreateRecordParams struct {
	EgressID   string
	Room       string
	FileName   string
	FilePath   string
	FileSize   int
	Duration   int
	RecordType string
	Status     string
	HLS        string
	Encode     int
	Uploader   string
}

// Create creates a new record
func (r *RecordRepository) Create(ctx context.Context, params CreateRecordParams) (int64, error) {
	dtmCreated := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO record_media
		(egressId, room, fileName, filePath, fileSize, duration, recordType, status, hls, encode, uploader, dtmCreated, startRecord)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		params.EgressID,
		params.Room,
		params.FileName,
		params.FilePath,
		params.FileSize,
		params.Duration,
		params.RecordType,
		params.Status,
		params.HLS,
		params.Encode,
		params.Uploader,
		dtmCreated,
		dtmCreated,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create record: %w", err)
	}

	return result.LastInsertId()
}

// GetByID gets record by ID
func (r *RecordRepository) GetByID(ctx context.Context, id int) (*models.RecordMedia, error) {
	var record models.RecordMedia
	query := `SELECT * FROM record_media WHERE id = ?`

	err := r.db.GetContext(ctx, &record, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	return &record, nil
}

// GetByEgressID gets record by egressId
func (r *RecordRepository) GetByEgressID(ctx context.Context, egressID string) (*models.RecordMedia, error) {
	var record models.RecordMedia
	query := `SELECT * FROM record_media WHERE egressId = ?`

	err := r.db.GetContext(ctx, &record, query, egressID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get record by egress ID: %w", err)
	}

	return &record, nil
}

// GetByRoom gets records by room
func (r *RecordRepository) GetByRoom(ctx context.Context, room string) ([]models.RecordMedia, error) {
	var records []models.RecordMedia
	query := `SELECT * FROM record_media WHERE room = ? ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &records, query, room)
	if err != nil {
		return nil, fmt.Errorf("failed to get records by room: %w", err)
	}

	return records, nil
}

// Update updates a record
func (r *RecordRepository) Update(ctx context.Context, id int, fileName, filePath, status string, fileSize, duration int) error {
	dtmUpdated := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE record_media SET fileName = ?, filePath = ?, fileSize = ?, duration = ?, status = ?, dtmUpdated = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, fileName, filePath, fileSize, duration, status, dtmUpdated, id)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	return nil
}

// UpdateStatus updates record status
func (r *RecordRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	dtmUpdated := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE record_media SET status = ?, dtmUpdated = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, status, dtmUpdated, id)
	if err != nil {
		return fmt.Errorf("failed to update record status: %w", err)
	}

	return nil
}

// UpdateByEgressID updates record by egress ID
func (r *RecordRepository) UpdateByEgressID(ctx context.Context, egressID, fileName, filePath, status string, fileSize, duration int) error {
	dtmUpdated := time.Now().Format("2006-01-02 15:04:05")
	dtmCompleted := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE record_media SET fileName = ?, filePath = ?, fileSize = ?, duration = ?, status = ?, dtmUpdated = ?, dtmCompleted = ?, endRecord = ? WHERE egressId = ?`

	_, err := r.db.ExecContext(ctx, query, fileName, filePath, fileSize, duration, status, dtmUpdated, dtmCompleted, dtmCompleted, egressID)
	if err != nil {
		return fmt.Errorf("failed to update record by egress ID: %w", err)
	}

	return nil
}

// GetFileHistory gets file history
func (r *RecordRepository) GetFileHistory(ctx context.Context, room string) ([]models.RecordMedia, error) {
	var records []models.RecordMedia
	query := `SELECT * FROM record_media`
	args := []interface{}{}

	if room != "" {
		query += ` WHERE room = ?`
		args = append(args, room)
	}
	query += ` ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &records, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get file history: %w", err)
	}

	return records, nil
}

// GetRecordQueue gets record queue (records with status 'recording')
func (r *RecordRepository) GetRecordQueue(ctx context.Context) ([]models.RecordMedia, error) {
	var records []models.RecordMedia
	query := `SELECT * FROM record_media WHERE status = 'recording' ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &records, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get record queue: %w", err)
	}

	return records, nil
}

// CheckEgressAvailable checks if egress is available (count < limit)
func (r *RecordRepository) CheckEgressAvailable(ctx context.Context, limit int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM record_media WHERE status = 'recording'`

	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return false, fmt.Errorf("failed to check egress availability: %w", err)
	}

	return count < limit, nil
}

// GetActiveRecordCount gets count of active recordings
func (r *RecordRepository) GetActiveRecordCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM record_media WHERE status = 'recording'`

	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get active record count: %w", err)
	}

	return count, nil
}

// Delete deletes a record
func (r *RecordRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM record_media WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	return nil
}

// UpdateEncodeStatus updates encode status
func (r *RecordRepository) UpdateEncodeStatus(ctx context.Context, id, encode int) error {
	query := `UPDATE record_media SET encode = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, encode, id)
	if err != nil {
		return fmt.Errorf("failed to update encode status: %w", err)
	}

	return nil
}
