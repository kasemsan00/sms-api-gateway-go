package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/models"
	"api-gateway-go/internal/repository"
	"api-gateway-go/pkg/utils"
)

var (
	ErrLinkNotFound    = errors.New("link not found")
	ErrLinkExpired     = errors.New("link has expired")
	ErrLinkDisabled    = errors.New("link is disabled")
	ErrOneTimeLinkUsed = errors.New("one-time link already used")
)

// CreateLinkOptions holds options for creating a link
type CreateLinkOptions struct {
	Room                  string
	UserType              string
	Mobile                string
	Share                 int
	RequireJoinPermission int
	RequireUserName       int
	RequirePassword       int
	Password              string
	OneTimeLink           int
	UserAgent             string
	DaysExpired           int
}

// LinkService handles link business logic
type LinkService struct {
	linkRepo *repository.LinkRepository
	roomRepo *repository.RoomRepository
	cfg      *config.Config
}

// NewLinkService creates a new LinkService
func NewLinkService(linkRepo *repository.LinkRepository, roomRepo *repository.RoomRepository, cfg *config.Config) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
		roomRepo: roomRepo,
		cfg:      cfg,
	}
}

// CreateLink creates a new link
func (s *LinkService) CreateLink(ctx context.Context, opts CreateLinkOptions) (*models.LinkConnect, error) {
	// Generate link ID
	linkID := utils.GenerateLinkID(s.cfg.CustomCharset)
	now := utils.FormatDateTimeNow()

	// Calculate expiration
	daysExpired := opts.DaysExpired
	if daysExpired == 0 {
		daysExpired = s.cfg.RoomDayDefaultTimeout
	}
	expiredAt := utils.AddDays(time.Now(), daysExpired)

	// Create link in database
	params := repository.CreateLinkParams{
		Share:                 opts.Share,
		Mobile:                opts.Mobile,
		LinkID:                linkID,
		Room:                  opts.Room,
		LinkType:              opts.UserType,
		RequireJoinPermission: opts.RequireJoinPermission,
		RequireUserName:       opts.RequireUserName,
		RequirePassword:       opts.RequirePassword,
		Password:              opts.Password,
		OneTimeLink:           opts.OneTimeLink,
		UserAgent:             opts.UserAgent,
		DtmCreated:            now,
		DtmExpired:            expiredAt,
	}

	_, err := s.linkRepo.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.GetLinkDetail(ctx, linkID, "")
}

// GetLinkDetail gets link details by linkID
func (s *LinkService) GetLinkDetail(ctx context.Context, linkID, room string) (*models.LinkConnect, error) {
	if linkID != "" {
		link, err := s.linkRepo.GetByLinkID(ctx, linkID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrLinkNotFound
			}
			return nil, err
		}
		return link, nil
	}

	if room != "" {
		// Get by room with default user type
		link, err := s.linkRepo.GetByRoomAndUserType(ctx, room, "guest")
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrLinkNotFound
			}
			return nil, err
		}
		return link, nil
	}

	return nil, ErrLinkNotFound
}

// GetLinkDetailByRoomAndType gets link by room and user type
func (s *LinkService) GetLinkDetailByRoomAndType(ctx context.Context, room, userType string) (*models.LinkConnect, error) {
	link, err := s.linkRepo.GetByRoomAndUserType(ctx, room, userType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}
	return link, nil
}

// GetShareURL gets the share URL for a room
func (s *LinkService) GetShareURL(ctx context.Context, room, userType string) (string, error) {
	return s.linkRepo.GetShareURL(ctx, room, userType)
}

// UpdateLatLng updates latitude and longitude
func (s *LinkService) UpdateLatLng(ctx context.Context, linkID string, latitude, longitude float64, accuracy int) error {
	return s.linkRepo.UpdateLatLng(ctx, linkID, latitude, longitude, accuracy)
}

// UpdatePatientLocation updates patient location
func (s *LinkService) UpdatePatientLocation(ctx context.Context, linkID string, lat, lng float64) error {
	return s.linkRepo.UpdatePatientLocation(ctx, linkID, lat, lng)
}

// CheckAndUpdateOneTimeLink checks and updates one-time link status
func (s *LinkService) CheckAndUpdateOneTimeLink(ctx context.Context, linkID string) error {
	status, err := s.linkRepo.GetOneTimeLinkStatus(ctx, linkID)
	if err != nil {
		return err
	}

	if status == 1 {
		return ErrOneTimeLinkUsed
	}

	return s.linkRepo.UpdateOneTimeLink(ctx, linkID, 1)
}

// GetLinkIDList gets links by room and mobile
func (s *LinkService) GetLinkIDList(ctx context.Context, room, mobile string) ([]models.LinkConnect, error) {
	return s.linkRepo.GetLinkIDList(ctx, room, mobile)
}

// UpdateLinkEnabled updates link enabled status
func (s *LinkService) UpdateLinkEnabled(ctx context.Context, room, linkType string, enabled int) error {
	return s.linkRepo.UpdateLinkEnabled(ctx, room, linkType, enabled)
}

// GetUserMobile gets user mobile by room
func (s *LinkService) GetUserMobile(ctx context.Context, room string) (string, error) {
	return s.linkRepo.GetUserMobile(ctx, room)
}

// UpdateUserAgent updates user agent
func (s *LinkService) UpdateUserAgent(ctx context.Context, linkID, userAgent, os string) error {
	return s.linkRepo.UpdateUserAgent(ctx, linkID, userAgent, os)
}

// UpdateErrorLocation updates error location
func (s *LinkService) UpdateErrorLocation(ctx context.Context, linkID, errorMsg string) error {
	return s.linkRepo.UpdateErrorLocation(ctx, linkID, errorMsg)
}

// UpdateErrorVideo updates error video
func (s *LinkService) UpdateErrorVideo(ctx context.Context, linkID, errorMsg string) error {
	return s.linkRepo.UpdateErrorVideo(ctx, linkID, errorMsg)
}

// GetLastLatLng gets last latitude/longitude
func (s *LinkService) GetLastLatLng(ctx context.Context, room, userType string, share int) (*models.LinkConnect, error) {
	return s.linkRepo.GetLastLatLng(ctx, room, userType, share)
}

// GetLatLngGroup gets lat/lng group
func (s *LinkService) GetLatLngGroup(ctx context.Context, room, userType string) ([]models.LinkConnect, error) {
	return s.linkRepo.GetLatLngGroup(ctx, room, userType)
}

// UpdateLinkUserName updates link username
func (s *LinkService) UpdateLinkUserName(ctx context.Context, linkID, userName string) error {
	return s.linkRepo.UpdateLinkUserName(ctx, linkID, userName)
}

// AutoLinkExpiredClose closes expired links
func (s *LinkService) AutoLinkExpiredClose(ctx context.Context) error {
	// This would typically query for expired links and disable them
	// For now, we'll implement a simple version
	return nil
}

// GetDomain gets domain for a service
func (s *LinkService) GetDomain(ctx context.Context, service int, sender, linkType, linkID string) (string, error) {
	// Get last domain index
	lastIndex, err := s.linkRepo.GetLastDomainIndex(ctx, linkID)
	if err != nil {
		lastIndex = 0
	}

	// Update domain index
	err = s.linkRepo.UpdateDomainIndex(ctx, linkID, lastIndex+1)
	if err != nil {
		return "", err
	}

	// Return API URL
	return s.cfg.APIURL, nil
}
