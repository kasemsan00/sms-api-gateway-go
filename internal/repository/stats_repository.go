package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// StatsRepository handles statistics database operations
type StatsRepository struct {
	db *sqlx.DB
}

// NewStatsRepository creates a new StatsRepository
func NewStatsRepository(db *sqlx.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// StatsSummary represents summary statistics
type StatsSummary struct {
	TotalRooms       int `db:"totalRooms" json:"totalRooms"`
	OpenRooms        int `db:"openRooms" json:"openRooms"`
	ClosedRooms      int `db:"closedRooms" json:"closedRooms"`
	TotalUsers       int `db:"totalUsers" json:"totalUsers"`
	ActiveUsers      int `db:"activeUsers" json:"activeUsers"`
	TotalLinks       int `db:"totalLinks" json:"totalLinks"`
	TotalRecords     int `db:"totalRecords" json:"totalRecords"`
	TotalCases       int `db:"totalCases" json:"totalCases"`
}

// DeviceStats represents device statistics
type DeviceStats struct {
	Device string `db:"device" json:"device"`
	Count  int    `db:"count" json:"count"`
}

// TypeStats represents type statistics
type TypeStats struct {
	Type  string `db:"type" json:"type"`
	Count int    `db:"count" json:"count"`
}

// UserStats represents user statistics
type UserStats struct {
	UserType string `db:"userType" json:"userType"`
	Count    int    `db:"count" json:"count"`
}

// CaseStats represents case statistics
type CaseStats struct {
	Status string `db:"status" json:"status"`
	Count  int    `db:"count" json:"count"`
}

// GetSummary gets summary statistics
func (r *StatsRepository) GetSummary(ctx context.Context, service int) (*StatsSummary, error) {
	var summary StatsSummary

	// Get room counts
	roomQuery := `SELECT
		COUNT(*) as totalRooms,
		SUM(CASE WHEN status = 'open' THEN 1 ELSE 0 END) as openRooms,
		SUM(CASE WHEN status = 'close' THEN 1 ELSE 0 END) as closedRooms
		FROM room_conference`
	if service > 0 {
		roomQuery += fmt.Sprintf(` WHERE service = %d`, service)
	}

	err := r.db.GetContext(ctx, &summary, roomQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get room stats: %w", err)
	}

	// Get user counts
	userQuery := `SELECT
		COUNT(*) as totalUsers,
		SUM(CASE WHEN status = 'connection' THEN 1 ELSE 0 END) as activeUsers
		FROM room_user`

	var userStats struct {
		TotalUsers  int `db:"totalUsers"`
		ActiveUsers int `db:"activeUsers"`
	}
	err = r.db.GetContext(ctx, &userStats, userQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	summary.TotalUsers = userStats.TotalUsers
	summary.ActiveUsers = userStats.ActiveUsers

	// Get link count
	linkQuery := `SELECT COUNT(*) FROM link_connect`
	err = r.db.GetContext(ctx, &summary.TotalLinks, linkQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get link stats: %w", err)
	}

	// Get record count
	recordQuery := `SELECT COUNT(*) FROM record_media`
	err = r.db.GetContext(ctx, &summary.TotalRecords, recordQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get record stats: %w", err)
	}

	// Get case count
	caseQuery := `SELECT COUNT(*) FROM case_data`
	if service > 0 {
		caseQuery += fmt.Sprintf(` WHERE service = %d`, service)
	}
	err = r.db.GetContext(ctx, &summary.TotalCases, caseQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get case stats: %w", err)
	}

	return &summary, nil
}

// GetDeviceStats gets device statistics
func (r *StatsRepository) GetDeviceStats(ctx context.Context, service int) ([]DeviceStats, error) {
	var stats []DeviceStats
	query := `SELECT
		CASE
			WHEN userAgent LIKE '%iPhone%' OR userAgent LIKE '%iPad%' THEN 'iOS'
			WHEN userAgent LIKE '%Android%' THEN 'Android'
			WHEN userAgent LIKE '%Windows%' THEN 'Windows'
			WHEN userAgent LIKE '%Mac%' THEN 'Mac'
			WHEN userAgent LIKE '%Linux%' THEN 'Linux'
			ELSE 'Other'
		END as device,
		COUNT(*) as count
		FROM link_connect
		GROUP BY device
		ORDER BY count DESC`

	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats: %w", err)
	}

	return stats, nil
}

// GetTypeStats gets link type statistics
func (r *StatsRepository) GetTypeStats(ctx context.Context, service int) ([]TypeStats, error) {
	var stats []TypeStats
	query := `SELECT linkType as type, COUNT(*) as count
		FROM link_connect
		WHERE linkType IS NOT NULL AND linkType != ''
		GROUP BY linkType
		ORDER BY count DESC`

	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get type stats: %w", err)
	}

	return stats, nil
}

// GetUserStats gets user type statistics
func (r *StatsRepository) GetUserStats(ctx context.Context, service int) ([]UserStats, error) {
	var stats []UserStats
	query := `SELECT userType, COUNT(*) as count
		FROM room_user
		WHERE userType IS NOT NULL AND userType != ''
		GROUP BY userType
		ORDER BY count DESC`

	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return stats, nil
}

// GetCaseStats gets case status statistics
func (r *StatsRepository) GetCaseStats(ctx context.Context, service int) ([]CaseStats, error) {
	var stats []CaseStats
	query := `SELECT status, COUNT(*) as count
		FROM case_data`
	if service > 0 {
		query += fmt.Sprintf(` WHERE service = %d`, service)
	}
	query += ` GROUP BY status ORDER BY count DESC`

	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get case stats: %w", err)
	}

	return stats, nil
}

// GetDailyStats gets daily statistics for a date range
func (r *StatsRepository) GetDailyStats(ctx context.Context, startDate, endDate string, service int) ([]map[string]interface{}, error) {
	query := `SELECT
		DATE(dtmCreated) as date,
		COUNT(*) as count
		FROM room_conference
		WHERE dtmCreated BETWEEN ? AND ?`
	args := []interface{}{startDate, endDate}

	if service > 0 {
		query += ` AND service = ?`
		args = append(args, service)
	}
	query += ` GROUP BY DATE(dtmCreated) ORDER BY date`

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily stats: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// GetMonthlyStats gets monthly statistics
func (r *StatsRepository) GetMonthlyStats(ctx context.Context, year int, service int) ([]map[string]interface{}, error) {
	query := `SELECT
		MONTH(dtmCreated) as month,
		COUNT(*) as count
		FROM room_conference
		WHERE YEAR(dtmCreated) = ?`
	args := []interface{}{year}

	if service > 0 {
		query += ` AND service = ?`
		args = append(args, service)
	}
	query += ` GROUP BY MONTH(dtmCreated) ORDER BY month`

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}
