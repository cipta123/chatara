package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/umess/backend/internal/middleware"
	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/internal/service"
	"github.com/umess/backend/pkg/utils"
)

type ConversationHandler struct {
	conversationService *service.ConversationService
}

func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{
		conversationService: service.NewConversationService(),
	}
}

func (h *ConversationHandler) CreateDirectConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID1, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		UserID2 int `json:"user_id_2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	conversation, err := h.conversationService.CreateDirectConversation(userID1, req.UserID2)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Conversation created successfully", conversation)
}

func (h *ConversationHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversations, err := h.conversationService.GetUserConversations(userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(w, "Conversations retrieved successfully", conversations)
}

func (h *ConversationHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
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

	conversation, err := h.conversationService.GetConversation(conversationID, userID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Conversation retrieved successfully", conversation)
}

func (h *ConversationHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req model.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	conversation, err := h.conversationService.CreateGroup(userID, req.Name, req.Description, req.ParticipantIDs)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Group created successfully", conversation)
}

func (h *ConversationHandler) AddGroupMember(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.conversationService.AddGroupMember(conversationID, req.UserID, userID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Member added successfully", nil)
}

func (h *ConversationHandler) RemoveGroupMember(w http.ResponseWriter, r *http.Request) {
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
	memberID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.conversationService.RemoveGroupMember(conversationID, memberID, userID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Member removed successfully", nil)
}

func (h *ConversationHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
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

	err = h.conversationService.LeaveGroup(conversationID, userID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Left group successfully", nil)
}

func (h *ConversationHandler) UpdateGroupInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.conversationService.UpdateGroupInfo(conversationID, userID, req.Name, req.Description)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, "Group info updated successfully", nil)
}


