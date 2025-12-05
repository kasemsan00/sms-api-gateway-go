package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// AddUserParams holds parameters for adding a user
type AddUserParams struct {
	Room       string
	Identity   string
	UserName   string
	UserType   string
	Status     string
	SocketID   string
	Color      string
	Conference int
	UserAgent  string
}

// AddUser adds a new user to the room
func (r *UserRepository) AddUser(ctx context.Context, params AddUserParams) error {
	dtmCurrent := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO room_user
		(userAgent, room, identity, color, userName, userType, status, socketId, conference, dtmCreated, dtmUpdated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		params.UserAgent,
		params.Room,
		params.Identity,
		params.Color,
		params.UserName,
		params.UserType,
		params.Status,
		params.SocketID,
		params.Conference,
		dtmCurrent,
		dtmCurrent,
	)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, room, identity, userName, color, socketID string) error {
	query := `UPDATE room_user SET userName = ?`
	args := []interface{}{userName}

	if socketID != "" {
		query += `, socketId = ?`
		args = append(args, socketID)
	}
	if color != "" {
		query += `, color = ?`
		args = append(args, color)
	}

	query += ` WHERE room = ? AND identity = ?`
	args = append(args, room, identity)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// UpdateUserStatus updates user status
func (r *UserRepository) UpdateUserStatus(ctx context.Context, room, identity, status string) error {
	query := `UPDATE room_user SET status = ?, camera = 1, microphone = 1 WHERE room = ? AND identity = ?`
	_, err := r.db.ExecContext(ctx, query, status, room, identity)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}

// UpdateUserType updates user type
func (r *UserRepository) UpdateUserType(ctx context.Context, room, identity, userType string) error {
	query := `UPDATE room_user SET userType = ? WHERE room = ? AND identity = ?`
	_, err := r.db.ExecContext(ctx, query, userType, room, identity)
	if err != nil {
		return fmt.Errorf("failed to update user type: %w", err)
	}
	return nil
}

// UpdateUserAgent updates user agent
func (r *UserRepository) UpdateUserAgent(ctx context.Context, room, identity, userAgent string) error {
	query := `UPDATE room_user SET userAgent = ? WHERE room = ? AND identity = ?`
	_, err := r.db.ExecContext(ctx, query, userAgent, room, identity)
	if err != nil {
		return fmt.Errorf("failed to update user agent: %w", err)
	}
	return nil
}

// UpdateUserCamera updates user camera status
func (r *UserRepository) UpdateUserCamera(ctx context.Context, identity string, camera bool) error {
	cameraInt := 0
	if camera {
		cameraInt = 1
	}
	query := `UPDATE room_user SET camera = ? WHERE identity = ?`
	_, err := r.db.ExecContext(ctx, query, cameraInt, identity)
	if err != nil {
		return fmt.Errorf("failed to update user camera: %w", err)
	}
	return nil
}

// UpdateUserMicrophone updates user microphone status
func (r *UserRepository) UpdateUserMicrophone(ctx context.Context, identity string, microphone bool) error {
	micInt := 0
	if microphone {
		micInt = 1
	}
	query := `UPDATE room_user SET microphone = ? WHERE identity = ?`
	_, err := r.db.ExecContext(ctx, query, micInt, identity)
	if err != nil {
		return fmt.Errorf("failed to update user microphone: %w", err)
	}
	return nil
}

// UpdateUserConference updates user conference status
func (r *UserRepository) UpdateUserConference(ctx context.Context, identity string, conference int) error {
	query := `UPDATE room_user SET conference = ? WHERE identity = ?`
	_, err := r.db.ExecContext(ctx, query, conference, identity)
	if err != nil {
		return fmt.Errorf("failed to update user conference: %w", err)
	}
	return nil
}

// UpdateUserDisconnect clears socket ID on disconnect
func (r *UserRepository) UpdateUserDisconnect(ctx context.Context, socketID string) error {
	query := `UPDATE room_user SET socketId = '', camera = 1, microphone = 1 WHERE socketId = ?`
	_, err := r.db.ExecContext(ctx, query, socketID)
	if err != nil {
		return fmt.Errorf("failed to update user disconnect: %w", err)
	}
	return nil
}

// GetUserDetail gets user detail by room and identity or socketID
func (r *UserRepository) GetUserDetail(ctx context.Context, room, identity, socketID string) (*models.RoomUser, error) {
	var user models.RoomUser
	var query string
	var args []interface{}

	if room != "" && identity != "" {
		query = `SELECT room, identity, color, userName, status, camera, microphone, userType, conference
			FROM room_user WHERE room = ? AND identity = ? LIMIT 1`
		args = []interface{}{room, identity}
	} else if socketID != "" {
		query = `SELECT room, identity, color, userName, status, camera, microphone, userType, conference
			FROM room_user WHERE socketId = ? LIMIT 1`
		args = []interface{}{socketID}
	} else {
		return nil, fmt.Errorf("invalid parameters: room and identity or socketID required")
	}

	err := r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user detail: %w", err)
	}

	return &user, nil
}

