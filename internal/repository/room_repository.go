package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// RoomRepository handles room database operations
type RoomRepository struct {
	db *sqlx.DB
}

// NewRoomRepository creates a new RoomRepository
func NewRoomRepository(db *sqlx.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// CreateRoomParams holds parameters for creating a room
type CreateRoomParams struct {
	Room                  string
	Service               int
	Status                string
	RoomType              string
	AutoRecord            int
	RecordType            string
	EncodingOptionsPreset string
	ChatEnabled           int
	WebSocketURL          string
	UserAgent             string
	DtmCreated            string
	DtmUpdated            string
	DtmExpired            string
}

// Create creates a new room
func (r *RoomRepository) Create(ctx context.Context, params CreateRoomParams) (int64, error) {
	query := `INSERT INTO room_conference
		(status, roomType, room, service, recordId, autoRecord, recordType, encodingOptionsPreset, chatEnabled, webSocketURL, userAgent, dtmCreated, dtmUpdated, dtmExpired)
		VALUES (?, ?, ?, ?, '', ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		params.Status,
		params.RoomType,
		params.Room,
		params.Service,
		params.AutoRecord,
		params.RecordType,
		params.EncodingOptionsPreset,
		params.ChatEnabled,
		params.WebSocketURL,
		params.UserAgent,
		params.DtmCreated,
		params.DtmUpdated,
		params.DtmExpired,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create room: %w", err)
	}

	return result.LastInsertId()
}

// GetByRoom gets room details by room name
func (r *RoomRepository) GetByRoom(ctx context.Context, room string) (*models.RoomConference, error) {
	var roomConf models.RoomConference
	query := `SELECT id, status, room, roomType, recordId, autoRecord, recordType, encodingOptionsPreset,
		chatEnabled, messageUnread, userAgent, dtmCreated, dtmStartRecord, dtmStopRecord, webSocketURL
		FROM room_conference WHERE room = ? LIMIT 1`

	err := r.db.GetContext(ctx, &roomConf, query, room)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	return &roomConf, nil
}

// GetByID gets room details by ID
func (r *RoomRepository) GetByID(ctx context.Context, id uint) (*models.RoomConference, error) {
	var roomConf models.RoomConference
	query := `SELECT * FROM room_conference WHERE id = ? LIMIT 1`

	err := r.db.GetContext(ctx, &roomConf, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get room by ID: %w", err)
	}

	return &roomConf, nil
}

// GetByStatus gets all rooms by status
func (r *RoomRepository) GetByStatus(ctx context.Context, status string) ([]models.RoomConference, error) {
	var rooms []models.RoomConference
	query := `SELECT * FROM room_conference WHERE status = ?`

	err := r.db.SelectContext(ctx, &rooms, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms by status: %w", err)
	}

	return rooms, nil
}

// UpdateStatus updates room status
func (r *RoomRepository) UpdateStatus(ctx context.Context, room, status string) error {
	query := `UPDATE room_conference SET status = ?, messageUnread = 0 WHERE room = ?`
	_, err := r.db.ExecContext(ctx, query, status, room)
	if err != nil {
		return fmt.Errorf("failed to update room status: %w", err)
	}
	return nil
}

// CloseRoom closes a room by setting status to 'close'
func (r *RoomRepository) CloseRoom(ctx context.Context, room string) (int64, error) {
	query := `UPDATE room_conference SET status = 'close' WHERE room = ?`
	result, err := r.db.ExecContext(ctx, query, room)
	if err != nil {
		return 0, fmt.Errorf("failed to close room: %w", err)
	}
	return result.RowsAffected()
}

// CloseAllRooms closes all open rooms
func (r *RoomRepository) CloseAllRooms(ctx context.Context) error {
	query := `UPDATE room_conference SET status = 'close' WHERE status = 'open'`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to close all rooms: %w", err)
	}
	return nil
}

// UpdateExpired updates room expiration time
func (r *RoomRepository) UpdateExpired(ctx context.Context, room, dtmExpired string) error {
	query := `UPDATE room_conference SET dtmExpired = ? WHERE room = ?`
	_, err := r.db.ExecContext(ctx, query, dtmExpired, room)
	if err != nil {
		return fmt.Errorf("failed to update room expiration: %w", err)
	}
	return nil
}

// UpdateRoomTypeParams holds parameters for updating room type
type UpdateRoomTypeParams struct {
	Room         string
	RoomType     string
	AutoRecord   int
	ChatEnabled  int
	WebSocketURL string
	UserAgent    string
	DtmUpdated   string
}

// UpdateRoomType updates room type and related fields
func (r *RoomRepository) UpdateRoomType(ctx context.Context, params UpdateRoomTypeParams) error {
	query := `UPDATE room_conference SET
		roomType = ?,
		autoRecord = ?,
		chatEnabled = ?,
		webSocketURL = ?,
		userAgent = ?,
		dtmUpdated = ?
		WHERE room = ?`

	_, err := r.db.ExecContext(ctx, query,
		params.RoomType,
		params.AutoRecord,
		params.ChatEnabled,
		params.WebSocketURL,
		params.UserAgent,
		params.DtmUpdated,
		params.Room,
	)
	if err != nil {
		return fmt.Errorf("failed to update room type: %w", err)
	}
	return nil
}

// UpdateTime updates room dtmUpdated
func (r *RoomRepository) UpdateTime(ctx context.Context, room string) error {
	dtmCurrent := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE room_conference SET dtmUpdated = ? WHERE room = ?`
	_, err := r.db.ExecContext(ctx, query, dtmCurrent, room)
	if err != nil {
		return fmt.Errorf("failed to update room time: %w", err)
	}
	return nil
}

// Delete deletes a room from the database
func (r *RoomRepository) Delete(ctx context.Context, room string) error {
	query := `DELETE FROM room_conference WHERE room = ?`
	_, err := r.db.ExecContext(ctx, query, room)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}

// GetExpiredRooms gets rooms that are open but expired
func (r *RoomRepository) GetExpiredRooms(ctx context.Context) ([]models.RoomConference, error) {
	var rooms []models.RoomConference
	dtmCurrent := time.Now().Format("2006-01-02 15:04:05")
	query := `SELECT room, dtmExpired FROM room_conference WHERE status = 'open' AND dtmExpired < ?`

	err := r.db.SelectContext(ctx, &rooms, query, dtmCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired rooms: %w", err)
	}

	return rooms, nil
}

// CheckRoomExpired checks if a room is expired
func (r *RoomRepository) CheckRoomExpired(ctx context.Context, room string) (bool, error) {
	dtmCurrent := time.Now().Format("2006-01-02 15:04:05")
	query := `SELECT COUNT(*) FROM room_conference WHERE dtmExpired < ? AND room = ? LIMIT 1`

	var count int
	err := r.db.GetContext(ctx, &count, query, dtmCurrent, room)
	if err != nil {
		return false, fmt.Errorf("failed to check room expiration: %w", err)
	}

	return count > 0, nil
}

// UpdateRecordStatus updates room record status
func (r *RoomRepository) UpdateRecordStatus(ctx context.Context, room string, status int) error {
	var query string
	if room == "all" {
		query = `UPDATE room_conference SET recordStatus = 0, recordId = ''`
		_, err := r.db.ExecContext(ctx, query)
		return err
	}

	query = `UPDATE room_conference SET recordStatus = ? WHERE room = ?`
	_, err := r.db.ExecContext(ctx, query, status, room)
	if err != nil {
		return fmt.Errorf("failed to update record status: %w", err)
	}
	return nil
}

// UpdateRecordID updates room record ID
func (r *RoomRepository) UpdateRecordID(ctx context.Context, room, recordID string) error {
	query := `UPDATE room_conference SET recordId = ? WHERE room = ?`
	_, err := r.db.ExecContext(ctx, query, recordID, room)
	if err != nil {
		return fmt.Errorf("failed to update record ID: %w", err)
	}
	return nil
}

// UpdateMessageUnread updates message unread status
func (r *RoomRepository) UpdateMessageUnread(ctx context.Context, room string, messageUnread int) error {
	query := `UPDATE room_conference SET messageUnread = ? WHERE room = ?`
	_, err := r.db.ExecContext(ctx, query, messageUnread, room)
	if err != nil {
		return fmt.Errorf("failed to update message unread: %w", err)
	}
	return nil
}

// GetCountUnreadMessage gets count of unread messages for a service
func (r *RoomRepository) GetCountUnreadMessage(ctx context.Context, service int) (int, error) {
	query := `SELECT COUNT(*) FROM room_conference
		LEFT JOIN case_data ON case_data.roomId = room_conference.id
		WHERE messageUnread = 1 AND case_data.service = ?`

	var count int
	err := r.db.GetContext(ctx, &count, query, service)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread message count: %w", err)
	}
	return count, nil
}

