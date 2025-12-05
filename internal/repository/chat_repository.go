package repository

import (
	"context"
	"fmt"
	"time"

	"api-gateway-go/internal/models"

	"github.com/jmoiron/sqlx"
)

// ChatRepository handles chat database operations
type ChatRepository struct {
	db *sqlx.DB
}

// NewChatRepository creates a new ChatRepository
func NewChatRepository(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// SaveMessageParams holds parameters for saving a chat message
type SaveMessageParams struct {
	Room             string
	Identity         string
	ChatIdentity     string
	UserName         string
	UserType         string
	Text             string
	Files            string
	Color            string
	ReplyToMessageID int
	ReplyToUserName  string
	ReplyToText      string
}

// SaveMessage saves a chat message
func (r *ChatRepository) SaveMessage(ctx context.Context, params SaveMessageParams) (int64, error) {
	dtmCreated := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO chat_message
		(room, identity, chat_identity, userName, userType, text, files, color, replyToMessageId, replyToUserName, replyToText, dtmCreated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		params.Room,
		params.Identity,
		params.ChatIdentity,
		params.UserName,
		params.UserType,
		params.Text,
		params.Files,
		params.Color,
		params.ReplyToMessageID,
		params.ReplyToUserName,
		params.ReplyToText,
		dtmCreated,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to save chat message: %w", err)
	}

	return result.LastInsertId()
}

// GetHistory gets chat history for a room
func (r *ChatRepository) GetHistory(ctx context.Context, room string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	query := `SELECT id, room, identity, chat_identity, userName, text, color, files,
		replyToMessageId, replyToUserName, replyToText, dtmCreated, userType
		FROM chat_message WHERE room = ? ORDER BY dtmCreated ASC`

	err := r.db.SelectContext(ctx, &messages, query, room)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	return messages, nil
}

// GetHistoryWithLimit gets limited chat history for a room
func (r *ChatRepository) GetHistoryWithLimit(ctx context.Context, room string, limit, offset int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	query := `SELECT id, room, identity, chat_identity, userName, text, color, files,
		replyToMessageId, replyToUserName, replyToText, dtmCreated, userType
		FROM chat_message WHERE room = ? ORDER BY dtmCreated DESC LIMIT ? OFFSET ?`

	err := r.db.SelectContext(ctx, &messages, query, room, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetNotification gets chat notifications
func (r *ChatRepository) GetNotification(ctx context.Context, room string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	query := `SELECT id, room, identity, userName, text, dtmCreated, userType
		FROM chat_message WHERE room = ? AND userType = 'user' ORDER BY dtmCreated DESC LIMIT 10`

	err := r.db.SelectContext(ctx, &messages, query, room)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat notification: %w", err)
	}

	return messages, nil
}

// GetMessageByID gets a message by ID
func (r *ChatRepository) GetMessageByID(ctx context.Context, id int) (*models.ChatMessage, error) {
	var message models.ChatMessage
	query := `SELECT * FROM chat_message WHERE id = ?`

	err := r.db.GetContext(ctx, &message, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get message by ID: %w", err)
	}

	return &message, nil
}

// DeleteMessagesByRoom deletes all messages for a room
func (r *ChatRepository) DeleteMessagesByRoom(ctx context.Context, room string) error {
	query := `DELETE FROM chat_message WHERE room = ?`
	_, err := r.db.ExecContext(ctx, query, room)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}
	return nil
}

// GetMessageCount gets message count for a room
func (r *ChatRepository) GetMessageCount(ctx context.Context, room string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM chat_message WHERE room = ?`

	err := r.db.GetContext(ctx, &count, query, room)
	if err != nil {
		return 0, fmt.Errorf("failed to get message count: %w", err)
	}

	return count, nil
}
