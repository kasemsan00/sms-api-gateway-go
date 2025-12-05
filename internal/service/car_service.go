package service

import (
	"context"
	"database/sql"
	"errors"

	"api-gateway-go/internal/models"
	"api-gateway-go/internal/repository"
)

var (
	ErrCarTaskNotFound = errors.New("car task not found")
)

// CreateCarTaskOptions holds options for creating a car task
type CreateCarTaskOptions struct {
	UID      string
	Status   string
	Mobile   string
	UserName string
	Room     string
}

// CarService handles car task business logic
type CarService struct {
	carRepo *repository.CarRepository
}

// NewCarService creates a new CarService
func NewCarService(carRepo *repository.CarRepository) *CarService {
	return &CarService{
		carRepo: carRepo,
	}
}

// CreateTask creates a new car task
func (s *CarService) CreateTask(ctx context.Context, opts CreateCarTaskOptions) (int64, error) {
	params := repository.CreateTaskParams{
		UID:      opts.UID,
		Status:   opts.Status,
		Mobile:   opts.Mobile,
		UserName: opts.UserName,
		Room:     opts.Room,
	}
	return s.carRepo.CreateTask(ctx, params)
}

// GetTaskDetail gets task details by ID
func (s *CarService) GetTaskDetail(ctx context.Context, id int) (*models.CarTrack, error) {
	task, err := s.carRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCarTaskNotFound
		}
		return nil, err
	}
	if task == nil {
		return nil, ErrCarTaskNotFound
	}
	return task, nil
}

// GetTaskByUID gets task by UID
func (s *CarService) GetTaskByUID(ctx context.Context, uid string) (*models.CarTrack, error) {
	task, err := s.carRepo.GetByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCarTaskNotFound
		}
		return nil, err
	}
	if task == nil {
		return nil, ErrCarTaskNotFound
	}
	return task, nil
}

// GetTaskByRoom gets task by room
func (s *CarService) GetTaskByRoom(ctx context.Context, room string) (*models.CarTrack, error) {
	task, err := s.carRepo.GetByRoom(ctx, room)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCarTaskNotFound
		}
		return nil, err
	}
	if task == nil {
		return nil, ErrCarTaskNotFound
	}
	return task, nil
}

// UpdateTask updates a car task status
func (s *CarService) UpdateTask(ctx context.Context, id int, status string) error {
	return s.carRepo.UpdateTask(ctx, id, status)
}

// UpdatePosition updates car position
func (s *CarService) UpdatePosition(ctx context.Context, userName, room string, lat, lng, accuracy, altitude, altitudeAccuracy float64, speed, heading int) error {
	return s.carRepo.UpdatePosition(ctx, userName, room, lat, lng, accuracy, altitude, altitudeAccuracy, speed, heading)
}

// GetTaskList gets list of car tasks
func (s *CarService) GetTaskList(ctx context.Context, status string, limit, offset int) ([]models.CarTrack, error) {
	if limit == 0 {
		limit = 100
	}
	return s.carRepo.GetTaskList(ctx, status, limit, offset)
}

// GetCarPosition gets car position by room
func (s *CarService) GetCarPosition(ctx context.Context, room string) (*models.CarTrack, error) {
	return s.carRepo.GetCarPosition(ctx, room)
}

// DeleteTask deletes a car task
func (s *CarService) DeleteTask(ctx context.Context, id int) error {
	return s.carRepo.Delete(ctx, id)
}

// GetUserLatLng gets user lat/lng by room
func (s *CarService) GetUserLatLng(ctx context.Context, room string) (*models.CarTrack, error) {
	return s.carRepo.GetUserLatLng(ctx, room)
}
