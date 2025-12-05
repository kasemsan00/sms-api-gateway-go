package service

import (
	"context"
	"database/sql"
	"errors"

	"api-gateway-go/internal/models"
	"api-gateway-go/internal/repository"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

// NotificationService handles notification business logic
type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(notificationRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

// GetAllNotifications gets all notifications
func (s *NotificationService) GetAllNotifications(ctx context.Context, limit, offset int) ([]models.Notification, error) {
	if limit > 0 {
		return s.notificationRepo.GetWithPagination(ctx, limit, offset)
	}
	return s.notificationRepo.GetAll(ctx)
}

// GetNotificationsByUserName gets notifications by username
func (s *NotificationService) GetNotificationsByUserName(ctx context.Context, userName string) ([]models.Notification, error) {
	return s.notificationRepo.GetByUserName(ctx, userName)
}

// GetUnreadNotifications gets unread notifications
func (s *NotificationService) GetUnreadNotifications(ctx context.Context) ([]models.Notification, error) {
	return s.notificationRepo.GetUnread(ctx)
}

// GetNotificationByID gets a notification by ID
func (s *NotificationService) GetNotificationByID(ctx context.Context, notificationID int) (*models.Notification, error) {
	notification, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}
	if notification == nil {
		return nil, ErrNotificationNotFound
	}
	return notification, nil
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(ctx context.Context, params repository.CreateNotificationParams) (int64, error) {
	return s.notificationRepo.Create(ctx, params)
}

// UpdateReadStatus marks a notification as read
func (s *NotificationService) UpdateReadStatus(ctx context.Context, notificationID int) error {
	return s.notificationRepo.UpdateReadStatus(ctx, notificationID, 1)
}

// GetUnreadCount gets count of unread notifications
func (s *NotificationService) GetUnreadCount(ctx context.Context) (int, error) {
	return s.notificationRepo.GetUnreadCount(ctx)
}

// MarkAllAsRead marks all notifications as read
func (s *NotificationService) MarkAllAsRead(ctx context.Context) error {
	return s.notificationRepo.MarkAllAsRead(ctx)
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID int) error {
	return s.notificationRepo.Delete(ctx, notificationID)
}
