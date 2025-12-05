package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// CarRepository handles car tracking database operations
type CarRepository struct {
	db *sqlx.DB
}

// NewCarRepository creates a new CarRepository
func NewCarRepository(db *sqlx.DB) *CarRepository {
	return &CarRepository{db: db}
}

// CreateTaskParams holds parameters for creating a car task
type CreateTaskParams struct {
	UID      string
	Status   string
	Mobile   string
	UserName string
	Room     string
}

// CreateTask creates a new car task
func (r *CarRepository) CreateTask(ctx context.Context, params CreateTaskParams) (int64, error) {
	dtmCreated := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO car_track (uid, status, mobile, userName, room, dtmCreated, dtmUpdated)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		params.UID,
		params.Status,
		params.Mobile,
		params.UserName,
		params.Room,
		dtmCreated,
		dtmCreated,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create car task: %w", err)
	}

	return result.LastInsertId()
}

// GetByID gets car task by ID
func (r *CarRepository) GetByID(ctx context.Context, id int) (*models.CarTrack, error) {
	var task models.CarTrack
	query := `SELECT * FROM car_track WHERE id = ?`

	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get car task: %w", err)
	}

	return &task, nil
}

// GetByUID gets car task by UID
func (r *CarRepository) GetByUID(ctx context.Context, uid string) (*models.CarTrack, error) {
	var task models.CarTrack
	query := `SELECT * FROM car_track WHERE uid = ?`

	err := r.db.GetContext(ctx, &task, query, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get car task by UID: %w", err)
	}

	return &task, nil
}

// GetByRoom gets car task by room
func (r *CarRepository) GetByRoom(ctx context.Context, room string) (*models.CarTrack, error) {
	var task models.CarTrack
	query := `SELECT * FROM car_track WHERE room = ? ORDER BY id DESC LIMIT 1`

	err := r.db.GetContext(ctx, &task, query, room)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get car task by room: %w", err)
	}

	return &task, nil
}

// UpdateTask updates a car task
func (r *CarRepository) UpdateTask(ctx context.Context, id int, status string) error {
	dtmUpdated := time.Now().Format("2006-01-02 15:04:05")
	var query string

	switch status {
	case "started":
		query = `UPDATE car_track SET status = ?, dtmStarted = ?, dtmUpdated = ? WHERE id = ?`
	case "arrived":
		query = `UPDATE car_track SET status = ?, dtmArrived = ?, dtmUpdated = ? WHERE id = ?`
	case "canceled":
		query = `UPDATE car_track SET status = ?, dtmCanceled = ?, dtmUpdated = ? WHERE id = ?`
	case "completed":
		query = `UPDATE car_track SET status = ?, dtmCompleted = ?, dtmUpdated = ? WHERE id = ?`
	default:
		query = `UPDATE car_track SET status = ?, dtmUpdated = ? WHERE id = ?`
		_, err := r.db.ExecContext(ctx, query, status, dtmUpdated, id)
		return err
	}

	_, err := r.db.ExecContext(ctx, query, status, dtmUpdated, dtmUpdated, id)
	if err != nil {
		return fmt.Errorf("failed to update car task: %w", err)
	}

	return nil
}

// UpdatePosition updates car position
func (r *CarRepository) UpdatePosition(ctx context.Context, userName, room string, lat, lng, accuracy, altitude, altitudeAccuracy float64, speed, heading int) error {
	dtmUpdated := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE car_track SET
		latitude = ?, longitude = ?, accuracy = ?, speed = ?, heading = ?,
		altitude = ?, altitudeAccuracy = ?, dtmUpdated = ?
		WHERE userName = ? AND room = ?`

	_, err := r.db.ExecContext(ctx, query,
		lat, lng, accuracy, speed, heading,
		altitude, altitudeAccuracy, dtmUpdated,
		userName, room,
	)
	if err != nil {
		return fmt.Errorf("failed to update car position: %w", err)
	}

	return nil
}

// GetCarPosition gets car position by room
func (r *CarRepository) GetCarPosition(ctx context.Context, room string) (*models.CarTrack, error) {
	var task models.CarTrack
	query := `SELECT latitude, longitude, accuracy, speed, heading, altitude, altitudeAccuracy
		FROM car_track WHERE room = ? ORDER BY dtmUpdated DESC LIMIT 1`

	err := r.db.GetContext(ctx, &task, query, room)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get car position: %w", err)
	}

	return &task, nil
}

// GetTaskList gets list of car tasks
func (r *CarRepository) GetTaskList(ctx context.Context, status string, limit, offset int) ([]models.CarTrack, error) {
	var tasks []models.CarTrack
	query := `SELECT * FROM car_track`
	args := []interface{}{}

	if status != "" {
		query += ` WHERE status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY dtmCreated DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &tasks, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get task list: %w", err)
	}

	return tasks, nil
}

// GetUserLatLng gets user lat/lng by room
func (r *CarRepository) GetUserLatLng(ctx context.Context, room string) (*models.CarTrack, error) {
	var task models.CarTrack
	query := `SELECT latitude, longitude FROM link_connect WHERE room = ? AND userType = 'user' AND latitude IS NOT NULL LIMIT 1`

	err := r.db.GetContext(ctx, &task, query, room)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user lat/lng: %w", err)
	}

	return &task, nil
}

// Delete deletes a car task
func (r *CarRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM car_track WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete car task: %w", err)
	}

	return nil
}
