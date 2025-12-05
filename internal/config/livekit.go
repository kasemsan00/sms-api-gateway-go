package config

import (
	"context"
	"fmt"
	"time"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"github.com/rs/zerolog/log"
)

// LiveKitManager manages LiveKit server connections
type LiveKitManager struct {
	cfg          *Config
	roomClient   *lksdk.RoomServiceClient
	egressClient *lksdk.EgressClient
}

var livekitManager *LiveKitManager

// InitLiveKit initializes the LiveKit clients
func InitLiveKit(cfg *Config) (*LiveKitManager, error) {
	if cfg.LiveKitHost == "" || cfg.LiveKitAPIKey == "" || cfg.LiveKitAPISecret == "" {
		log.Warn().Msg("LiveKit configuration is incomplete, skipping LiveKit initialization")
		return &LiveKitManager{cfg: cfg}, nil
	}

	roomClient := lksdk.NewRoomServiceClient(cfg.LiveKitHost, cfg.LiveKitAPIKey, cfg.LiveKitAPISecret)
	egressClient := lksdk.NewEgressClient(cfg.LiveKitHost, cfg.LiveKitAPIKey, cfg.LiveKitAPISecret)

	livekitManager = &LiveKitManager{
		cfg:          cfg,
		roomClient:   roomClient,
		egressClient: egressClient,
	}

	log.Info().Msgf("Connected to LiveKit: %s", cfg.LiveKitHost)

	return livekitManager, nil
}

// GetLiveKit returns the LiveKit manager
func GetLiveKit() *LiveKitManager {
	return livekitManager
}

// RoomClient returns the Room Service client
func (lm *LiveKitManager) RoomClient() *lksdk.RoomServiceClient {
	return lm.roomClient
}

// EgressClient returns the Egress client
func (lm *LiveKitManager) EgressClient() *lksdk.EgressClient {
	return lm.egressClient
}

