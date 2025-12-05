package service

import (
	"context"

	"api-gateway-go/internal/repository"
)

// StatsService handles statistics business logic
type StatsService struct {
	statsRepo *repository.StatsRepository
}

// NewStatsService creates a new StatsService
func NewStatsService(statsRepo *repository.StatsRepository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
	}
}

// GetSummary gets statistics summary
func (s *StatsService) GetSummary(ctx context.Context, service int) (*repository.StatsSummary, error) {
	return s.statsRepo.GetSummary(ctx, service)
}

// GetDeviceStats gets device statistics
func (s *StatsService) GetDeviceStats(ctx context.Context, service int) ([]repository.DeviceStats, error) {
	return s.statsRepo.GetDeviceStats(ctx, service)
}

// GetTypeStats gets type statistics
func (s *StatsService) GetTypeStats(ctx context.Context, service int) ([]repository.TypeStats, error) {
	return s.statsRepo.GetTypeStats(ctx, service)
}

// GetUserStats gets user statistics
func (s *StatsService) GetUserStats(ctx context.Context, service int) ([]repository.UserStats, error) {
	return s.statsRepo.GetUserStats(ctx, service)
}

// GetCaseStats gets case statistics
func (s *StatsService) GetCaseStats(ctx context.Context, service int) ([]repository.CaseStats, error) {
	return s.statsRepo.GetCaseStats(ctx, service)
}

// GetDailyStats gets daily statistics
func (s *StatsService) GetDailyStats(ctx context.Context, startDate, endDate string, service int) ([]map[string]interface{}, error) {
	return s.statsRepo.GetDailyStats(ctx, startDate, endDate, service)
}

// GetMonthlyStats gets monthly statistics
func (s *StatsService) GetMonthlyStats(ctx context.Context, year, service int) ([]map[string]interface{}, error) {
	return s.statsRepo.GetMonthlyStats(ctx, year, service)
}
