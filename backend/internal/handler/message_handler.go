package handler

import (
	"encoding/json"
	"log"
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

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	messageID, err := strconv.Atoi(vars["messageId"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	// Verify conversation exists and user is participant
	_, err = h.conversationService.GetConversation(conversationID, userID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Conversation not found")
		return
	}

	// Delete message
	err = h.messageService.DeleteMessage(messageID, userID)
	if err != nil {
		if err.Error() == "message not found" || err.Error() == "unauthorized to delete message" {
			utils.SendError(w, http.StatusForbidden, err.Error())
		} else {
			utils.SendError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// Broadcast deletion via WebSocket
	conversation, err := h.conversationService.GetConversation(conversationID, userID)
	if err == nil {
		var participantIDs []int
		for _, participant := range conversation.Participants {
			participantIDs = append(participantIDs, participant.ID)
		}
		if len(participantIDs) > 0 {
			h.hub.SendToUsers(participantIDs, map[string]interface{}{
				"type":           "message_deleted",
				"message_id":     messageID,
				"conversation_id": conversationID,
			})
		}
	}

	utils.SendSuccess(w, "Message deleted successfully", map[string]interface{}{
		"message_id": messageID,
	})
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
		log.Printf("Failed to get conversation %d for broadcast: %v", message.ConversationID, err)
		return // Failed to get conversation
	}

	// Prepare message data
	msgData := map[string]interface{}{
		"type":    "new_message",
		"message": message,
	}

	// Send to ALL participants (including sender for confirmation)
	var participantIDs []int
	for _, participant := range conversation.Participants {
		participantIDs = append(participantIDs, participant.ID)
	}

	log.Printf("Broadcasting message %d from user %d to %d participants in conversation %d: %v", 
		message.ID, message.SenderID, len(participantIDs), message.ConversationID, participantIDs)

	if len(participantIDs) > 0 {
		err := h.hub.SendToUsers(participantIDs, msgData)
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
		}
	}
}
