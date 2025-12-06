package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/umess/backend/internal/middleware"
	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/internal/service"
	"github.com/umess/backend/internal/websocket"
	"github.com/umess/backend/pkg/utils"
)

type MessageHandler struct {
	messageService      *service.MessageService
	conversationService *service.ConversationService
	hub                 *websocket.Hub
}

func NewMessageHandler(hub *websocket.Hub) *MessageHandler {
	return &MessageHandler{
		messageService:      service.NewMessageService(),
		conversationService: service.NewConversationService(),
		hub:                 hub,
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req model.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	message, err := h.messageService.SendMessage(userID, &req)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Broadcast message via WebSocket to conversation participants
	h.broadcastMessage(message)

	utils.SendSuccess(w, "Message sent successfully", message)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	conversationID, err := strconv.Atoi(vars["conversationId"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	messages, err := h.messageService.GetMessages(conversationID, userID, limit, offset)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Messages retrieved successfully", messages)
}

func (h *MessageHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	conversationID, err := strconv.Atoi(vars["conversationId"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	err = h.messageService.MarkAsRead(conversationID, userID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Messages marked as read", nil)
}

func (h *MessageHandler) broadcastMessage(message *model.Message) {
	// Get conversation participants
	conversation, err := h.conversationService.GetConversation(message.ConversationID, message.SenderID)
	if err != nil {
		return // Failed to get conversation
	}

	// Prepare message data
	msgData := map[string]interface{}{
		"type":    "new_message",
		"message": message,
	}

	// Send to all participants except sender
	var participantIDs []int
	for _, participant := range conversation.Participants {
		if participant.ID != message.SenderID {
			participantIDs = append(participantIDs, participant.ID)
		}
	}

	if len(participantIDs) > 0 {
		h.hub.SendToUsers(participantIDs, msgData)
	}
}
