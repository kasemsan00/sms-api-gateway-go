package service

import (
	"context"

	"api-gateway-go/internal/models"
	"api-gateway-go/internal/repository"
)

// ChatService handles chat business logic
type ChatService struct {
	chatRepo *repository.ChatRepository
}

// NewChatService creates a new ChatService
func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

// GetChatHistory gets chat history for a room
func (s *ChatService) GetChatHistory(ctx context.Context, room string, limit int) ([]models.ChatMessage, error) {
	if limit > 0 {
		return s.chatRepo.GetHistoryWithLimit(ctx, room, limit, 0)
	}
	return s.chatRepo.GetHistory(ctx, room)
}

// AddChatMessage adds a new chat message
func (s *ChatService) AddChatMessage(ctx context.Context, params repository.SaveMessageParams) (int64, error) {
	return s.chatRepo.SaveMessage(ctx, params)
}

// GetChatNotification gets chat notification messages
func (s *ChatService) GetChatNotification(ctx context.Context, room string) ([]models.ChatMessage, error) {
	return s.chatRepo.GetNotification(ctx, room)
}

// GetMessageCount gets message count for a room
func (s *ChatService) GetMessageCount(ctx context.Context, room string) (int, error) {
	return s.chatRepo.GetMessageCount(ctx, room)
}

// DeleteMessagesByRoom deletes all messages for a room
func (s *ChatService) DeleteMessagesByRoom(ctx context.Context, room string) error {
	return s.chatRepo.DeleteMessagesByRoom(ctx, room)
}