// IsUserInRoom checks if user exists in room
func (r *UserRepository) IsUserInRoom(ctx context.Context, room, identity string) (int, error) {
	query := `SELECT COUNT(*) FROM room_user WHERE room = ? AND identity = ?`
	var count int
	err := r.db.GetContext(ctx, &count, query, room, identity)
	if err != nil {
		return 0, fmt.Errorf("failed to check user in room: %w", err)
	}
	return count, nil
}

// GetUserAlreadyInRoom checks if user is already in room
func (r *UserRepository) GetUserAlreadyInRoom(ctx context.Context, room, identity string) (bool, error) {
	count, err := r.IsUserInRoom(ctx, room, identity)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListUsersInRoom lists all users in a room with specific status
func (r *UserRepository) ListUsersInRoom(ctx context.Context, room, status string) ([]models.RoomUser, error) {
	var users []models.RoomUser
	query := `SELECT identity, userName, camera, microphone, userType, conference
		FROM room_user WHERE room = ? AND status = ?`

	err := r.db.SelectContext(ctx, &users, query, room, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list users in room: %w", err)
	}
	return users, nil
}

// GetUserRoomAdmin gets admin users in a room
func (r *UserRepository) GetUserRoomAdmin(ctx context.Context, room string) ([]models.RoomUser, error) {
	var users []models.RoomUser
	query := `SELECT socketId FROM room_user WHERE room = ? AND userType = 'admin' AND socketId != ''`

	err := r.db.SelectContext(ctx, &users, query, room)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin users: %w", err)
	}
	return users, nil
}

// GetSocketIDFromIdentity gets socket ID from identity
func (r *UserRepository) GetSocketIDFromIdentity(ctx context.Context, identity string) (string, error) {
	var socketID string
	query := `SELECT socketId FROM room_user WHERE identity = ?`

	err := r.db.GetContext(ctx, &socketID, query, identity)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get socket ID: %w", err)
	}
	return socketID, nil
}

// GetRandomColor gets a random color from color_scheme table
func (r *UserRepository) GetRandomColor(ctx context.Context) (string, error) {
	var colorHex string
	query := `SELECT color_hex FROM color_scheme ORDER BY RAND() LIMIT 1`

	err := r.db.GetContext(ctx, &colorHex, query)
	if err != nil {
		return "", fmt.Errorf("failed to get random color: %w", err)
	}
	return colorHex, nil
}

// AgentList gets count of admin users in a room that are connected
func (r *UserRepository) AgentList(ctx context.Context, room string) (int, error) {
	query := `SELECT COUNT(*) FROM room_user WHERE userType = 'admin' AND room = ? AND status = 'connection'`
	var count int
	err := r.db.GetContext(ctx, &count, query, room)
	if err != nil {
		return 0, fmt.Errorf("failed to count agents: %w", err)
	}
	return count, nil
}

// InitUserExist sets all connected users to disconnect status
func (r *UserRepository) InitUserExist(ctx context.Context) error {
	query := `UPDATE room_user SET status = 'disconnect' WHERE status = 'connection'`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to init user exist: %w", err)
	}
	return nil
}

// GetUserAgent gets user agent by identity
func (r *UserRepository) GetUserAgent(ctx context.Context, identity string) (string, error) {
	var userAgent string
	query := `SELECT userAgent FROM room_user WHERE identity = ? ORDER BY id DESC LIMIT 1`

	err := r.db.GetContext(ctx, &userAgent, query, identity)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get user agent: %w", err)
	}
	return userAgent, nil
}

// UpdateIdentityAndUserName updates identity and username by socketID
func (r *UserRepository) UpdateIdentityAndUserName(ctx context.Context, socketID, identity, userName, status string) error {
	query := `UPDATE room_user SET identity = ?, userName = ?, status = ? WHERE socketId = ?`
	_, err := r.db.ExecContext(ctx, query, identity, userName, status, socketID)
	if err != nil {
		return fmt.Errorf("failed to update identity and username: %w", err)
	}
	return nil
}

// UpdateCameraMicrophoneStatus updates camera/microphone status
func (r *UserRepository) UpdateCameraMicrophoneStatus(ctx context.Context, room, identity, status string) error {
	query := `UPDATE room_user SET cameraMicrophoneStatus = ? WHERE room = ? AND identity = ?`
	_, err := r.db.ExecContext(ctx, query, status, room, identity)
	if err != nil {
		return fmt.Errorf("failed to update camera/microphone status: %w", err)
	}
	return nil
}

// GetCameraMicrophoneStatus gets camera/microphone status
func (r *UserRepository) GetCameraMicrophoneStatus(ctx context.Context, room, identity string) (string, error) {
	var status string
	query := `SELECT cameraMicrophoneStatus FROM room_user WHERE room = ? AND identity = ? LIMIT 1`

	err := r.db.GetContext(ctx, &status, query, room, identity)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get camera/microphone status: %w", err)
	}
	return status, nil
}
