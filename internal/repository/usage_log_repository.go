package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// UsageLogRepository handles usage log database operations
type UsageLogRepository struct {
	db *sqlx.DB
}

// NewUsageLogRepository creates a new UsageLogRepository
func NewUsageLogRepository(db *sqlx.DB) *UsageLogRepository {
	return &UsageLogRepository{db: db}
}

// AddStatusLogParams holds parameters for adding a status log
type AddStatusLogParams struct {
	LinkID    string
	LinkType  string
	Mobile    string
	Room      string
	Identity  string
	UserName  string
	UserType  string
	Status    string
	Latitude  float64
	Longitude float64
	UserAgent string
	Data      string
}

// AddStatusLog adds a status log
func (r *UsageLogRepository) AddStatusLog(ctx context.Context, params AddStatusLogParams) (int64, error) {
	dtmCreated := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO usage_status_log
		(linkID, linkType, mobile, room, identity, userName, userType, status, latitude, longitude, userAgent, data, dtmCreated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		params.LinkID,
		params.LinkType,
		params.Mobile,
		params.Room,
		params.Identity,
		params.UserName,
		params.UserType,
		params.Status,
		params.Latitude,
		params.Longitude,
		params.UserAgent,
		params.Data,
		dtmCreated,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to add status log: %w", err)
	}

	return result.LastInsertId()
}

// AddDataLog adds a data log
func (r *UsageLogRepository) AddDataLog(ctx context.Context, data interface{}) (int64, error) {
	dtmCreated := time.Now().Format("2006-01-02 15:04:05")

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal data: %w", err)
	}

	query := `INSERT INTO data_log (data, dtmCreated) VALUES (?, ?)`

	result, err := r.db.ExecContext(ctx, query, string(dataJSON), dtmCreated)
	if err != nil {
		return 0, fmt.Errorf("failed to add data log: %w", err)
	}

	return result.LastInsertId()
}

// GetStatusLogsByLinkID gets status logs by linkID
func (r *UsageLogRepository) GetStatusLogsByLinkID(ctx context.Context, linkID string) ([]models.UsageStatusLog, error) {
	var logs []models.UsageStatusLog
	query := `SELECT * FROM usage_status_log WHERE linkID = ? ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &logs, query, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get status logs: %w", err)
	}

	return logs, nil
}

// GetStatusLogsByRoom gets status logs by room
func (r *UsageLogRepository) GetStatusLogsByRoom(ctx context.Context, room string) ([]models.UsageStatusLog, error) {
	var logs []models.UsageStatusLog
	query := `SELECT * FROM usage_status_log WHERE room = ? ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &logs, query, room)
	if err != nil {
		return nil, fmt.Errorf("failed to get status logs by room: %w", err)
	}

	return logs, nil
}

// GetAgent gets agent username by room
func (r *UsageLogRepository) GetAgent(ctx context.Context, room string) (string, error) {
	var userName string
	query := `SELECT userName FROM usage_status_log WHERE room = ? AND userType = 'admin' ORDER BY dtmCreated DESC LIMIT 1`

	err := r.db.GetContext(ctx, &userName, query, room)
	if err != nil {
		return "", nil // Return empty if not found
	}

	return userName, nil
}

// GetRoomStatus gets room status from usage logs
func (r *UsageLogRepository) GetRoomStatus(ctx context.Context, room string) (string, error) {
	var status string
	query := `SELECT status FROM usage_status_log WHERE room = ? ORDER BY dtmCreated DESC LIMIT 1`

	err := r.db.GetContext(ctx, &status, query, room)
	if err != nil {
		return "unknown", nil // Return unknown if not found
	}

	return status, nil
}

// GetCRMLinkStatusLog gets CRM link status log
func (r *UsageLogRepository) GetCRMLinkStatusLog(ctx context.Context, linkID string) ([]models.UsageStatusLog, error) {
	var logs []models.UsageStatusLog
	query := `SELECT * FROM usage_status_log WHERE linkID = ? ORDER BY dtmCreated ASC`

	err := r.db.SelectContext(ctx, &logs, query, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get CRM link status log: %w", err)
	}

	return logs, nil
}

// GetDataLogs gets data logs with pagination
func (r *UsageLogRepository) GetDataLogs(ctx context.Context, limit, offset int) ([]models.DataLog, error) {
	var logs []models.DataLog
	query := `SELECT * FROM data_log ORDER BY dtmCreated DESC LIMIT ? OFFSET ?`

	err := r.db.SelectContext(ctx, &logs, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get data logs: %w", err)
	}

	return logs, nil
}

// DeleteOldLogs deletes logs older than specified days
func (r *UsageLogRepository) DeleteOldLogs(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02 15:04:05")

	query := `DELETE FROM usage_status_log WHERE dtmCreated < ?`
	_, err := r.db.ExecContext(ctx, query, cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to delete old status logs: %w", err)
	}

	query = `DELETE FROM data_log WHERE dtmCreated < ?`
	_, err = r.db.ExecContext(ctx, query, cutoffDate)
	if err != nil {
		return fmt.Errorf("failed to delete old data logs: %w", err)
	}

	return nil
}