// CreateRoom creates a new LiveKit room
func (lm *LiveKitManager) CreateRoom(ctx context.Context, name string, emptyTimeout, maxParticipants uint32) (*livekit.Room, error) {
	if lm.roomClient == nil {
		return nil, fmt.Errorf("LiveKit room client not initialized")
	}

	req := &livekit.CreateRoomRequest{
		Name:            name,
		EmptyTimeout:    emptyTimeout,
		MaxParticipants: maxParticipants,
	}

	room, err := lm.roomClient.CreateRoom(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	log.Info().Msgf("Created LiveKit room: %s", name)
	return room, nil
}

// DeleteRoom deletes a LiveKit room
func (lm *LiveKitManager) DeleteRoom(ctx context.Context, name string) error {
	if lm.roomClient == nil {
		return fmt.Errorf("LiveKit room client not initialized")
	}

	req := &livekit.DeleteRoomRequest{
		Room: name,
	}

	_, err := lm.roomClient.DeleteRoom(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	log.Info().Msgf("Deleted LiveKit room: %s", name)
	return nil
}

// ListRooms lists all active rooms
func (lm *LiveKitManager) ListRooms(ctx context.Context) ([]*livekit.Room, error) {
	if lm.roomClient == nil {
		return nil, fmt.Errorf("LiveKit room client not initialized")
	}

	req := &livekit.ListRoomsRequest{}
	resp, err := lm.roomClient.ListRooms(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %w", err)
	}

	return resp.Rooms, nil
}

// ListParticipants lists participants in a room
func (lm *LiveKitManager) ListParticipants(ctx context.Context, room string) ([]*livekit.ParticipantInfo, error) {
	if lm.roomClient == nil {
		return nil, fmt.Errorf("LiveKit room client not initialized")
	}

	req := &livekit.ListParticipantsRequest{
		Room: room,
	}

	resp, err := lm.roomClient.ListParticipants(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list participants: %w", err)
	}

	return resp.Participants, nil
}

// GetParticipant gets a specific participant
func (lm *LiveKitManager) GetParticipant(ctx context.Context, room, identity string) (*livekit.ParticipantInfo, error) {
	if lm.roomClient == nil {
		return nil, fmt.Errorf("LiveKit room client not initialized")
	}

	req := &livekit.RoomParticipantIdentity{
		Room:     room,
		Identity: identity,
	}

	participant, err := lm.roomClient.GetParticipant(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get participant: %w", err)
	}

	return participant, nil
}

// RemoveParticipant removes a participant from a room
func (lm *LiveKitManager) RemoveParticipant(ctx context.Context, room, identity string) error {
	if lm.roomClient == nil {
		return fmt.Errorf("LiveKit room client not initialized")
	}

	req := &livekit.RoomParticipantIdentity{
		Room:     room,
		Identity: identity,
	}

	_, err := lm.roomClient.RemoveParticipant(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	log.Info().Msgf("Removed participant %s from room %s", identity, room)
	return nil
}

// UpdateParticipant updates participant metadata
func (lm *LiveKitManager) UpdateParticipant(ctx context.Context, room, identity, metadata string) (*livekit.ParticipantInfo, error) {
	if lm.roomClient == nil {
		return nil, fmt.Errorf("LiveKit room client not initialized")
	}

	req := &livekit.UpdateParticipantRequest{
		Room:     room,
		Identity: identity,
		Metadata: metadata,
	}

	participant, err := lm.roomClient.UpdateParticipant(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update participant: %w", err)
	}

	return participant, nil
}

// MutePublishedTrack mutes/unmutes a published track
func (lm *LiveKitManager) MutePublishedTrack(ctx context.Context, room, identity, trackSid string, muted bool) (*livekit.MuteRoomTrackResponse, error) {
	if lm.roomClient == nil {
		return nil, fmt.Errorf("LiveKit room client not initialized")
	}

	req := &livekit.MuteRoomTrackRequest{
		Room:     room,
		Identity: identity,
		TrackSid: trackSid,
		Muted:    muted,
	}

	resp, err := lm.roomClient.MutePublishedTrack(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to mute track: %w", err)
	}

	return resp, nil
}

// GenerateToken generates an access token for a participant
func (lm *LiveKitManager) GenerateToken(identity, name, room string, grants *lksdk.VideoGrant, metadata string, ttl time.Duration) (string, error) {
	at := lksdk.NewAccessToken(lm.cfg.LiveKitAPIKey, lm.cfg.LiveKitAPISecret)

	at.SetIdentity(identity).
		SetName(name).
		SetMetadata(metadata).
		SetValidFor(ttl).
		AddGrant(grants)

	return at.ToJWT()
}

// StartRoomCompositeEgress starts room composite egress recording
func (lm *LiveKitManager) StartRoomCompositeEgress(ctx context.Context, room string, output *livekit.EncodedFileOutput) (*livekit.EgressInfo, error) {
	if lm.egressClient == nil {
		return nil, fmt.Errorf("LiveKit egress client not initialized")
	}

	req := &livekit.RoomCompositeEgressRequest{
		RoomName: room,
		Output: &livekit.RoomCompositeEgressRequest_File{
			File: output,
		},
	}

	info, err := lm.egressClient.StartRoomCompositeEgress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start room composite egress: %w", err)
	}

	log.Info().Msgf("Started room composite egress for room %s: %s", room, info.EgressId)
	return info, nil
}

// StopEgress stops an active egress
func (lm *LiveKitManager) StopEgress(ctx context.Context, egressID string) (*livekit.EgressInfo, error) {
	if lm.egressClient == nil {
		return nil, fmt.Errorf("LiveKit egress client not initialized")
	}

	req := &livekit.StopEgressRequest{
		EgressId: egressID,
	}

	info, err := lm.egressClient.StopEgress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to stop egress: %w", err)
	}

	log.Info().Msgf("Stopped egress: %s", egressID)
	return info, nil
}

// ListEgress lists active egresses
func (lm *LiveKitManager) ListEgress(ctx context.Context, room string) ([]*livekit.EgressInfo, error) {
	if lm.egressClient == nil {
		return nil, fmt.Errorf("LiveKit egress client not initialized")
	}

	req := &livekit.ListEgressRequest{
		RoomName: room,
	}

	resp, err := lm.egressClient.ListEgress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list egress: %w", err)
	}

	return resp.Items, nil
}

// Health checks LiveKit health
func (lm *LiveKitManager) Health(ctx context.Context) error {
	if lm.roomClient == nil {
		return fmt.Errorf("LiveKit not configured")
	}

	_, err := lm.ListRooms(ctx)
	return err
}
