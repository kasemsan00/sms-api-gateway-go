package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// LinkRepository handles link database operations
type LinkRepository struct {
	db *sqlx.DB
}

// NewLinkRepository creates a new LinkRepository
func NewLinkRepository(db *sqlx.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

// CreateLinkParams holds parameters for creating a link
type CreateLinkParams struct {
	Share                 int
	Mobile                string
	LinkID                string
	Room                  string
	RecordID              int
	CrmSender             string
	UserType              string
	LinkType              string
	UserName              string
	IsAdmin               int
	RequireJoinPermission int
	RequireUserName       int
	RequirePassword       int
	OneTimeLink           int
	Password              string
	DtmCreated            string
	DtmExpired            string
	UserAgent             string
}

// Create creates a new link
func (r *LinkRepository) Create(ctx context.Context, params CreateLinkParams) (int64, error) {
	query := `INSERT INTO link_connect(
		share, mobile, linkID, room, recordId, crmSender, userType, linkType, userName,
		isAdmin, requireJoinPermission, requireUserName, requirePassword,
		oneTimeLink, password, dtmCreated, dtmExpired, userAgent
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		params.Share,
		params.Mobile,
		params.LinkID,
		params.Room,
		params.RecordID,
		params.CrmSender,
		params.UserType,
		params.LinkType,
		params.UserName,
		params.IsAdmin,
		params.RequireJoinPermission,
		params.RequireUserName,
		params.RequirePassword,
		params.OneTimeLink,
		params.Password,
		params.DtmCreated,
		params.DtmExpired,
		params.UserAgent,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create link: %w", err)
	}

	return result.LastInsertId()
}

// GetByLinkID gets link detail by linkID
func (r *LinkRepository) GetByLinkID(ctx context.Context, linkID string) (*models.LinkConnect, error) {
	var link models.LinkConnect
	query := `SELECT linkID, room, enabled, mobile, isAdmin, userName, userType, linkType, mobile,
		requireJoinPermission, crmSender, requireUserName, requirePassword, oneTimeLink, dtmCreated, dtmExpired
		FROM link_connect WHERE linkID = ?`

	err := r.db.GetContext(ctx, &link, query, linkID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get link by ID: %w", err)
	}

	return &link, nil
}

// GetByRoomAndUserType gets link by room and userType
func (r *LinkRepository) GetByRoomAndUserType(ctx context.Context, room, userType string) (*models.LinkConnect, error) {
	var link models.LinkConnect
	query := `SELECT linkID, room, isAdmin, enabled, userName, userType, linkType, mobile, requireJoinPermission,
		crmSender, requireUserName, requirePassword, oneTimeLink, dtmCreated, dtmExpired,
		latitude, longitude, patientLatitude, patientLongitude
		FROM link_connect WHERE room = ?`
	args := []interface{}{room}

	if userType != "" {
		query += ` AND userType = ?`
		args = append(args, userType)
	}
	query += ` ORDER BY id DESC LIMIT 1`

	err := r.db.GetContext(ctx, &link, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get link by room: %w", err)
	}

	return &link, nil
}

// GetShareURL gets share URL for a room
func (r *LinkRepository) GetShareURL(ctx context.Context, room, userType string) (string, error) {
	var linkID string
	query := `SELECT linkID FROM link_connect WHERE room = ? AND share = 1 AND userName = ?`

	err := r.db.GetContext(ctx, &linkID, query, room, userType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get share URL: %w", err)
	}

	return linkID, nil
}

// UpdateLatLng updates latitude and longitude
func (r *LinkRepository) UpdateLatLng(ctx context.Context, linkID string, latitude, longitude float64, accuracy int) error {
	query := `UPDATE link_connect SET latitude = ?, longitude = ?, accuracy = ? WHERE linkID = ?`
	_, err := r.db.ExecContext(ctx, query, latitude, longitude, accuracy, linkID)
	if err != nil {
		return fmt.Errorf("failed to update lat/lng: %w", err)
	}
	return nil
}

// UpdatePatientLocation updates patient location
func (r *LinkRepository) UpdatePatientLocation(ctx context.Context, linkID string, patientLat, patientLng float64) error {
	patientUpdated := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE link_connect SET patientLatitude = ?, patientLongitude = ?, patientUpdated = ? WHERE linkID = ?`
	_, err := r.db.ExecContext(ctx, query, patientLat, patientLng, patientUpdated, linkID)
	if err != nil {
		return fmt.Errorf("failed to update patient location: %w", err)
	}
	return nil
}

// GetOneTimeLinkStatus gets one time link status
func (r *LinkRepository) GetOneTimeLinkStatus(ctx context.Context, linkID string) (int, error) {
	var status int
	query := `SELECT oneTimeLink FROM link_connect WHERE linkID = ?`

	err := r.db.GetContext(ctx, &status, query, linkID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get one time link status: %w", err)
	}

	return status, nil
}

// UpdateOneTimeLink updates one time link status
func (r *LinkRepository) UpdateOneTimeLink(ctx context.Context, linkID string, status int) error {
	query := `UPDATE link_connect SET oneTimeLink = ? WHERE linkID = ?`
	_, err := r.db.ExecContext(ctx, query, status, linkID)
	if err != nil {
		return fmt.Errorf("failed to update one time link: %w", err)
	}
	return nil
}

// GetLastLatLng gets last latitude/longitude for a room
func (r *LinkRepository) GetLastLatLng(ctx context.Context, room, userType string, share int) (*models.LinkConnect, error) {
	var link models.LinkConnect
	query := `SELECT latitude, longitude FROM link_connect WHERE room = ? AND userType = ? AND share = ? LIMIT 1`

	err := r.db.GetContext(ctx, &link, query, room, userType, share)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last lat/lng: %w", err)
	}

	return &link, nil
}