// GetServiceID gets service ID by room name
func (r *RoomRepository) GetServiceID(ctx context.Context, room string) (int, error) {
	query := `SELECT case_data.service FROM case_data
		LEFT JOIN room_conference ON case_data.roomId = room_conference.id
		WHERE room_conference.room = ?`

	var service int
	err := r.db.GetContext(ctx, &service, query, room)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get service ID: %w", err)
	}
	return service, nil
}

// UpdateWebSocketURL updates room WebSocket URL
func (r *RoomRepository) UpdateWebSocketURL(ctx context.Context, room, url string) error {
	query := `UPDATE room_conference SET webSocketURL = ? WHERE room = ?`
	_, err := r.db.ExecContext(ctx, query, url, room)
	if err != nil {
		return fmt.Errorf("failed to update WebSocket URL: %w", err)
	}
	return nil
}

// UpdateRoomStartedFinished updates room started/finished timestamps
func (r *RoomRepository) UpdateRoomStartedFinished(ctx context.Context, room string, started, finished bool) error {
	var query string
	dtmCurrent := time.Now().Format("2006-01-02 15:04:05")

	if started {
		query = `UPDATE room_conference SET dtmRoomStarted = ? WHERE room = ?`
	} else if finished {
		query = `UPDATE room_conference SET dtmRoomFinished = ? WHERE room = ?`
	} else {
		return nil
	}

	_, err := r.db.ExecContext(ctx, query, dtmCurrent, room)
	if err != nil {
		return fmt.Errorf("failed to update room timestamps: %w", err)
	}
	return nil
}

// UpdateStartStopRecord updates start/stop record timestamps
func (r *RoomRepository) UpdateStartStopRecord(ctx context.Context, room string, isStart bool) error {
	dtmCurrent := time.Now().Format("2006-01-02 15:04:05")
	var query string

	if isStart {
		query = `UPDATE room_conference SET dtmStartRecord = ? WHERE room = ?`
	} else {
		query = `UPDATE room_conference SET dtmStopRecord = ? WHERE room = ?`
	}

	_, err := r.db.ExecContext(ctx, query, dtmCurrent, room)
	if err != nil {
		return fmt.Errorf("failed to update record timestamps: %w", err)
	}
	return nil
}
