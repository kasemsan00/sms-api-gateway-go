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

	"github.com/livekit/protocol/auth"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyInRoom = errors.New("user already in room")
)

// GenerateUserOptions holds options for generating a user
type GenerateUserOptions struct {
	LinkID    string
	Room      string
	UserName  string
	UserType  string
	UserAgent string
	ServiceID int
}

// GenerateUserResult holds the result of generating a user
type GenerateUserResult struct {
	Token    string `json:"token"`
	Identity string `json:"identity"`
	Room     string `json:"room"`
	URL      string `json:"url"`
	UserName string `json:"userName,omitempty"`
	Color    string `json:"color,omitempty"`
}

// UserService handles user business logic
type UserService struct {
	userRepo   *repository.UserRepository
	roomRepo   *repository.RoomRepository
	livekitMgr *config.LiveKitManager
	cfg        *config.Config
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repository.UserRepository, roomRepo *repository.RoomRepository, livekitMgr *config.LiveKitManager, cfg *config.Config) *UserService {
	return &UserService{
		userRepo:   userRepo,
		roomRepo:   roomRepo,
		livekitMgr: livekitMgr,
		cfg:        cfg,
	}
}

// GenerateUser generates a new user with LiveKit token
func (s *UserService) GenerateUser(ctx context.Context, opts GenerateUserOptions) (*GenerateUserResult, error) {
	// Generate identity
	identity := utils.GenerateIdentity()

	// Get random color
	color, err := s.userRepo.GetRandomColor(ctx)
	if err != nil {
		color = "#000000" // Default color
	}

	// Set default username if not provided
	userName := opts.UserName
	if userName == "" {
		userName = utils.GenerateGuestName()
	}

	// Add user to database
	err = s.userRepo.AddUser(ctx, repository.AddUserParams{
		Room:       opts.Room,
		Identity:   identity,
		UserName:   userName,
		UserType:   opts.UserType,
		Status:     "connect",
		SocketID:   "",
		Color:      color,
		Conference: 0,
		UserAgent:  opts.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	// Generate LiveKit token
	token, err := s.GenerateLiveKitToken(ctx, identity, userName, opts.Room, opts.UserType)
	if err != nil {
		return nil, err
	}

	// Get LiveKit URL
	url := s.cfg.LiveKitHost

	return &GenerateUserResult{
		Token:    token,
		Identity: identity,
		Room:     opts.Room,
		URL:      url,
		UserName: userName,
		Color:    color,
	}, nil
}

// GenerateLiveKitToken generates a LiveKit access token
func (s *UserService) GenerateLiveKitToken(ctx context.Context, identity, name, room, userType string) (string, error) {
	if s.livekitMgr == nil {
		return "", errors.New("LiveKit not configured")
	}

	// Create video grants based on user type
	grants := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}

	// Set additional permissions based on user type
	if userType == "admin" || userType == "host" {
		grants.RoomAdmin = true
		grants.CanPublish = utils.BoolPtr(true)
		grants.CanPublishData = utils.BoolPtr(true)
	} else if userType == "viewer" {
		grants.CanPublish = utils.BoolPtr(false)
		grants.CanPublishData = utils.BoolPtr(false)
	} else {
		grants.CanPublish = utils.BoolPtr(true)
		grants.CanPublishData = utils.BoolPtr(true)
	}

	// Generate token with 24 hour TTL
	return s.livekitMgr.GenerateToken(identity, name, room, grants, "", 24*time.Hour)
}

// GetUserDetail gets user details
func (s *UserService) GetUserDetail(ctx context.Context, room, identity, socketID string) (*models.RoomUser, error) {
	user, err := s.userRepo.GetUserDetail(ctx, room, identity, socketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetUserAlreadyInRoom checks if user is already in room
func (s *UserService) GetUserAlreadyInRoom(ctx context.Context, room, identity string) (bool, error) {
	return s.userRepo.GetUserAlreadyInRoom(ctx, room, identity)
}

// UpdateUserStatus updates user status
func (s *UserService) UpdateUserStatus(ctx context.Context, room, identity, status string) error {
	return s.userRepo.UpdateUserStatus(ctx, room, identity, status)
}

// UpdateUserType updates user type
func (s *UserService) UpdateUserType(ctx context.Context, room, identity, userType string) error {
	return s.userRepo.UpdateUserType(ctx, room, identity, userType)
}

// UpdateUserConference updates user conference status
func (s *UserService) UpdateUserConference(ctx context.Context, identity string, conference int) error {
	return s.userRepo.UpdateUserConference(ctx, identity, conference)
}

// UpdateUserCamera updates user camera status
func (s *UserService) UpdateUserCamera(ctx context.Context, identity string, camera bool) error {
	return s.userRepo.UpdateUserCamera(ctx, identity, camera)
}

// UpdateUserMicrophone updates user microphone status
func (s *UserService) UpdateUserMicrophone(ctx context.Context, identity string, microphone bool) error {
	return s.userRepo.UpdateUserMicrophone(ctx, identity, microphone)
}

// ListUsersInRoom lists all users in a room
func (s *UserService) ListUsersInRoom(ctx context.Context, room, status string) ([]models.RoomUser, error) {
	return s.userRepo.ListUsersInRoom(ctx, room, status)
}

// GetUserRoomAdmin gets admin users in a room
func (s *UserService) GetUserRoomAdmin(ctx context.Context, room string) ([]models.RoomUser, error) {
	return s.userRepo.GetUserRoomAdmin(ctx, room)
}

// AgentList gets count of admin users connected
func (s *UserService) AgentList(ctx context.Context, room string) (int, error) {
	return s.userRepo.AgentList(ctx, room)
}

// RemoveParticipant removes a participant from LiveKit
func (s *UserService) RemoveParticipant(ctx context.Context, room, identity string) error {
	if s.livekitMgr == nil || s.livekitMgr.RoomClient() == nil {
		return errors.New("LiveKit not configured")
	}

	// Update user status in database
	err := s.userRepo.UpdateUserStatus(ctx, room, identity, "disconnect")
	if err != nil {
		// Log but continue
	}

	// Remove from LiveKit
	return s.livekitMgr.RemoveParticipant(ctx, room, identity)
}

// MutePublishedTrack mutes/unmutes a track
func (s *UserService) MutePublishedTrack(ctx context.Context, room, identity, trackSid string, muted bool) error {
	if s.livekitMgr == nil || s.livekitMgr.RoomClient() == nil {
		return errors.New("LiveKit not configured")
	}

	_, err := s.livekitMgr.MutePublishedTrack(ctx, room, identity, trackSid, muted)
	return err
}

// ListParticipants lists participants in a LiveKit room
func (s *UserService) ListParticipants(ctx context.Context, room string) ([]interface{}, error) {
	if s.livekitMgr == nil || s.livekitMgr.RoomClient() == nil {
		return nil, errors.New("LiveKit not configured")
	}

	participants, err := s.livekitMgr.ListParticipants(ctx, room)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice
	result := make([]interface{}, len(participants))
	for i, p := range participants {
		result[i] = p
	}

	return result, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, room, identity, userName, color, socketID string) error {
	return s.userRepo.UpdateUser(ctx, room, identity, userName, color, socketID)
}

// UpdateUserDisconnect clears socket ID on disconnect
func (s *UserService) UpdateUserDisconnect(ctx context.Context, socketID string) error {
	return s.userRepo.UpdateUserDisconnect(ctx, socketID)
}

// GetSocketIDFromIdentity gets socket ID from identity
func (s *UserService) GetSocketIDFromIdentity(ctx context.Context, identity string) (string, error) {
	return s.userRepo.GetSocketIDFromIdentity(ctx, identity)
}

// InitUserExist sets all connected users to disconnect status
func (s *UserService) InitUserExist(ctx context.Context) error {
	return s.userRepo.InitUserExist(ctx)
}
