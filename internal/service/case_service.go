package service

import (
	"context"
	"database/sql"
	"errors"

	"api-gateway-go/internal/models"
	"api-gateway-go/internal/repository"
)

var (
	ErrCaseNotFound = errors.New("case not found")
)

// CreateCaseOptions holds options for creating a case
type CreateCaseOptions struct {
	CaseID          int
	Service         int
	RoomID          int
	OperationNumber string
	Status          string
	HN              string
	PatientMobile   string
	MobileCreated   string
	CaseType        string
	UserName        string
	Organization    string
}

// CaseService handles case business logic
type CaseService struct {
	caseRepo *repository.CaseRepository
}

// NewCaseService creates a new CaseService
func NewCaseService(caseRepo *repository.CaseRepository) *CaseService {
	return &CaseService{
		caseRepo: caseRepo,
	}
}

// CreateCase creates a new case
func (s *CaseService) CreateCase(ctx context.Context, opts CreateCaseOptions) (int64, error) {
	params := repository.CreateCaseParams{
		CaseID:          opts.CaseID,
		Service:         opts.Service,
		RoomID:          opts.RoomID,
		OperationNumber: opts.OperationNumber,
		Status:          opts.Status,
		HN:              opts.HN,
		PatientMobile:   opts.PatientMobile,
		MobileCreated:   opts.MobileCreated,
		CaseType:        opts.CaseType,
		UserName:        opts.UserName,
		Organization:    opts.Organization,
	}
	return s.caseRepo.Create(ctx, params)
}

// GetCaseByID gets a case by ID
func (s *CaseService) GetCaseByID(ctx context.Context, id uint) (*models.CaseData, error) {
	caseData, err := s.caseRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCaseNotFound
		}
		return nil, err
	}
	if caseData == nil {
		return nil, ErrCaseNotFound
	}
	return caseData, nil
}

// GetCaseByCaseID gets a case by case ID
func (s *CaseService) GetCaseByCaseID(ctx context.Context, caseID int) (*models.CaseData, error) {
	caseData, err := s.caseRepo.GetByCaseID(ctx, caseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCaseNotFound
		}
		return nil, err
	}
	if caseData == nil {
		return nil, ErrCaseNotFound
	}
	return caseData, nil
}

// GetCaseByRoomID gets a case by room ID
func (s *CaseService) GetCaseByRoomID(ctx context.Context, roomID int) (*models.CaseData, error) {
	caseData, err := s.caseRepo.GetByRoomID(ctx, roomID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCaseNotFound
		}
		return nil, err
	}
	if caseData == nil {
		return nil, ErrCaseNotFound
	}
	return caseData, nil
}

// UpdateCase updates a case
func (s *CaseService) UpdateCase(ctx context.Context, id uint, status, hn, patientMobile, caseType, userName string) error {
	return s.caseRepo.Update(ctx, id, status, hn, patientMobile, caseType, userName)
}

// UpdateCaseStatus updates case status
func (s *CaseService) UpdateCaseStatus(ctx context.Context, caseID int, status string) error {
	return s.caseRepo.UpdateStatus(ctx, caseID, status)
}

// GetCaseHistory gets case history with pagination
func (s *CaseService) GetCaseHistory(ctx context.Context, service, limit, offset int) ([]models.CaseData, error) {
	if limit == 0 {
		limit = 100
	}
	return s.caseRepo.GetHistory(ctx, service, limit, offset)
}

// GetCaseHistoryCount gets case history count
func (s *CaseService) GetCaseHistoryCount(ctx context.Context, service int) (int, error) {
	return s.caseRepo.GetHistoryCount(ctx, service)
}

// GetCasesByService gets cases by service
func (s *CaseService) GetCasesByService(ctx context.Context, service int) ([]models.CaseData, error) {
	return s.caseRepo.GetByService(ctx, service)
}

// GetRoomNameByCaseID gets room name by case ID and service
func (s *CaseService) GetRoomNameByCaseID(ctx context.Context, caseID, service int) (string, error) {
	return s.caseRepo.GetRoomName(ctx, caseID, service)
}

// DeleteCase deletes a case
func (s *CaseService) DeleteCase(ctx context.Context, id uint) error {
	return s.caseRepo.Delete(ctx, id)
}
