package service

import (
	"context"
	"time"

	"api-gateway-go/internal/config"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// CronStatus represents cron job status
type CronStatus struct {
	Running   bool          `json:"running"`
	Jobs      []CronJobInfo `json:"jobs"`
	LastCheck time.Time     `json:"lastCheck"`
}

// CronJobInfo represents info about a cron job
type CronJobInfo struct {
	Name     string    `json:"name"`
	Schedule string    `json:"schedule"`
	LastRun  time.Time `json:"lastRun"`
	NextRun  time.Time `json:"nextRun"`
}

// CrontabService handles cron jobs
type CrontabService struct {
	cron        *cron.Cron
	roomService *RoomService
	linkService *LinkService
	livekitMgr  *config.LiveKitManager
	cfg         *config.Config
	status      *CronStatus
}

// NewCrontabService creates a new CrontabService
func NewCrontabService(roomService *RoomService, linkService *LinkService, livekitMgr *config.LiveKitManager, cfg *config.Config) *CrontabService {
	// Create cron with Bangkok timezone
	loc, _ := time.LoadLocation("Asia/Bangkok")
	c := cron.New(cron.WithLocation(loc))

	return &CrontabService{
		cron:        c,
		roomService: roomService,
		linkService: linkService,
		livekitMgr:  livekitMgr,
		cfg:         cfg,
		status: &CronStatus{
			Running: false,
			Jobs:    []CronJobInfo{},
		},
	}
}

// InitCronJobs initializes all cron jobs
func (s *CrontabService) InitCronJobs() error {
	// Room cleanup - every 30 minutes
	_, err := s.cron.AddFunc("*/30 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		log.Info().Msg("Running room cleanup cron job")
		if err := s.roomService.AutoRoomExpiredClose(ctx); err != nil {
			log.Error().Err(err).Msg("Room cleanup cron job failed")
		}
	})
	if err != nil {
		return err
	}

	// Link cleanup - every 30 minutes
	_, err = s.cron.AddFunc("*/30 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		log.Info().Msg("Running link cleanup cron job")
		if err := s.linkService.AutoLinkExpiredClose(ctx); err != nil {
			log.Error().Err(err).Msg("Link cleanup cron job failed")
		}
	})
	if err != nil {
		return err
	}

	// LiveKit health check - every 10 minutes
	_, err = s.cron.AddFunc("*/10 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		log.Debug().Msg("Running LiveKit health check")
		if s.livekitMgr != nil {
			if err := s.livekitMgr.Health(ctx); err != nil {
				log.Error().Err(err).Msg("LiveKit health check failed")
			}
		}
	})
	if err != nil {
		return err
	}

	// Start the cron scheduler
	s.cron.Start()
	s.status.Running = true
	s.status.LastCheck = time.Now()

	log.Info().Msg("Cron jobs initialized and started")

	return nil
}

// GetStatus returns the cron status
func (s *CrontabService) GetStatus() *CronStatus {
	s.status.LastCheck = time.Now()

	// Get job info from cron entries
	entries := s.cron.Entries()
	s.status.Jobs = make([]CronJobInfo, len(entries))

	jobNames := []string{"Room Cleanup", "Link Cleanup", "LiveKit Health Check"}
	for i, entry := range entries {
		name := "Unknown"
		if i < len(jobNames) {
			name = jobNames[i]
		}
		s.status.Jobs[i] = CronJobInfo{
			Name:    name,
			NextRun: entry.Next,
		}
	}

	return s.status
}

// Stop stops the cron scheduler
func (s *CrontabService) Stop() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		s.status.Running = false
		log.Info().Msg("Cron jobs stopped")
	}
}

// Cleanup cleans up resources
func (s *CrontabService) Cleanup() error {
	s.Stop()
	return nil
}
