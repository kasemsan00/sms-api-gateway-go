package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// CaseRepository handles case database operations
type CaseRepository struct {
	db *sqlx.DB
}

// NewCaseRepository creates a new CaseRepository
func NewCaseRepository(db *sqlx.DB) *CaseRepository {
	return &CaseRepository{db: db}
}

// CreateCaseParams holds parameters for creating a case
type CreateCaseParams struct {
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

// Create creates a new case
func (r *CaseRepository) Create(ctx context.Context, params CreateCaseParams) (int64, error) {
	dtmCreated := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO case_data
		(caseId, service, roomId, operationNumber, status, hn, patientMobile, mobileCreated, caseType, userName, organization, dtmCreated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		params.CaseID,
		params.Service,
		params.RoomID,
		params.OperationNumber,
		params.Status,
		params.HN,
		params.PatientMobile,
		params.MobileCreated,
		params.CaseType,
		params.UserName,
		params.Organization,
		dtmCreated,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create case: %w", err)
	}

	return result.LastInsertId()
}

// GetByCaseID gets case by caseId
func (r *CaseRepository) GetByCaseID(ctx context.Context, caseID int) (*models.CaseData, error) {
	var caseData models.CaseData
	query := `SELECT * FROM case_data WHERE caseId = ?`

	err := r.db.GetContext(ctx, &caseData, query, caseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get case: %w", err)
	}

	return &caseData, nil
}

// GetByID gets case by ID
func (r *CaseRepository) GetByID(ctx context.Context, id uint) (*models.CaseData, error) {
	var caseData models.CaseData
	query := `SELECT * FROM case_data WHERE id = ?`

	err := r.db.GetContext(ctx, &caseData, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get case by ID: %w", err)
	}

	return &caseData, nil
}

// Update updates a case
func (r *CaseRepository) Update(ctx context.Context, id uint, status, hn, patientMobile, caseType, userName string) error {
	query := `UPDATE case_data SET status = ?, hn = ?, patientMobile = ?, caseType = ?, userName = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, status, hn, patientMobile, caseType, userName, id)
	if err != nil {
		return fmt.Errorf("failed to update case: %w", err)
	}

	return nil
}

// UpdateStatus updates case status
func (r *CaseRepository) UpdateStatus(ctx context.Context, caseID int, status string) error {
	query := `UPDATE case_data SET status = ? WHERE caseId = ?`

	_, err := r.db.ExecContext(ctx, query, status, caseID)
	if err != nil {
		return fmt.Errorf("failed to update case status: %w", err)
	}

	return nil
}

// GetHistory gets case history with pagination
func (r *CaseRepository) GetHistory(ctx context.Context, service, limit, offset int) ([]models.CaseData, error) {
	var cases []models.CaseData
	query := `SELECT * FROM case_data WHERE service = ? ORDER BY dtmCreated DESC LIMIT ? OFFSET ?`

	err := r.db.SelectContext(ctx, &cases, query, service, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get case history: %w", err)
	}

	return cases, nil
}

// GetHistoryCount gets case history count
func (r *CaseRepository) GetHistoryCount(ctx context.Context, service int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM case_data WHERE service = ?`

	err := r.db.GetContext(ctx, &count, query, service)
	if err != nil {
		return 0, fmt.Errorf("failed to get case count: %w", err)
	}

	return count, nil
}

// GetRoomName gets room name by caseId and service
func (r *CaseRepository) GetRoomName(ctx context.Context, caseID, service int) (string, error) {
	var room string
	query := `SELECT room_conference.room FROM case_data
		LEFT JOIN room_conference ON case_data.roomId = room_conference.id
		WHERE case_data.caseId = ? AND case_data.service = ?`

	err := r.db.GetContext(ctx, &room, query, caseID, service)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get room name: %w", err)
	}

	return room, nil
}

// GetByRoomID gets case by roomId
func (r *CaseRepository) GetByRoomID(ctx context.Context, roomID int) (*models.CaseData, error) {
	var caseData models.CaseData
	query := `SELECT * FROM case_data WHERE roomId = ?`

	err := r.db.GetContext(ctx, &caseData, query, roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get case by room ID: %w", err)
	}

	return &caseData, nil
}

// GetByService gets cases by service
func (r *CaseRepository) GetByService(ctx context.Context, service int) ([]models.CaseData, error) {
	var cases []models.CaseData
	query := `SELECT * FROM case_data WHERE service = ? ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &cases, query, service)
	if err != nil {
		return nil, fmt.Errorf("failed to get cases by service: %w", err)
	}

	return cases, nil
}

// Delete deletes a case
func (r *CaseRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM case_data WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete case: %w", err)
	}

	return nil
}
