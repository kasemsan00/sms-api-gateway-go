package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// RadioRepository handles radio database operations
type RadioRepository struct {
	db *sqlx.DB
}

// NewRadioRepository creates a new RadioRepository
func NewRadioRepository(db *sqlx.DB) *RadioRepository {
	return &RadioRepository{db: db}
}

// GetDevices gets all radio devices
func (r *RadioRepository) GetDevices(ctx context.Context, limit int) ([]models.RadioDevices, error) {
	var devices []models.RadioDevices
	query := `SELECT * FROM radio_devices ORDER BY dtmUpdated DESC LIMIT ?`

	err := r.db.SelectContext(ctx, &devices, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get radio devices: %w", err)
	}

	return devices, nil
}

// GetDeviceByID gets device by ID
func (r *RadioRepository) GetDeviceByID(ctx context.Context, id int) (*models.RadioDevices, error) {
	var device models.RadioDevices
	query := `SELECT * FROM radio_devices WHERE id = ?`

	err := r.db.GetContext(ctx, &device, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get radio device: %w", err)
	}

	return &device, nil
}

// GetDeviceByDeviceID gets device by deviceId
func (r *RadioRepository) GetDeviceByDeviceID(ctx context.Context, deviceID string) (*models.RadioDevices, error) {
	var device models.RadioDevices
	query := `SELECT * FROM radio_devices WHERE deviceId = ?`

	err := r.db.GetContext(ctx, &device, query, deviceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get radio device by deviceId: %w", err)
	}

	return &device, nil
}

// CreateDevice creates a new radio device
func (r *RadioRepository) CreateDevice(ctx context.Context, deviceID, radioNo, channel, status string) (int64, error) {
	dtmCreated := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO radio_devices (deviceId, radioNo, channel, status, dtmCreated, dtmUpdated)
		VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, deviceID, radioNo, channel, status, dtmCreated, dtmCreated)
	if err != nil {
		return 0, fmt.Errorf("failed to create radio device: %w", err)
	}

	return result.LastInsertId()
}

// UpdateDevice updates a radio device
func (r *RadioRepository) UpdateDevice(ctx context.Context, id int, radioNo, channel, status string) error {
	dtmUpdated := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE radio_devices SET radioNo = ?, channel = ?, status = ?, dtmUpdated = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, radioNo, channel, status, dtmUpdated, id)
	if err != nil {
		return fmt.Errorf("failed to update radio device: %w", err)
	}

	return nil
}

// UpdateDeviceLocation updates device location
func (r *RadioRepository) UpdateDeviceLocation(ctx context.Context, deviceID string, latitude, longitude float64) error {
	dtmUpdated := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE radio_devices SET latitude = ?, longitude = ?, dtmUpdated = ? WHERE deviceId = ?`

	_, err := r.db.ExecContext(ctx, query, latitude, longitude, dtmUpdated, deviceID)
	if err != nil {
		return fmt.Errorf("failed to update device location: %w", err)
	}

	return nil
}

// DeleteDevice deletes a radio device
func (r *RadioRepository) DeleteDevice(ctx context.Context, id int) error {
	query := `DELETE FROM radio_devices WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete radio device: %w", err)
	}

	return nil
}

// GetLocations gets all radio locations
func (r *RadioRepository) GetLocations(ctx context.Context, limit int) ([]models.RadioLocations, error) {
	var locations []models.RadioLocations
	query := `SELECT * FROM radio_locations ORDER BY dtmCreated DESC LIMIT ?`

	err := r.db.SelectContext(ctx, &locations, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get radio locations: %w", err)
	}

	return locations, nil
}

// GetLocationByRadioNo gets location by radio number
func (r *RadioRepository) GetLocationByRadioNo(ctx context.Context, radioNo string) (*models.RadioLocations, error) {
	var location models.RadioLocations
	query := `SELECT * FROM radio_locations WHERE radioNo = ? ORDER BY dtmCreated DESC LIMIT 1`

	err := r.db.GetContext(ctx, &location, query, radioNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get radio location: %w", err)
	}

	return &location, nil
}

// GetLocationByID gets location by ID
func (r *RadioRepository) GetLocationByID(ctx context.Context, id int) (*models.RadioLocations, error) {
	var location models.RadioLocations
	query := `SELECT * FROM radio_locations WHERE id = ?`

	err := r.db.GetContext(ctx, &location, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get radio location by ID: %w", err)
	}

	return &location, nil
}

// CreateLocation creates a new radio location
func (r *RadioRepository) CreateLocation(ctx context.Context, radioNo string, latitude, longitude float64) (int64, error) {
	dtmCreated := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO radio_locations (radioNo, latitude, longitude, dtmCreated)
		VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, radioNo, latitude, longitude, dtmCreated)
	if err != nil {
		return 0, fmt.Errorf("failed to create radio location: %w", err)
	}

	return result.LastInsertId()
}