// GetUserMobile gets user mobile by room
func (r *LinkRepository) GetUserMobile(ctx context.Context, room string) (string, error) {
	var mobile string
	query := `SELECT mobile FROM link_connect WHERE room = ? AND userType = 'user' LIMIT 1`

	err := r.db.GetContext(ctx, &mobile, query, room)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get user mobile: %w", err)
	}

	return mobile, nil
}

// UpdateUserAgent updates user agent
func (r *LinkRepository) UpdateUserAgent(ctx context.Context, linkID, userAgent, os string) error {
	query := `UPDATE link_connect SET userAgent = ?, os = ? WHERE linkID = ?`
	_, err := r.db.ExecContext(ctx, query, userAgent, os, linkID)
	if err != nil {
		return fmt.Errorf("failed to update user agent: %w", err)
	}
	return nil
}

// UpdateErrorLocation updates error location
func (r *LinkRepository) UpdateErrorLocation(ctx context.Context, linkID, errorMsg string) error {
	query := `UPDATE link_connect SET errorLocation = ? WHERE linkID = ?`
	_, err := r.db.ExecContext(ctx, query, errorMsg, linkID)
	if err != nil {
		return fmt.Errorf("failed to update error location: %w", err)
	}
	return nil
}

// UpdateErrorVideo updates error video
func (r *LinkRepository) UpdateErrorVideo(ctx context.Context, linkID, errorMsg string) error {
	query := `UPDATE link_connect SET errorVideo = ? WHERE linkID = ?`
	_, err := r.db.ExecContext(ctx, query, errorMsg, linkID)
	if err != nil {
		return fmt.Errorf("failed to update error video: %w", err)
	}
	return nil
}

// GetLinkIDList gets links by room and mobile
func (r *LinkRepository) GetLinkIDList(ctx context.Context, room, mobile string) ([]models.LinkConnect, error) {
	var links []models.LinkConnect
	query := `SELECT * FROM link_connect WHERE room = ? AND mobile = ?`

	err := r.db.SelectContext(ctx, &links, query, room, mobile)
	if err != nil {
		return nil, fmt.Errorf("failed to get link list: %w", err)
	}

	return links, nil
}

// UpdateDomainIndex updates domain index
func (r *LinkRepository) UpdateDomainIndex(ctx context.Context, linkID string, index int) error {
	query := `UPDATE link_connect SET domainIndex = ? WHERE linkID = ?`
	_, err := r.db.ExecContext(ctx, query, index, linkID)
	if err != nil {
		return fmt.Errorf("failed to update domain index: %w", err)
	}
	return nil
}

// GetLastDomainIndex gets the last domain index
func (r *LinkRepository) GetLastDomainIndex(ctx context.Context, excludeLinkID string) (int, error) {
	var domainIndex sql.NullInt32
	query := `SELECT IFNULL(domainIndex, 0) as domainIndex FROM link_connect WHERE linkID != ? ORDER BY id DESC LIMIT 1`

	err := r.db.GetContext(ctx, &domainIndex, query, excludeLinkID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get last domain index: %w", err)
	}

	if domainIndex.Valid {
		return int(domainIndex.Int32), nil
	}
	return 0, nil
}

// UpdateLinkEnabled updates link enabled status
func (r *LinkRepository) UpdateLinkEnabled(ctx context.Context, room, linkType string, enabled int) error {
	query := `UPDATE link_connect SET enabled = ? WHERE room = ? AND linkType = ?`
	_, err := r.db.ExecContext(ctx, query, enabled, room, linkType)
	if err != nil {
		return fmt.Errorf("failed to update link enabled: %w", err)
	}
	return nil
}

// GetUserAgentByLinkID gets user agent by linkID
func (r *LinkRepository) GetUserAgentByLinkID(ctx context.Context, linkID string) (string, error) {
	var userAgent string
	query := `SELECT userAgent FROM link_connect WHERE linkID = ? LIMIT 1`

	err := r.db.GetContext(ctx, &userAgent, query, linkID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get user agent: %w", err)
	}

	return userAgent, nil
}

// GetUserAgentByID gets user agent by ID
func (r *LinkRepository) GetUserAgentByID(ctx context.Context, id int) (string, error) {
	var userAgent string
	query := `SELECT userAgent FROM link_connect WHERE id = ? LIMIT 1`

	err := r.db.GetContext(ctx, &userAgent, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get user agent by ID: %w", err)
	}

	return userAgent, nil
}

// GetLatLngGroup gets lat/lng group by room and userType
func (r *LinkRepository) GetLatLngGroup(ctx context.Context, room, userType string) ([]models.LinkConnect, error) {
	var links []models.LinkConnect
	query := `SELECT id, mobile, share, userName, accuracy, latitude, longitude, errorLocation
		FROM link_connect WHERE room = ? AND userType = ? ORDER BY id DESC`

	err := r.db.SelectContext(ctx, &links, query, room, userType)
	if err != nil {
		return nil, fmt.Errorf("failed to get lat/lng group: %w", err)
	}

	return links, nil
}

// UpdateLinkUserName updates link username
func (r *LinkRepository) UpdateLinkUserName(ctx context.Context, linkID, userName string) error {
	query := `UPDATE link_connect SET userName = ? WHERE linkID = ?`
	_, err := r.db.ExecContext(ctx, query, userName, linkID)
	if err != nil {
		return fmt.Errorf("failed to update link username: %w", err)
	}
	return nil
}
