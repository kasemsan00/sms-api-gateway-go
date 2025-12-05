package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"api-gateway-go/internal/config"
	"api-gateway-go/pkg/utils"
)

// FileInfo represents uploaded file info
type FileInfo struct {
	OriginalName string `json:"originalName"`
	FileName     string `json:"fileName"`
	FilePath     string `json:"filePath"`
	FileSize     int64  `json:"fileSize"`
	MimeType     string `json:"mimeType"`
	URL          string `json:"url"`
}

// FileService handles file operations
type FileService struct {
	cfg        *config.Config
	uploadPath string
}

// NewFileService creates a new FileService
func NewFileService(cfg *config.Config) *FileService {
	return &FileService{
		cfg:        cfg,
		uploadPath: "./uploads",
	}
}

// SaveFile saves an uploaded file
func (s *FileService) SaveFile(file *multipart.FileHeader, subPath string) (*FileInfo, error) {
	// Check file size
	if file.Size > s.cfg.FileSizeLimit {
		return nil, fmt.Errorf("file size exceeds limit of %d bytes", s.cfg.FileSizeLimit)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	newFileName := utils.GenerateRandomString(16, utils.AlphanumericCharset) + ext

	// Create directory if it doesn't exist
	saveDir := filepath.Join(s.uploadPath, subPath)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, err
	}

	// Open source file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Create destination file
	dstPath := filepath.Join(saveDir, newFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy file content
	written, err := io.Copy(dst, src)
	if err != nil {
		os.Remove(dstPath)
		return nil, err
	}

	// Get MIME type from extension
	mimeType := s.getMimeType(ext)

	return &FileInfo{
		OriginalName: file.Filename,
		FileName:     newFileName,
		FilePath:     dstPath,
		FileSize:     written,
		MimeType:     mimeType,
		URL:          s.cfg.APIURL + "/" + subPath + "/" + newFileName,
	}, nil
}

// SaveVideo saves an uploaded video file
func (s *FileService) SaveVideo(file *multipart.FileHeader, room string) (*FileInfo, error) {
	return s.SaveFile(file, "videos")
}

// SaveImage saves an uploaded image file
func (s *FileService) SaveImage(file *multipart.FileHeader) (*FileInfo, error) {
	return s.SaveFile(file, "images")
}

// DeleteFile deletes a file
func (s *FileService) DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// FileExists checks if a file exists
func (s *FileService) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// GetRecordPath returns the record file path
func (s *FileService) GetRecordPath() string {
	return s.cfg.RecordPath
}

// GenerateRecordFilename generates a filename for a recording
func (s *FileService) GenerateRecordFilename(room string) string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s_%s.mp4", room, timestamp)
}

// getMimeType returns MIME type based on file extension
func (s *FileService) getMimeType(ext string) string {
	ext = strings.ToLower(ext)
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	if mime, ok := mimeTypes[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}

// VideoListParams represents parameters for video list query
type VideoListParams struct {
	Page      int
	Limit     int
	LinkType  string
	Search    string
	Mobile    string
	Agent     string
	StartDate string
	EndDate   string
}

// VideoInfo represents video information
type VideoInfo struct {
	ID        int    `json:"id"`
	Room      string `json:"room"`
	FileName  string `json:"fileName"`
	FilePath  string `json:"filePath"`
	FileSize  int64  `json:"fileSize"`
	Duration  int    `json:"duration"`
	URL       string `json:"url"`
	CreatedAt string `json:"createdAt"`
}

// GetVideoList gets list of videos with pagination
// This is a placeholder implementation - actual implementation would query database
func (s *FileService) GetVideoList(ctx context.Context, params VideoListParams) ([]VideoInfo, int, error) {
	// Placeholder - This would normally query the database for video records
	videos := []VideoInfo{}
	total := 0

	return videos, total, nil
}
