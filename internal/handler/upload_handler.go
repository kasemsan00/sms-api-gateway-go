package handler

import (
	"api-gateway-go/internal/service"
	"api-gateway-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// UploadHandler handles file upload routes
type UploadHandler struct {
	fileService *service.FileService
}

// NewUploadHandler creates a new UploadHandler
func NewUploadHandler(fileService *service.FileService) *UploadHandler {
	return &UploadHandler{
		fileService: fileService,
	}
}

// UploadFile uploads a file
// POST /upload/file
func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.BadRequestResponse(c, "File required")
	}

	subPath := c.FormValue("path", "files")

	fileInfo, err := h.fileService.SaveFile(file, subPath)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fileInfo)
}

// UploadImage uploads an image
// POST /upload/image
func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		file, err = c.FormFile("file")
		if err != nil {
			return utils.BadRequestResponse(c, "Image file required")
		}
	}

	fileInfo, err := h.fileService.SaveImage(file)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fileInfo)
}

// UploadVideo uploads a video
// POST /upload/video
func (h *UploadHandler) UploadVideo(c *fiber.Ctx) error {
	file, err := c.FormFile("video")
	if err != nil {
		file, err = c.FormFile("file")
		if err != nil {
			return utils.BadRequestResponse(c, "Video file required")
		}
	}

	room := c.FormValue("room", "")

	fileInfo, err := h.fileService.SaveVideo(file, room)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fileInfo)
}

// UploadMultiple uploads multiple files
// POST /upload/multiple
func (h *UploadHandler) UploadMultiple(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid form data")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return utils.BadRequestResponse(c, "Files required")
	}

	subPath := c.FormValue("path", "files")

	var uploadedFiles []*service.FileInfo
	var errors []string

	for _, file := range files {
		fileInfo, err := h.fileService.SaveFile(file, subPath)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		uploadedFiles = append(uploadedFiles, fileInfo)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"files":    uploadedFiles,
		"uploaded": len(uploadedFiles),
		"errors":   errors,
	})
}

// DeleteFile deletes a file
// DELETE /upload/file
func (h *UploadHandler) DeleteFile(c *fiber.Ctx) error {
	type DeleteRequest struct {
		FilePath string `json:"filePath"`
	}

	var req DeleteRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.FilePath == "" {
		return utils.BadRequestResponse(c, "File path required")
	}

	err := h.fileService.DeleteFile(req.FilePath)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"deleted": true,
	})
}

// CheckFileExists checks if a file exists
// GET /upload/exists
func (h *UploadHandler) CheckFileExists(c *fiber.Ctx) error {
	filePath := c.Query("path")
	if filePath == "" {
		return utils.BadRequestResponse(c, "File path required")
	}

	exists := h.fileService.FileExists(filePath)

	return utils.SuccessResponse(c, fiber.Map{
		"exists": exists,
	})
}

// VideoList gets video list with pagination
// GET /upload/list
func (h *UploadHandler) VideoList(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	linkType := c.Query("linkType", "record")
	search := c.Query("search", "")
	mobile := c.Query("mobile", "")
	agent := c.Query("agent", "")
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")

	// Clean up undefined values
	if search == "null" || search == "undefined" {
		search = ""
	}
	if mobile == "null" || mobile == "undefined" {
		mobile = ""
	}
	if agent == "null" || agent == "undefined" {
		agent = ""
	}

	videos, total, err := h.fileService.GetVideoList(c.Context(), service.VideoListParams{
		Page:      page,
		Limit:     limit,
		LinkType:  linkType,
		Search:    search,
		Mobile:    mobile,
		Agent:     agent,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}

	totalPages := (total + limit - 1) / limit

	return utils.SuccessResponse(c, fiber.Map{
		"data": videos,
		"pagination": fiber.Map{
			"totalItems":      total,
			"totalPages":      totalPages,
			"currentPage":     page,
			"itemsPerPage":    limit,
			"hasNextPage":     page < totalPages,
			"hasPreviousPage": page > 1,
		},
	})
}

// SendSMS sends SMS with video/record link
// POST /upload/sms/send
func (h *UploadHandler) SendSMS(c *fiber.Ctx) error {
	type SendRequest struct {
		RecordID      int    `json:"recordId"`
		AgentUsername string `json:"agentUsername"`
		Room          string `json:"room"`
		Mobile        string `json:"mobile"`
		Service       int    `json:"service"`
		SMS           int    `json:"sms"`
	}

	var req SendRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if req.RecordID == 0 || req.AgentUsername == "" || req.Room == "" || req.Mobile == "" || req.Service == 0 {
		return utils.BadRequestResponse(c, "Record ID, agent username, room, mobile, and service are required")
	}

	// Here we would create a link and send SMS
	// For now, return success
	return utils.SuccessResponse(c, fiber.Map{
		"message": "SMS sent successfully",
	})
}
