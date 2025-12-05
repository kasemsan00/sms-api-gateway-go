package service

import (
	"context"
	"database/sql"
	"errors"

	"api-gateway-go/internal/models"
	"api-gateway-go/internal/repository"
)

var (
	ErrRadioDeviceNotFound   = errors.New("radio device not found")
	ErrRadioLocationNotFound = errors.New("radio location not found")
)

// RadioService handles radio device and location business logic
type RadioService struct {
	radioRepo *repository.RadioRepository
}

// NewRadioService creates a new RadioService
func NewRadioService(radioRepo *repository.RadioRepository) *RadioService {
	return &RadioService{
		radioRepo: radioRepo,
	}
}

// GetDevices gets all radio devices
func (s *RadioService) GetDevices(ctx context.Context, limit int) ([]models.RadioDevices, error) {
	if limit == 0 {
		limit = 100
	}
	return s.radioRepo.GetDevices(ctx, limit)
}

// GetDeviceByID gets a device by ID
func (s *RadioService) GetDeviceByID(ctx context.Context, id int) (*models.RadioDevices, error) {
	device, err := s.radioRepo.GetDeviceByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRadioDeviceNotFound
		}
		return nil, err
	}
	if device == nil {
		return nil, ErrRadioDeviceNotFound
	}
	return device, nil
}

// GetDeviceByDeviceID gets a device by device ID
func (s *RadioService) GetDeviceByDeviceID(ctx context.Context, deviceID string) (*models.RadioDevices, error) {
	device, err := s.radioRepo.GetDeviceByDeviceID(ctx, deviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRadioDeviceNotFound
		}
		return nil, err
	}
	if device == nil {
		return nil, ErrRadioDeviceNotFound
	}
	return device, nil
}

// CreateDevice creates a new radio device
func (s *RadioService) CreateDevice(ctx context.Context, deviceID, radioNo, channel, status string) (int64, error) {
	return s.radioRepo.CreateDevice(ctx, deviceID, radioNo, channel, status)
}

// UpdateDevice updates a radio device
func (s *RadioService) UpdateDevice(ctx context.Context, id int, radioNo, channel, status string) error {
	return s.radioRepo.UpdateDevice(ctx, id, radioNo, channel, status)
}

// UpdateDeviceLocation updates device location
func (s *RadioService) UpdateDeviceLocation(ctx context.Context, deviceID string, latitude, longitude float64) error {
	return s.radioRepo.UpdateDeviceLocation(ctx, deviceID, latitude, longitude)
}

// DeleteDevice deletes a radio device
func (s *RadioService) DeleteDevice(ctx context.Context, id int) error {
	return s.radioRepo.DeleteDevice(ctx, id)
}

// GetLocations gets all radio locations
func (s *RadioService) GetLocations(ctx context.Context, limit int) ([]models.RadioLocations, error) {
	if limit == 0 {
		limit = 100
	}
	return s.radioRepo.GetLocations(ctx, limit)
}

// GetLocationByRadioNo gets location by radio number
func (s *RadioService) GetLocationByRadioNo(ctx context.Context, radioNo string) (*models.RadioLocations, error) {
	location, err := s.radioRepo.GetLocationByRadioNo(ctx, radioNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRadioLocationNotFound
		}
		return nil, err
	}
	if location == nil {
		return nil, ErrRadioLocationNotFound
	}
	return location, nil
}

// CreateLocation creates a new radio location
func (s *RadioService) CreateLocation(ctx context.Context, radioNo string, latitude, longitude float64) (int64, error) {
	return s.radioRepo.CreateLocation(ctx, radioNo, latitude, longitude)
}
