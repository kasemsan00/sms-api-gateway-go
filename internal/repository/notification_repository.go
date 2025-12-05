package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// NotificationRepository handles notification database operations
type NotificationRepository struct {
	db *sqlx.DB
}

// NewNotificationRepository creates a new NotificationRepository
func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// CreateNotificationParams holds parameters for creating a notification
type CreateNotificationParams struct {
	UserName         string
	Mobile           string
	Message          string
	CaseID           int
	NotificationType string
	RelatedURL       string
}

// Create creates a new notification
func (r *NotificationRepository) Create(ctx context.Context, params CreateNotificationParams) (int64, error) {
	dtmCreated := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO notification
		(userName, mobile, message, caseId, notificationType, relatedUrl, dtmCreated)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		params.UserName,
		params.Mobile,
		params.Message,
		params.CaseID,
		params.NotificationType,
		params.RelatedURL,
		dtmCreated,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create notification: %w", err)
	}

	return result.LastInsertId()
}

// GetByID gets notification by ID
func (r *NotificationRepository) GetByID(ctx context.Context, notificationID int) (*models.Notification, error) {
	var notification models.Notification
	query := `SELECT * FROM notification WHERE notificationId = ?`

	err := r.db.GetContext(ctx, &notification, query, notificationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	return &notification, nil
}

// GetAll gets all notifications
func (r *NotificationRepository) GetAll(ctx context.Context) ([]models.Notification, error) {
	var notifications []models.Notification
	query := `SELECT * FROM notification ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &notifications, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all notifications: %w", err)
	}

	return notifications, nil
}

// GetUnread gets all unread notifications
func (r *NotificationRepository) GetUnread(ctx context.Context) ([]models.Notification, error) {
	var notifications []models.Notification
	query := `SELECT * FROM notification WHERE ` + "`read`" + ` = 0 ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &notifications, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %w", err)
	}

	return notifications, nil
}

// UpdateReadStatus updates notification read status
func (r *NotificationRepository) UpdateReadStatus(ctx context.Context, notificationID int, read int) error {
	dtmRead := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE notification SET ` + "`read`" + ` = ?, dtmRead = ? WHERE notificationId = ?`

	_, err := r.db.ExecContext(ctx, query, read, dtmRead, notificationID)
	if err != nil {
		return fmt.Errorf("failed to update read status: %w", err)
	}

	return nil
}

// Delete deletes a notification
func (r *NotificationRepository) Delete(ctx context.Context, notificationID int) error {
	query := `DELETE FROM notification WHERE notificationId = ?`

	_, err := r.db.ExecContext(ctx, query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}

// GetByUserName gets notifications by username
func (r *NotificationRepository) GetByUserName(ctx context.Context, userName string) ([]models.Notification, error) {
	var notifications []models.Notification
	query := `SELECT * FROM notification WHERE userName = ? ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &notifications, query, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications by username: %w", err)
	}

	return notifications, nil
}

// GetByCaseID gets notifications by case ID
func (r *NotificationRepository) GetByCaseID(ctx context.Context, caseID int) ([]models.Notification, error) {
	var notifications []models.Notification
	query := `SELECT * FROM notification WHERE caseId = ? ORDER BY dtmCreated DESC`

	err := r.db.SelectContext(ctx, &notifications, query, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications by case ID: %w", err)
	}

	return notifications, nil
}

// GetWithPagination gets notifications with pagination
func (r *NotificationRepository) GetWithPagination(ctx context.Context, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	query := `SELECT * FROM notification ORDER BY dtmCreated DESC LIMIT ? OFFSET ?`

	err := r.db.SelectContext(ctx, &notifications, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}

// GetUnreadCount gets count of unread notifications
func (r *NotificationRepository) GetUnreadCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notification WHERE ` + "`read`" + ` = 0`

	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

// MarkAllAsRead marks all notifications as read
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context) error {
	dtmRead := time.Now().Format("2006-01-02 15:04:05")
	query := `UPDATE notification SET ` + "`read`" + ` = 1, dtmRead = ? WHERE ` + "`read`" + ` = 0`

	_, err := r.db.ExecContext(ctx, query, dtmRead)
	if err != nil {
		return fmt.Errorf("failed to mark all as read: %w", err)
	}

	return nil
}
