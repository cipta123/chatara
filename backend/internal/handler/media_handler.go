package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/umess/backend/internal/middleware"
	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/internal/service"
	"github.com/umess/backend/pkg/config"
	"github.com/umess/backend/pkg/utils"
)

type MediaHandler struct {
	messageService *service.MessageService
}

func NewMediaHandler() *MediaHandler {
	return &MediaHandler{
		messageService: service.NewMessageService(),
	}
}

func (h *MediaHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "No image file provided")
		return
	}
	defer file.Close()

	// Validate file type
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		utils.SendError(w, http.StatusBadRequest, "File must be an image")
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := config.AppConfig.MediaStorage
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Generate URL (for development, use local path)
	mediaURL := fmt.Sprintf("/api/media/%s", filename)

	response := map[string]interface{}{
		"url":      mediaURL,
		"file_path": filePath,
		"file_name": header.Filename,
		"file_size": header.Size,
		"mime_type": header.Header.Get("Content-Type"),
	}

	utils.SendSuccess(w, "Image uploaded successfully", response)
}

func (h *MediaHandler) ServeMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	filename := vars["filename"]

	if filename == "" {
		utils.SendError(w, http.StatusBadRequest, "Filename required")
		return
	}

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		utils.SendError(w, http.StatusBadRequest, "Invalid filename")
		return
	}

	filePath := filepath.Join(config.AppConfig.MediaStorage, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.SendError(w, http.StatusNotFound, "File not found")
		return
	}

	// Serve file
	http.ServeFile(w, r, filePath)
}

func (h *MediaHandler) SendMessageWithImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		ConversationID int    `json:"conversation_id"`
		MediaURL       string `json:"media_url"`
		Caption        string `json:"caption,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	messageReq := &model.SendMessageRequest{
		ConversationID: req.ConversationID,
		Content:        req.Caption,
		Type:           model.MessageTypeImage,
		MediaURL:       req.MediaURL,
	}

	message, err := h.messageService.SendMessage(userID, messageReq)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Message with image sent successfully", message)
}


